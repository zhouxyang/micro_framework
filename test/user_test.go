package test

import (
	"context"
	"micro_framework/db"
	"micro_framework/pb/user"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func initUserDB() (*gorm.DB, error) {
	myDB, err := initDB()
	if err != nil {
		return nil, errors.Errorf("%v", err)
	}
	b := db.User{
		UserID: UserID,
		Name:   "zy",
	}
	if err := myDB.Save(&b).Error; err != nil {
		return nil, errors.Errorf("s.Update error:%v", err)
	}
	return myDB, nil
}

func cleanUserDB(myDB *gorm.DB) error {
	if err := myDB.Delete(db.User{}, "userid=?", UserID).Error; err != nil {
		return errors.Errorf("s.Update error:%v", err)
	}
	if err := myDB.Close(); err != nil {
		return errors.Errorf("db close error:%v", err)
	}
	return nil
}

func TestUser(t *testing.T) {
	myDB, err := initUserDB()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer cleanUserDB(myDB)
	conn, err := getConn()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer conn.Close()
	userClient := user.NewUserClient(conn)
	userParam := &user.GetUserRequest{UserID: UserID}
	userRes, err := userClient.GetUser(context.Background(), userParam)
	if err != nil {
		t.Errorf("userClient.GetUser error: %v", err)
		return
	}
	if userRes.Result.Code != 1000 {
		t.Errorf("reply.Result.Code return code:%v, msg:%v", userRes.Result.Code, userRes.Result.Msg)
		return
	}
	if userRes.UserName != "zy" {
		t.Errorf("get username:%v from:%v", userRes.UserName, UserID)
		return
	}
}
