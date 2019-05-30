package user

import (
	"context"
	"route_guide/cmd"
	"route_guide/configfile"
	"route_guide/db/mysql"

	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"

	pb "route_guide/pb/user"
)

func init() {
	//注册服务初始化函数
	cmd.RegisterService("UserService", InitServer)

}

// InitServer 初始化MyService服务
func InitServer(grpcServer *grpc.Server, config *configfile.Config) error {
	myDB, err := configfile.InitDB(config.UserService.UserDB)
	if err != nil {
		return err
	}
	//orm.RunSyncdb("default", false, true)
	srv := &UserServer{
		config:  config,
		UserDao: &mysql.UserDao{MyDB: myDB},
	}
	pb.RegisterUserServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	//	reflection.Register(grpcServer)
	return nil
}

// UserServer 自定义服务结构体
type UserServer struct {
	*mysql.UserDao
	config *configfile.Config
}

//GetUser 校验用户
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	log := cmd.GetLog(ctx)
	user, err := s.UserDao.GetUserByUserID(log, req.UserID)
	if err != nil || user == nil {
		return &pb.UserResponse{Result: &pb.Result{Code: 5001, Msg: "userid not found"}}, err
	}
	log.Infof("get user:%+v from db", req)
	return &pb.UserResponse{
		UserID:   user.UserID,
		UserName: user.Name,
		Result:   &pb.Result{Code: 1000, Msg: "ok"},
	}, nil
}
