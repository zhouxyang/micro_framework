package product

import (
	"context"
	"fmt"
	"micro_framework/cmd"
	"micro_framework/configfile"
	"micro_framework/db"

	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"

	pb "micro_framework/pb/product"
)

func init() {
	//注册服务初始化函数
	cmd.RegisterService("ProductService", InitServer)

}

// InitServer 初始化MyService服务
func InitServer(grpcServer *grpc.Server, config *configfile.Config) error {
	myDB, err := db.InitDB(config.ProductService.ProductDB)
	if err != nil {
		return err
	}
	myDB.AutoMigrate(&db.Product{})
	myDB.LogMode(true)
	srv := &ProductServer{
		config:     config,
		ProductDao: myDB,
	}
	pb.RegisterProductServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	//	reflection.Register(grpcServer)
	//熔断器设置
	hystrix.ConfigureCommand("ProductService", hystrix.CommandConfig{
		Timeout:                1 * 1000, //超时
		MaxConcurrentRequests:  10,       //最大并发的请求数
		RequestVolumeThreshold: 5,        //请求量阈值
		SleepWindow:            3 * 1000, //熔断开启多久尝试发起一次请求
		ErrorPercentThreshold:  1,        //误差阈值百分比
	})
	return nil
}

// ProductServer 自定义服务结构体
type ProductServer struct {
	ProductDao *db.MyDB
	config     *configfile.Config
}

//GetProduct 校验用户
func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	log := cmd.GetLog(ctx)
	s.ProductDao.SetLogger(log)
	var product db.Product
	if err := s.ProductDao.Find(ctx, &product, "productid=?", req.ProductID).Error; err != nil {
		return &pb.ProductResponse{Result: &pb.Result{Code: 5000, Msg: "productid is invalid"}}, err
	}
	log.Infof("get product:%+v from db", product)
	return &pb.ProductResponse{
		ProductID:    product.ProductID,
		ProductName:  product.ProductName,
		ProductPrice: fmt.Sprintf("%v", product.ProductPrice),
		Result:       &pb.Result{Code: 1000, Msg: "ok"},
	}, nil
}
