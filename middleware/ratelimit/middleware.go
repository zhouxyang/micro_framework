package ratelimit

import (
	"github.com/rcrowley/go-metrics"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Limiter defines the interface to perform request rate limiting.
// If Limit function return true, the request will be rejected.
// Otherwise, the request will pass.
type Limiter struct {
	Meter   metrics.Meter
	Counter metrics.Counter
	Limit   int64
}

// UnaryServerInterceptor returns a new unary server interceptors that performs request rate limiting.
func UnaryServerInterceptor(limiter *Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		limiter.Meter.Mark(1)
		limiter.Counter.Inc(1)
		defer limiter.Counter.Dec(1)
		if limiter.Counter.Count() > limiter.Limit {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleare, please retry later.", info.FullMethod)
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func StreamServerInterceptor(limiter *Limiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		limiter.Meter.Mark(1)
		limiter.Counter.Inc(1)
		defer limiter.Counter.Dec(1)
		if limiter.Counter.Count() > limiter.Limit {
			return status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleare, please retry later.", info.FullMethod)
		}
		return handler(srv, stream)
	}
}
