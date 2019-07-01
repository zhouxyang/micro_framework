package main

import (
	"flag"
	"fmt"
	"micro_framework/cmd"
	"micro_framework/configfile"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	grpc_requestid "micro_framework/middleware/grpc_requestid"

	_ "micro_framework/cmd/balance"
	_ "micro_framework/cmd/order"
	_ "micro_framework/cmd/product"
	_ "micro_framework/cmd/user"
)

var version string

func main() {
	flag.Parse()
	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "show version",
		Long:  "show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", version)
		},
	}

	cmdRun := &cobra.Command{
		Use:   "run config.toml",
		Short: "run server",
		Long:  "run server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("need config file")
			}

			err := start(args[0], false)
			if err != nil {
				fmt.Printf("start grpcserver error: %+v\n", err)
			}
			return nil
		},
	}

	cmdReload := &cobra.Command{
		Use:   "reload config.toml",
		Short: "reload server",
		Long:  "reload server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("need config file")
			}

			err := start(args[0], true)
			if err != nil {
				fmt.Printf("reload grpcserver error: %+v\n", err)
			}
			return nil
		},
	}
	rootCmd := &cobra.Command{
		Use:  "grpc server",
		Long: "grpc server",
	}

	rootCmd.AddCommand(cmdRun, cmdVersion, cmdReload)
	rootCmd.Execute()
}

func start(filename string, reload bool) error {
	conf, err := configfile.GetConfig(filename)
	if err != nil {
		return err
	}
	log, err := cmd.InitGrpcLog(conf.LogPath)
	if err != nil {
		return err
	}
	//log.Infof("%+v", *conf)
	var lis net.Listener
	if reload {
		f := os.NewFile(3, "")
		lis, err = net.FileListener(f)
		if err != nil {
			return err
		}
	} else {
		lis, err = net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return err
		}
	}
	if err := registerEtcd(log, conf); err != nil {
		return err
	}
	startServer(log, lis, conf)
	return nil
}

func registerEtcd(log *logrus.Entry, conf *configfile.Config) error {
	cli, err := cmd.GetEtcdClient(conf.EtcdHost, conf.EtcdPort)
	if err != nil {
		log.Infof("GetEtcdClient err:%v", err)
		return err
	}
	// 创建命名解析
	r := &etcdnaming.GRPCResolver{Client: cli}

	// 将本服务注册添加etcd
	err = r.Update(context.TODO(), conf.ServerName, naming.Update{Op: naming.Add, Addr: fmt.Sprintf("%s:%d", conf.Host, conf.Port)})
	if err != nil {
		log.Infof("[测] update etcd err:%v", err)
		return err
	}
	return nil

}

func unRegisterEtcd(log *logrus.Entry, conf *configfile.Config) error {
	cli, err := cmd.GetEtcdClient(conf.EtcdHost, conf.EtcdPort)
	if err != nil {
		log.Infof("GetEtcdClient err:%v", err)
		return err
	}
	// 创建命名解析
	r := &etcdnaming.GRPCResolver{Client: cli}
	// 将本服务注册添加etcd中
	err = r.Update(context.TODO(), conf.ServerName, naming.Update{Op: naming.Delete, Addr: fmt.Sprintf("%s:%d", conf.Host, conf.Port)})
	if err != nil {
		log.Infof("update etcd err:%v", err)
		return err
	}
	return nil

}

func startServer(log *logrus.Entry, lis net.Listener, conf *configfile.Config) {
	collector, err := zipkin.NewHTTPCollector(conf.ZipkinHTTPEndpoint)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err: %v", err)
	}

	recorder := zipkin.NewRecorder(collector, true, fmt.Sprintf("%v:%v", conf.Host, conf.Port), "micro_framework")
	tracer, err := zipkin.NewTracer(recorder, zipkin.ClientServerSameSpan(true))
	if err != nil {
		log.Fatalf("zipkin.NewTracer err: %v", err)
	}
	opentracing.InitGlobalTracer(tracer)
	// 配置gprc中间件
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(log),
			grpc_requestid.StreamServerInterceptor(log),
			grpc_recovery.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(log),
			grpc_requestid.UnaryServerInterceptor(log),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
		)),
	)

	// 初始化并注册提供的微服务
	for _, initServer := range cmd.ServerDrivers {
		if err := initServer(grpcServer, conf); err != nil {
			log.Infof("initServer error:%v", err)
		}
	}
	// 优雅关闭以及热重启
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGTERM, syscall.SIGUSR2)
	go func() {
		for sig := range c {
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				log.Infof("Got signal: %v", sig)
				unRegisterEtcd(log, conf)
				grpcServer.GracefulStop()
				return
			case syscall.SIGUSR2:
				log.Infof("Got signal: %v", sig)
				unRegisterEtcd(log, conf)
				if err := reload(lis); err != nil {
					log.Infof("reload fail:%v", err)
					return
				}
				log.Infof("graceful reload")
				grpcServer.GracefulStop()
				return
			}
		}

	}()
	log.Infof("started")
	grpcServer.Serve(lis)
}

func reload(listener net.Listener) error {
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		return fmt.Errorf("listener is not tcp listener")
	}

	f, err := tl.File()
	if err != nil {
		return err
	}

	args := []string{"reload", os.Args[2]}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// put socket FD at the first entry
	cmd.ExtraFiles = []*os.File{f}
	return cmd.Start()
}
