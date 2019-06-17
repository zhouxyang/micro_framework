package order

import (
	"context"
	sql "github.com/go-sql-driver/mysql"
	"micro_framework/cmd"
	"micro_framework/configfile"
	"micro_framework/db"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"

	"micro_framework/pb/balance"
	pb "micro_framework/pb/order"
	"micro_framework/pb/product"
	"micro_framework/pb/user"
)

func init() {
	//注册服务初始化函数
	cmd.RegisterService("OrderService", InitServer)

}

// InitServer 初始化MyService服务
func InitServer(grpcServer *grpc.Server, config *configfile.Config) error {
	myDB, err := gorm.Open("mysql", config.OrderService.OrderDB)
	if err != nil {
		return err
	}
	myDB.AutoMigrate(&db.Product{})
	myDB.LogMode(true)
	srv := &OrderServer{
		config:   config,
		OrderDao: myDB,
	}
	pb.RegisterOrderServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	//reflection.Register(grpcServer)
	return nil
}

// OrderServer 自定义服务结构体
type OrderServer struct {
	OrderDao *gorm.DB
	config   *configfile.Config
}

// CreateOrder
func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.OrderCreateRequest) (*pb.OrderResponse, error) {
	log := cmd.GetLog(ctx)
	s.OrderDao.SetLogger(log)
	u1 := uuid.NewV4().String()
	orderResponse := &pb.OrderResponse{
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
	userRes := &user.UserResponse{}
	userRes, err = userClient.GetUser(ctx, userParam)
	if err != nil {
		log.Infof("userClient.GetUser error:%v", err)
		return &pb.OrderResponse{
			Result: &pb.Result{
				Code: 5001,
				Msg:  "userid is invalid",
			},
		}, err
	}
	if userRes.Result.Code != 1000 {
		log.Infof("userClient.GetUser return code:%v, msg:%v", userRes.Result.Code, userRes.Result.Msg)
		return &pb.OrderResponse{
			Result: &pb.Result{
				Code: 5001,
				Msg:  "userid is invalid",
			},
		}, err
	}
	log.Infof("get userid:%v, username: %v", userRes.UserID, userRes.UserName)
	// 2. check product async
	conn, err = cmd.GetGrpcClientConn(ctx, s.config, s.config.OrderService.ProductService, log)
	if err != nil {
		log.Infof("GetEtcdClient error:%v", err)
		return &pb.OrderResponse{Result: &pb.Result{Code: 5002, Msg: "internal error"}}, err
	}
	defer conn.Close()

	productClient := product.NewProductClient(conn.GrpcConn)
	prdChan := make(chan *db.Product, len(req.ProductID))
	for _, prdID := range req.ProductID {
		go func(productID string) {
			prd := &db.Product{}
			defer func() {
				prdChan <- prd
			}()
			productParam := &product.GetProductRequest{ProductID: productID}
			productRes, err := productClient.GetProduct(ctx, productParam)
			if err != nil {
				log.Infof("product.GetProduct error:%v", err)
				return
			}
			if productRes.Result.Code != 1000 {
				log.Infof("productClient.GetProduct return code:%v, msg:%v", productRes.Result.Code, productRes.Result.Msg)
				return
			}

			priceDecimal, err := decimal.NewFromString(productRes.ProductPrice)
			if err != nil {
				log.Infof("decimal.NewFromString error:%v", err)
				return
			}
			prd.ProductPrice = priceDecimal
			prd.ProductID = productRes.ProductID
		}(prdID)
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
	for i := 0; i < len(req.ProductID); i++ {
		prd := <-prdChan
		id, priceDecimal := prd.ProductID, prd.ProductPrice
		if id == "" {
			continue
		}
		deduct := &balance.DeductRequest{
			UserID:    userID,
			ProductID: id,
			Price:     priceDecimal.String(),
		}
		if err := stream.Send(deduct); err != nil {
			log.Infof("Failed to send a deduct: %v", err)
			continue
		}

		// 入库order表
		o := db.Order{
			ProductID:  id,
			OrderID:    u1,
			CreateTime: sql.NullTime{Time: time.Now(), Valid: true},
		}
		o.ProductPrice, _ = priceDecimal.Float64()
		if err := s.OrderDao.Create(&o).Error; err != nil {
			log.Infof("Failed to insert order: %v", err)
			continue
		}
		orderResponse.OrderMsg.ProductID = append(orderResponse.OrderMsg.ProductID, id)
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
	orderResponse.OrderMsg.TotolMoney = total.String()
	orderResponse.Balance = reply.Balance
	return orderResponse, nil
}

func (s *OrderServer) QueryOrder(ctx context.Context, req *pb.OrderQueryRequest) (*pb.OrderResponse, error) {
	log := cmd.GetLog(ctx)
	s.OrderDao.SetLogger(log)
	orderResponse := &pb.OrderResponse{
		OrderMsg: &pb.OrderMsg{
			OrderID:   req.OrderID,
			ProductID: make([]string, 0, 8),
		},
		Result: &pb.Result{
			Code: 1000,
			Msg:  "OK",
		},
	}

	orders := make([]db.Order, 0, 8)
	if err := s.OrderDao.Find(&orders, "orderid=?", req.OrderID).Error; err != nil {
		return &pb.OrderResponse{Result: &pb.Result{Code: 5004, Msg: "query order error"}}, err
	}
	total := decimal.Zero
	for _, o := range orders {
		total = total.Add(decimal.NewFromFloat(o.ProductPrice))
		orderResponse.OrderMsg.CreateTime = o.CreateTime.Time.String()
		orderResponse.OrderMsg.ProductID = append(orderResponse.OrderMsg.ProductID, o.ProductID)

	}
	orderResponse.OrderMsg.TotolMoney = total.String()
	return orderResponse, nil
}
