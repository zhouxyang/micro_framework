package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"route_guide/cmd"
	grpc_requestid "route_guide/middleware/grpc_requestid"
	"syscall"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/naming"

	_ "route_guide/cmd/route"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

func main() {
	flag.Parse()
	log := cmd.InitGrpcLog()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	register(log)
	start(log, lis)
}
func register(log *logrus.Entry) {
	cli, err := cmd.GetEtcdClient("localhost", 2379)
	if err != nil {
		log.Infof("GetEtcdClient err:", err)
		return
	}
	// 创建命名解析
	r := &etcdnaming.GRPCResolver{Client: cli}
	// 将本服务注册添加etcd中，服务名称为myService，服务地址为本机10000端口

	for name := range cmd.ServerDrivers {
		err = r.Update(context.TODO(), name, naming.Update{Op: naming.Add, Addr: fmt.Sprintf("127.0.0.1:%d", *port)})
		if err != nil {
			log.Infof("[测] update etcd err:", err)
			return
		}
	}

}

func unRegister(log *logrus.Entry) {
	cli, err := cmd.GetEtcdClient("localhost", 2379)
	if err != nil {
		log.Infof("GetEtcdClient err:", err)
		return
	}
	// 创建命名解析
	r := &etcdnaming.GRPCResolver{Client: cli}
	for name := range cmd.ServerDrivers {
		// 将本服务注册添加etcd中，服务名称为myService，服务地址为本机10000端口
		err = r.Update(context.TODO(), name, naming.Update{Op: naming.Delete, Addr: fmt.Sprintf("127.0.0.1:%d", *port)})
		if err != nil {
			log.Infof("[测] update etcd err:", err)
			return
		}
	}

}

func start(log *logrus.Entry, lis net.Listener) {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(log),
			grpc_requestid.StreamServerInterceptor(log),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(log),
			grpc_requestid.UnaryServerInterceptor(log),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	for _, register := range cmd.ServerDrivers {
		register(grpcServer)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for s := range c {
			unRegister(log)
			log.Infof("Got signal: %v", s)
			grpcServer.GracefulStop()
		}

	}()
	log.Infof("started")
	grpcServer.Serve(lis)
}
