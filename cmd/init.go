package cmd

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"os"
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
