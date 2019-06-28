package test

import (
	"flag"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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
	var opts []grpc.DialOption
	creds, err := credentials.NewClientTLSFromFile(*Crt, "")
	if err != nil {
		return nil, errors.Errorf("fail to new: %v", err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, gerr := grpc.Dial(*ServerAddr, opts...)
	if gerr != nil {
		return nil, errors.Errorf("fail to dial: %v", gerr)
	}
	return conn, nil
}
