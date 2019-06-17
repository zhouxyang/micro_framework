package balance

import (
	"io"
	"micro_framework/cmd"
	"micro_framework/configfile"
	"micro_framework/db"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"

	pb "micro_framework/pb/balance"
)

func init() {
	//注册服务初始化函数
	cmd.RegisterService("BalanceService", InitServer)

}

// InitServer 初始化MyService服务
func InitServer(grpcServer *grpc.Server, config *configfile.Config) error {
	myDB, err := gorm.Open("mysql", config.BalanceService.BalanceDB)
	if err != nil {
		return err
	}
	myDB.AutoMigrate(&db.Balance{})
	myDB.LogMode(true)
	srv := &BalanceServer{
		config:     config,
		BalanceDao: myDB,
	}
	pb.RegisterBalanceServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	//reflection.Register(grpcServer)
	return nil
}

// BalanceServer 自定义服务结构体
type BalanceServer struct {
	BalanceDao *gorm.DB
	config     *configfile.Config
}

//Deduct 客户端流方式
func (s *BalanceServer) Deduct(stream pb.Balance_DeductServer) error {
	balance := decimal.Zero
	log := cmd.GetLog(stream.Context())
	s.BalanceDao.SetLogger(log)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.DeductResponse{
				Balance: balance.String(),
				Result: &pb.Result{
					Code: 1000,
					Msg:  "ok",
				},
			})
		}
		if err != nil {
			return stream.SendAndClose(&pb.DeductResponse{
				Result: &pb.Result{
					Code: 5002,
					Msg:  err.Error(),
				},
			})
		}
		//1. 先查询余额
		b := db.Balance{}
		if err := s.BalanceDao.Find(&b, "userid=?", req.UserID).Error; err != nil {
			log.Infof("s.BalanceDao.GetBalanceByUserID:%v", err)
			continue
		}
		decimalPrice, aerr := decimal.NewFromString(req.Price)
		if aerr != nil {
			log.Infof("decimal.NewFromString error:%v", aerr)
			continue
		}
		//2. 进行余额扣除
		decimalBalance := decimal.NewFromFloat(b.Balance)
		decimalBalance = decimalBalance.Sub(decimalPrice)
		if decimalBalance.Cmp(decimal.Zero) < 0 {
			return stream.SendAndClose(&pb.DeductResponse{
				Result: &pb.Result{
					Code: 5005,
					Msg:  "balance is not enough",
				},
			})

		}
		b.Balance, _ = decimalBalance.Float64()
		if err := s.BalanceDao.Save(&b).Error; err != nil {
			log.Infof("s.Update error:%v", err)
			continue
		}
		balance = decimalBalance
	}
}
