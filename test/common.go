package test

import (
	"flag"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//go test -v -count=1  test/*

var (
	ServerAddr = flag.String("server_addr", "micro-framework.mydomain.com:443", "The server address in the format of host:port")
	Crt        = flag.String("crt", "../script/ingress/server.crt", "file path of server.crt")
	Database   = flag.String("database", "root:test@tcp(localhost:32306)/grpc?charset=utf8&parseTime=true&loc=Local&timeout=3s", "db config")
)

//UserID create new userID
var (
	UserID     = uuid.NewV4().String()
	ProductID  = uuid.NewV4().String()
	ProductID2 = uuid.NewV4().String()
	ProductID3 = uuid.NewV4().String()
)

func initDB() (*gorm.DB, error) {
	myDB, err := gorm.Open("mysql", *Database)
	if err != nil {
		return nil, errors.Errorf("%v", err)
	}
	return myDB, nil

}

func getConn() (*grpc.ClientConn, error) {
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile(*Crt, "")
	if err != nil {
		return nil, errors.Errorf("fail to new: %v", err)
	}
	entry := logrus.New()
	entry.Formatter = &logrus.JSONFormatter{}
	entry.Out = os.Stdout

	log := entry.WithFields(logrus.Fields{
		"pid": os.Getpid(),
	})
	conn, gerr := grpc.Dial(*ServerAddr,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_logrus.UnaryClientInterceptor(log),
			),
		),
	)
	if gerr != nil {
		return nil, errors.Errorf("fail to dial: %v", gerr)
	}
	return conn, nil
}
