package main

import (
	//"crypto/tls"
	"flag"
	"fmt"
	"log"
	//	"github.com/coreos/etcd/clientv3"
	//	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	//"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/metadata"
	"route_guide/pb/order"
)

/*
返回码对照表：
1000, ok
5000, productid is invalid
5001, userid is invalid
5002, internal error
5003, db error
5004, deduct error
*/

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	/*
		//此处是通过etcd服务发现连server
		cli, cerr := clientv3.NewFromURL("http://localhost:2379")
		if cerr != nil {
			log.Fatalf("NewFromURL err: %v", cerr)
		}
		r := &etcdnaming.GRPCResolver{Client: cli}
		b := grpc.RoundRobin(r)
	*/
	log.Println(*serverAddr)
	var opts []grpc.DialOption
	/*
		//此处是跳过证书验证的连接server，即insecure
		var tlsConf tls.Config
		tlsConf.InsecureSkipVerify = true
		creds := credentials.NewTLS(&tlsConf)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	*/
	/*
		// 此处是使用客户端中间件
		conn, gerr := grpc.Dial(*serverAddr, grpc.WithInsecure(),grpc.WithBlock(), grpc.WithUnaryInterceptor(UnaryClientInterceptor),
			grpc.WithStreamInterceptor(StreamClientInterceptor))
	*/

	/*
		//此处是明文传输 即plaintext
		opts = append(opts, grpc.WithInsecure())
	*/
	creds, err := credentials.NewClientTLSFromFile("./script/server.crt", "")
	if err != nil {
		log.Fatalf("fail to new: %v", err)
		return
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, gerr := grpc.Dial(*serverAddr, opts...)
	if gerr != nil {
		log.Fatalf("fail to dial: %v", gerr)
		return
	}
	defer conn.Close()
	client := order.NewOrderClient(conn)

	// Looking for a valid feature
	createReq := &order.OrderCreateRequest{
		UserID:    "123",
		ProductID: []string{"1", "2", "3"},
	}
	createResp, err := client.CreateOrder(context.Background(), createReq)
	if err != nil {
		log.Fatalf("fail to CreateOrder:%v", err)
		return
	}
	fmt.Printf("%+v\n", createResp.OrderMsg)
	fmt.Printf("%+v\n", createResp.Result)
	fmt.Printf("%+v\n", createResp.Balance)

	queryReq := &order.OrderQueryRequest{
		OrderID: createResp.OrderMsg.OrderID,
	}
	queryResp, err := client.QueryOrder(context.Background(), queryReq)
	if err != nil {
		log.Fatalf("fail to QueryOrder:%v", err)
		return
	}

	fmt.Printf("%+v\n", queryResp.OrderMsg)
	fmt.Printf("%+v\n", queryResp.Result)
	fmt.Printf("%+v\n", queryResp.Balance)
}
