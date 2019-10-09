package user

import (
	"context"
	"micro_framework/cmd"
	"micro_framework/configfile"
	//"micro_framework/db/mysql"
	"micro_framework/db"

	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"

	pb "micro_framework/pb/user"
)

func init() {
	//注册服务初始化函数
	cmd.RegisterService("UserService", InitServer)
}

// InitServer 初始化MyService服务
func InitServer(grpcServer *grpc.Server, config *configfile.Config) error {
	myDB, err := db.InitDB(config.UserService.UserDB)
	if err != nil {
		return err
	}
	myDB.AutoMigrate(&db.User{})
	myDB.LogMode(true)
	srv := &UserServer{
		config:  config,
		UserDao: myDB,
	}
	pb.RegisterUserServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	//	reflection.Register(grpcServer)
	//熔断器设置
	hystrix.ConfigureCommand("UserService", hystrix.CommandConfig{
		Timeout:                1 * 1000, //超时
		MaxConcurrentRequests:  5,        //最大并发的请求数
		RequestVolumeThreshold: 5,        //请求量阈值
		SleepWindow:            3 * 1000, //熔断开启多久尝试发起一次请求
		ErrorPercentThreshold:  1,        //误差阈值百分比
	})
	return nil
}

// UserServer 自定义服务结构体
type UserServer struct {
	config  *configfile.Config
	UserDao *db.MyDB
}

//GetUser 校验用户
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	log := cmd.GetLog(ctx)
	s.UserDao.SetLogger(log)
	var user db.User
	if err := s.UserDao.Find(ctx, &user, "userid=?", req.UserID).Error; err != nil {
		return &pb.UserResponse{Result: &pb.Result{Code: 5001, Msg: "userid not found"}}, err
	}
	log.Infof("get user:%+v from db", req)
	return &pb.UserResponse{
		UserID:   user.UserID,
		UserName: user.Name,
		Result:   &pb.Result{Code: 1000, Msg: "ok"},
	}, nil
}
