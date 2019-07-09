package requestdump

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor 从客户端拿到matadata中的traceid
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := ctxlogrus.Extract(ctx)
		u1 := ""
		// Read metadata from client.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Infof("UnaryEcho: failed to get metadata")
		}
		if t, ok := md["traceid"]; ok {
			//log.Infof("traceid from metadata:%v", t)
			for _, e := range t {
				u1 = e
			}
		}
		if u1 == "" {
			u1 = uuid.NewV4().String()
		}
		ctxlogrus.AddFields(ctx, logrus.Fields{
			"traceid": u1,
		})
		return handler(ctx, req)
	}
}

// StreamServerInterceptor 从客户端拿到matadata中的traceid
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ctx := stream.Context()
		log := ctxlogrus.Extract(ctx)
		u1 := ""
		// Read metadata from client.
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			log.Infof("ServerStreamingEcho: failed to get metadata")
		}
		if t, ok := md["traceid"]; ok {
			for _, e := range t {
				u1 = e
			}
		}
		if u1 == "" {
			u1 = uuid.NewV4().String()
		}
		ctxlogrus.AddFields(ctx, logrus.Fields{
			"traceid": u1,
		})
		wrap := grpc_middleware.WrapServerStream(stream)
		wrap.WrappedContext = ctx
		return handler(srv, wrap)
	}
}

// UnaryClientInterceptor 把客户端的traceid传到服务端
func UnaryClientInterceptor(log *logrus.Entry, opts ...grpc_logrus.Option) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log := ctxlogrus.Extract(ctx)
		traceid, ok := log.Data["traceid"].(string)
		if !ok {
			log.Infof("get traceid is nil")
			traceid = uuid.NewV4().String()
		}
		// Read metadata from client.
		md := metadata.Pairs("traceid", traceid)
		newCtx := metadata.NewOutgoingContext(ctx, md)
		//log.Infof("before invoker. method: %+v, request:%+v, traceid:%v", method, req, traceid)
		err := invoker(newCtx, method, req, reply, cc, opts...)
		//log.Infof("after invoker. reply: %+v", reply)
		return err
	}
}

// StreamClientInterceptor 把客户端traceid传到服务端
func StreamClientInterceptor(log *logrus.Entry, opts ...grpc_logrus.Option) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		traceid, ok := log.Data["traceid"].(string)
		if !ok {
			traceid = uuid.NewV4().String()
		}
		// Read metadata from client.
		md := metadata.Pairs("traceid", traceid)
		newCtx := metadata.NewOutgoingContext(ctx, md)
		//log.Infof("before invoker. method: %+v, StreamDesc:%+v, traceid:%v", method, desc, traceid)
		clientStream, err := streamer(newCtx, desc, cc, method, opts...)
		//log.Infof("before invoker. method: %+v", method)
		return clientStream, err
	}
}
