package product

import (
	"context"
	"micro_framework/cmd"
	"micro_framework/configfile"
	"micro_framework/db/mysql"

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
	myDB, err := configfile.InitDB(config.ProductService.ProductDB)
	if err != nil {
		return err
	}
	srv := &ProductServer{
		config:     config,
		ProductDao: &mysql.ProductDao{MyDB: myDB},
	}
	pb.RegisterProductServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	//	reflection.Register(grpcServer)
	return nil
}

// ProductServer 自定义服务结构体
type ProductServer struct {
	*mysql.ProductDao
	config *configfile.Config
}

//GetProduct 校验用户
func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	log := cmd.GetLog(ctx)
	product, err := s.ProductDao.GetProductByProductID(log, req.ProductID)
	if err != nil {
		return &pb.ProductResponse{Result: &pb.Result{Code: 5003, Msg: "db error"}}, err
	}
	if product == nil {
		return &pb.ProductResponse{Result: &pb.Result{Code: 5000, Msg: "productid is invalid"}}, err
	}
	log.Infof("get product:%+v from db", product)
	return &pb.ProductResponse{
		ProductID:    product.ProductID,
		ProductName:  product.ProductName,
		ProductPrice: product.ProductPrice.String(),
		Result:       &pb.Result{Code: 1000, Msg: "ok"},
	}, nil
}
