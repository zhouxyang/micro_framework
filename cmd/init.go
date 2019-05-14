package cmd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	grpc_requestid "route_guide/middleware/grpc_requestid"
	"time"
)

func InitGrpcLog() *logrus.Entry {
	// init grpc log
	/*LogFile := path.Join(Config.DockerInsideLogDir, path.Clean(Config.GrpcLogFile))
	grpcLog, err := os.OpenFile(LogFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("open logfile failed: %v(%v)", LogFile, err)
		os.Exit(-2)
	}
	*/
	entry := logrus.New()
	entry.Formatter = &logrus.JSONFormatter{}
	entry.Out = os.Stdout

	log := entry.WithFields(logrus.Fields{
		"pid": os.Getpid(),
	})

	log.Infof("init grpclog success")
	return log
}

func GetLog(ctx context.Context) *logrus.Entry {
	return ctxlogrus.Extract(ctx)
}

// GetEtcdClient 创建etcd连接
func GetEtcdClient(ip string, port int) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		// etcd集群成员节点列表
		Endpoints:   []string{fmt.Sprintf("%s:%d", ip, port)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		werr := errors.Wrap(err, "connect etcd err")
		return nil, werr
	}
	return cli, nil
}

// GetGrpcConn 创建grpc连接
func GetGrpcConn(ctx context.Context, service string, cli *clientv3.Client, log *logrus.Entry) (*grpc.ClientConn, error) {
	r := &etcdnaming.GRPCResolver{Client: cli}
	b := grpc.RoundRobin(r)
	conn, err := grpc.DialContext(
		ctx,
		service,
		grpc.WithInsecure(),
		grpc.WithBalancer(b),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_logrus.UnaryClientInterceptor(log),
				grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(3)),
				grpc_requestid.UnaryClientInterceptor(log),
			),
		),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				grpc_logrus.StreamClientInterceptor(log),
				grpc_retry.StreamClientInterceptor(grpc_retry.WithMax(3)),
				grpc_requestid.StreamClientInterceptor(log),
			),
		),
	)

	if err != nil {
		errors.Wrap(err, "fail to dial")
		return nil, err
	}
	return conn, err
}
