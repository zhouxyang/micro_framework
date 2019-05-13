package grpcrequestid

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds logrus.Entry to the context.
func UnaryServerInterceptor(entry *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		u1 := uuid.NewV4().String()
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
		u1 := uuid.NewV4().String()
		ctxlogrus.AddFields(ctx, logrus.Fields{
			"requestid": u1,
		})
		wrap := grpc_middleware.WrapServerStream(stream)
		wrap.WrappedContext = ctx
		return handler(srv, wrap)
	}
}
