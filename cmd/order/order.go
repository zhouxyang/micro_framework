package order

import (
	"context"
	sql "github.com/go-sql-driver/mysql"
	"route_guide/cmd"
	"route_guide/configfile"
	"route_guide/db"
	"route_guide/db/mysql"
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"

	"route_guide/pb/balance"
	pb "route_guide/pb/order"
	"route_guide/pb/product"
	"route_guide/pb/user"
)

func init() {
	//注册服务初始化函数
	cmd.RegisterService("OrderService", InitServer)

}

// InitServer 初始化MyService服务
func InitServer(grpcServer *grpc.Server, config *configfile.Config) error {
	myDB, err := configfile.InitDB(config.OrderService.OrderDB)
	if err != nil {
		return err
	}
	srv := &OrderServer{
		config:   config,
		OrderDao: &mysql.OrderDao{MyDB: myDB},
	}
	pb.RegisterOrderServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	//reflection.Register(grpcServer)
	return nil
}

// OrderServer 自定义服务结构体
type OrderServer struct {
	*mysql.OrderDao
	config *configfile.Config
}

// CreateOrder
func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.OrderCreateRequest) (*pb.OrderResponse, error) {
	log := cmd.GetLog(ctx)
	u1 := uuid.NewV4().String()
	response := &pb.OrderResponse{
		OrderMsg: &pb.OrderMsg{
			OrderID:    u1,
			ProductID:  make([]string, 0, len(req.ProductID)),
			CreateTime: time.Now().String(),
		},
		Result: &pb.Result{
			Code: 1000,
			Msg:  "OK",
		},
	}
	userID := req.UserID
	// 1. check user
	conn, err := cmd.GetGrpcClientConn(ctx, s.config, s.config.OrderService.UserService, log)
	if err != nil {
		log.Infof("GetEtcdClient error:%v", err)
		return &pb.OrderResponse{
			Result: &pb.Result{
				Code: 5002,
				Msg:  "internal error",
			},
		}, err
	}
	defer conn.Close()
	userClient := user.NewUserClient(conn.GrpcConn)
	userParam := &user.GetUserRequest{UserID: req.UserID}
	res, err := userClient.GetUser(ctx, userParam)
	if err != nil {
		log.Infof("userClient.GetUser error:%v", err)
		return &pb.OrderResponse{
			Result: &pb.Result{
				Code: 5001,
				Msg:  "userid is invalid",
			},
		}, err
	}
	if res.Result.Code != 1000 {
		log.Infof("userClient.GetUser return code:%v, msg:%v", res.Result.Code, res.Result.Msg)
		return &pb.OrderResponse{
			Result: &pb.Result{
				Code: 5001,
				Msg:  "userid is invalid",
			},
		}, err
	}
	log.Infof("get userid:%v, username: %v", res.UserID, res.UserName)
	// 2. check product
	conn, err = cmd.GetGrpcClientConn(ctx, s.config, s.config.OrderService.ProductService, log)
	if err != nil {
		log.Infof("GetEtcdClient error:%v", err)
		return &pb.OrderResponse{Result: &pb.Result{Code: 5002, Msg: "internal error"}}, err
	}
	defer conn.Close()

	productIDMap := make(map[string]string)
	for _, productID := range req.ProductID {
		productClient := product.NewProductClient(conn.GrpcConn)
		productParam := &product.GetProductRequest{ProductID: productID}
		res, err := productClient.GetProduct(ctx, productParam)
		if err != nil {
			log.Infof("product.GetProduct error:%v", err)
			return &pb.OrderResponse{Result: &pb.Result{Code: 5000, Msg: "productid is invalid"}}, err
		}
		if res.Result.Code != 1000 {
			log.Infof("productClient.GetProduct return code:%v, msg:%v", res.Result.Code, res.Result.Msg)
			continue
		}
		productIDMap[res.ProductID] = res.ProductPrice
	}

	// 3. deduct
	conn, err = cmd.GetGrpcClientConn(ctx, s.config, s.config.OrderService.BalanceService, log)
	if err != nil {
		log.Infof("GetEtcdClient error:%v", err)
		return &pb.OrderResponse{Result: &pb.Result{Code: 5002, Msg: "internal error"}}, err
	}
	defer conn.Close()
	balanceClient := balance.NewBalanceClient(conn.GrpcConn)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	stream, err := balanceClient.Deduct(ctx)
	if err != nil {
		log.Infof("balanceClient.Deduct error: %v", err)
		return &pb.OrderResponse{Result: &pb.Result{Code: 5004, Msg: "deduct error"}}, err
	}
	total := decimal.Zero
	for id, price := range productIDMap {
		priceDecimal, err := decimal.NewFromString(price)
		if err != nil {
			log.Infof("decimal.NewFromString error:%v", err)
			continue
		}
		deduct := &balance.DeductRequest{
			UserID:    userID,
			ProductID: id,
			Price:     price,
		}
		if err := stream.Send(deduct); err != nil {
			log.Infof("Failed to send a deduct: %v", err)
			continue
		}

		// 入库order表
		o := db.Order{
			ProductID:  id,
			OrderPrice: priceDecimal,
			OrderID:    u1,
			CreateTime: sql.NullTime{Time: time.Now(), Valid: true},
		}
		if err := s.OrderDao.InsertOrder(log, &o); err != nil {
			log.Infof("OrderDao.InsertOrder error:%v", err)
			continue
		}
		response.OrderMsg.ProductID = append(response.OrderMsg.ProductID, id)
		total = total.Add(priceDecimal)
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Infof("stream.CloseAndRecv error", err)
		return &pb.OrderResponse{Result: &pb.Result{Code: 5004, Msg: "deduct error"}}, err
	}
	if reply.Result.Code != 1000 {
		log.Infof("reply.Result.Code return code:%v, msg:%v", reply.Result.Code, reply.Result.Msg)
		return &pb.OrderResponse{Result: &pb.Result{Code: 5004, Msg: "deduct error"}}, nil
	}
	response.OrderMsg.TotolMoney = total.String()
	response.Balance = reply.Balance
	return response, nil
}

func (s *OrderServer) QueryOrder(ctx context.Context, req *pb.OrderQueryRequest) (*pb.OrderResponse, error) {
	log := cmd.GetLog(ctx)

	response := &pb.OrderResponse{
		OrderMsg: &pb.OrderMsg{
			OrderID:   req.OrderID,
			ProductID: make([]string, 0, 8),
		},
		Result: &pb.Result{
			Code: 1000,
			Msg:  "OK",
		},
	}
	orders, err := s.OrderDao.GetOrderByOrderID(log, req.OrderID)
	if err != nil {
		log.Infof("s.QueryTable error:%v", err)
		return &pb.OrderResponse{Result: &pb.Result{Code: 5004, Msg: "deduct error"}}, err
	}
	total := decimal.Zero
	for _, o := range orders {
		total = total.Add(o.OrderPrice)
		response.OrderMsg.CreateTime = o.CreateTime.Time.String()
		response.OrderMsg.ProductID = append(response.OrderMsg.ProductID, o.ProductID)
	}
	return response, nil
}
