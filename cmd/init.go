package cmd

import (
	"context"
	"fmt"
	"micro_framework/configfile"
	"micro_framework/pkg/caller"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	//grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_requestid "micro_framework/middleware/grpc_requestid"
)

func InitGrpcLog(filename string) (*logrus.Entry, error) {
	// init grpc log
	grpcLog, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		werr := errors.Wrap(err, fmt.Sprintf("open logfile failed: %v", filename))
		return nil, werr
	}
	entry := logrus.New()
	entry.Formatter = &logrus.JSONFormatter{}
	hookCaller := caller.NewHook(&caller.CallerHookOptions{
		Flags:      caller.Llongfile,
		EnableLine: true,
		EnableFile: true,
	})
	entry.Hooks[logrus.InfoLevel] = append(entry.Hooks[logrus.InfoLevel], hookCaller)
	entry.Out = grpcLog

	log := entry.WithFields(logrus.Fields{
		"pid": os.Getpid(),
	})
	log.Infof("init grpclog success")
	return log, nil
}

// GetLog 从ctx拿到log
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

// GetGrpcConnByEtcd 创建grpc连接
func GetGrpcConnByEtcd(ctx context.Context, service string, cli *clientv3.Client, log *logrus.Entry) (*grpc.ClientConn, error) {
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
				//grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(3)),
				grpc_requestid.UnaryClientInterceptor(log),
				grpc_opentracing.UnaryClientInterceptor(),
			),
		),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				grpc_logrus.StreamClientInterceptor(log),
				//grpc_retry.StreamClientInterceptor(grpc_retry.WithMax(3)),
				grpc_requestid.StreamClientInterceptor(log),
				grpc_opentracing.StreamClientInterceptor(),
			),
		),
	)

	if err != nil {
		errors.Wrap(err, "fail to dial")
		return nil, err
	}
	return conn, err
}

// GrpcClientConn 封装的grpc连接
type GrpcClientConn struct {
	GrpcConn *grpc.ClientConn
	EtcdConn *clientv3.Client
}

// GetGrpcConnByDomain 创建grpc连接
func GetGrpcConnByDomain(ctx context.Context, service string, log *logrus.Entry) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(
		ctx,
		service,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_logrus.UnaryClientInterceptor(log),
				//grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(3)),
				grpc_requestid.UnaryClientInterceptor(log),
				grpc_opentracing.UnaryClientInterceptor(),
			),
		),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				grpc_logrus.StreamClientInterceptor(log),
				//grpc_retry.StreamClientInterceptor(grpc_retry.WithMax(3)),
				grpc_requestid.StreamClientInterceptor(log),
				grpc_opentracing.StreamClientInterceptor(),
			),
		),
	)

	if err != nil {
		errors.Wrap(err, "fail to dial")
		return nil, err
	}
	return conn, err
}

// GetGrpcClientConn 创建grpc客户端连接，1. 先判断在etcd是否找得到该服务，2. 如果服务不在etcd中，则service为外部服务域名
func GetGrpcClientConn(ctx context.Context, conf *configfile.Config, service string, log *logrus.Entry) (*GrpcClientConn, error) {
	cli, err := GetEtcdClient(conf.EtcdHost, conf.EtcdPort)
	if err != nil {
		return nil, errors.Wrap(err, "GetEtcdClient error")
	}
	resp, err := cli.Get(context.TODO(), service, clientv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "kv.get error")
	}
	//log.Infof("get service:%s return resp:%+v, count:%d", service, resp, resp.Count)
	if resp.Count > 0 {
		grpcConn, err := GetGrpcConnByEtcd(ctx, service, cli, log)
		if err != nil {
			return nil, errors.Wrap(err, "GetGrpcConnByEtcd")
		}
		return &GrpcClientConn{
			GrpcConn: grpcConn,
			EtcdConn: cli,
		}, nil

	}
	grpcConn, err := GetGrpcConnByDomain(ctx, service, log)
	if err != nil {
		return nil, errors.Wrap(err, "GetGrpcConnByDomain")
	}
	return &GrpcClientConn{
		GrpcConn: grpcConn,
		EtcdConn: cli,
	}, nil

}

// Close 关闭etcd连接，grpc连接
func (c *GrpcClientConn) Close() {
	c.EtcdConn.Close()
	c.GrpcConn.Close()
}
