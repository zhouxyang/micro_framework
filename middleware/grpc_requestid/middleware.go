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

// UnaryServerInterceptor 从客户端拿到matadata中的requestid
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := ctxlogrus.Extract(ctx)
		u1 := ""
		// Read metadata from client.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Infof("UnaryEcho: failed to get metadata")
		}
		if t, ok := md["requestid"]; ok {
			//log.Infof("requestid from metadata:%v", t)
			for _, e := range t {
				u1 = e
			}
		}
		if u1 == "" {
			u1 = uuid.NewV4().String()
		}
		ctxlogrus.AddFields(ctx, logrus.Fields{
			"requestid": u1,
		})
		return handler(ctx, req)
	}
}

// StreamServerInterceptor 从客户端拿到matadata中的requestid
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
		if t, ok := md["requestid"]; ok {
			for _, e := range t {
				u1 = e
			}
		}
		if u1 == "" {
			u1 = uuid.NewV4().String()
		}
		ctxlogrus.AddFields(ctx, logrus.Fields{
			"requestid": u1,
		})
		wrap := grpc_middleware.WrapServerStream(stream)
		wrap.WrappedContext = ctx
		return handler(srv, wrap)
	}
}

// UnaryClientInterceptor 把客户端的requestid传到服务端
func UnaryClientInterceptor(log *logrus.Entry, opts ...grpc_logrus.Option) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log := ctxlogrus.Extract(ctx)
		requestid, ok := log.Data["requestid"].(string)
		if !ok {
			log.Infof("get requestid is nil")
			requestid = uuid.NewV4().String()
		}
		// Read metadata from client.
		md := metadata.Pairs("requestid", requestid)
		newCtx := metadata.NewOutgoingContext(ctx, md)
		//log.Infof("before invoker. method: %+v, request:%+v, requestid:%v", method, req, requestid)
		err := invoker(newCtx, method, req, reply, cc, opts...)
		//log.Infof("after invoker. reply: %+v", reply)
		return err
	}
}

// StreamClientInterceptor 把客户端requestid传到服务端
func StreamClientInterceptor(log *logrus.Entry, opts ...grpc_logrus.Option) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requestid, ok := log.Data["requestid"].(string)
		if !ok {
			requestid = uuid.NewV4().String()
		}
		// Read metadata from client.
		md := metadata.Pairs("requestid", requestid)
		newCtx := metadata.NewOutgoingContext(ctx, md)
		//log.Infof("before invoker. method: %+v, StreamDesc:%+v, requestid:%v", method, desc, requestid)
		clientStream, err := streamer(newCtx, desc, cc, method, opts...)
		//log.Infof("before invoker. method: %+v", method)
		return clientStream, err
	}
}
