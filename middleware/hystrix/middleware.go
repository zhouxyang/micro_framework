package hystrix

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor 通过hystrix进行熔断
func UnaryClientInterceptor(commandName string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log := ctxlogrus.Extract(ctx)
		err := hystrix.Do(commandName,
			func() (err error) {
				return invoker(ctx, method, req, reply, cc, opts...)
			},
			func(err error) error {
				log.Infof("request is blocked")
				return err
			})

		if err != nil {
			return err
		}
		return nil
	}
}

// StreamClientInterceptor 通过hystrix进行熔断
func StreamClientInterceptor(log *logrus.Entry, commandName string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var clientStream grpc.ClientStream
		err := hystrix.Do(commandName,
			func() (err error) {
				clientStream, err = streamer(ctx, desc, cc, method, opts...)
				return err
			},
			func(err error) error {
				log.Infof("request is blocked")
				return err
			})
		if err != nil {
			return nil, err
		}
		return clientStream, err
	}
}
