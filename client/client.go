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
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, gerr := grpc.Dial(*serverAddr, opts...)
	if gerr != nil {
		log.Fatalf("fail to dial: %v", gerr)
	}
	defer conn.Close()
	client := order.NewOrderClient(conn)

	// Looking for a valid feature
	req := &order.OrderCreateRequest{
		UserID:    "123",
		ProductID: []string{"1", "2", "3"},
	}
	resp, err := client.CreateOrder(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp.OrderMsg)
	fmt.Printf("%+v\n", resp.Result)
	fmt.Printf("%+v\n", resp.Balance)
	// Feature missing.
	//printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})
	// Looking for features between 40, -75 and 42, -73.
	/*printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})*/
	// RecordRoute
	//runRecordRoute(client)

	// RouteChat
	//runRouteChat(client)
}
