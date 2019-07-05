package grpcrequestid

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	proto "github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
)

func init() {
	marshaler = jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		OrigName:     true,
	}
}

var marshaler jsonpb.Marshaler

type protoMessage struct {
	msg proto.Message
}

func (m protoMessage) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := marshaler.Marshal(&buf, m.msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnaryServerInterceptor 打印请求与响应内容
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logger := ctxlogrus.Extract(ctx)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Infof("UnaryServerInterceptor: failed to get metadata")
			return
		}

		var addr string
		if pr, ok := peer.FromContext(ctx); ok {
			if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
				addr = tcpAddr.IP.String()
			} else {
				addr = pr.Addr.String()
			}
		}
		reqMsg, ok := req.(proto.Message)
		if !ok {
			logger.Infof("request dump:%v", fmt.Sprintf("not proto.Message: %v", req))
			return
		}
		resp, err = handler(ctx, req)
		respMsg, ok := resp.(proto.Message)
		if !ok {
			logger.Infof("request dump:%v", fmt.Sprintf("not proto.Message: %v", resp))
			return
		}

		log := logger.WithFields(logrus.Fields{
			"header": md,
			"addr":   addr,
			"req":    protoMessage{msg: reqMsg},
			"resp":   protoMessage{msg: respMsg},
		})

		log.Infof("request dump")
		return
	}
}
