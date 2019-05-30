package balance

import (
	"io"
	"route_guide/cmd"
	"route_guide/configfile"
	"route_guide/db/mysql"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"

	pb "route_guide/pb/balance"
)

func init() {
	//注册服务初始化函数
	cmd.RegisterService("BalanceService", InitServer)

}

// InitServer 初始化MyService服务
func InitServer(grpcServer *grpc.Server, config *configfile.Config) error {
	myDB, err := configfile.InitDB(config.BalanceService.BalanceDB)
	if err != nil {
		return err
	}
	srv := &BalanceServer{
		config:     config,
		BalanceDao: &mysql.BalanceDao{MyDB: myDB},
	}
	pb.RegisterBalanceServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	//reflection.Register(grpcServer)
	return nil
}

// BalanceServer 自定义服务结构体
type BalanceServer struct {
	*mysql.BalanceDao
	config *configfile.Config
}

//Deduct 客户端流方式
func (s *BalanceServer) Deduct(stream pb.Balance_DeductServer) error {
	balance := decimal.Zero
	log := cmd.GetLog(stream.Context())
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
		b, err := s.BalanceDao.GetBalanceByUserID(log, req.UserID)
		if err != nil {
			log.Infof("s.BalanceDao.GetBalanceByUserID:%v", err)
			continue
		}
		decimalPrice, err := decimal.NewFromString(req.Price)
		if err != nil {
			log.Infof("decimal.NewFromString error:%v", err)
			continue
		}
		//2. 进行余额扣除
		b.Balance = b.Balance.Sub(decimalPrice)
		if err := s.BalanceDao.UpdateBalanceByUserID(log, req.UserID, b.Balance); err != nil {
			log.Infof("s.Update error:%v", err)
			continue
		}
		balance = b.Balance
	}
}
