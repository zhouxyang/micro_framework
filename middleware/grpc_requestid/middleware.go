package grpcrequestid

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"route_guide/cmd"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds logrus.Entry to the context.
func UnaryServerInterceptor(entry *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := cmd.GetLog(ctx)
		u1 := ""
		// Read metadata from client.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Infof("UnaryEcho: failed to get metadata")
		}
		if t, ok := md["requestid"]; ok {
			log.Infof("requestid from metadata:%v", t)
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

// StreamServerInterceptor returns grpc.StreamServerInterceptor
func StreamServerInterceptor(entry *logrus.Entry) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ctx := stream.Context()
		log := cmd.GetLog(ctx)
		u1 := ""
		// Read metadata from client.
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			log.Infof("ServerStreamingEcho: failed to get metadata")
		}
		if t, ok := md["requestid"]; ok {
			log.Infof("request from metadata: %v", t)
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
