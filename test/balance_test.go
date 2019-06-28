package test

import (
	"context"
	"micro_framework/db"
	"micro_framework/pb/balance"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func initBalanceDB() (*gorm.DB, error) {
	myDB, err := initDB()
	if err != nil {
		return nil, errors.Errorf("%v", err)
	}
	b := db.Balance{
		UserID:  UserID,
		Balance: 20,
	}
	if err := myDB.Save(&b).Error; err != nil {
		return nil, errors.Errorf("s.Update error:%v", err)
	}
	return myDB, nil
}

func cleanBalanceDB(myDB *gorm.DB) error {
	if err := myDB.Delete(db.Balance{}, "userid=?", UserID).Error; err != nil {
		return errors.Errorf("s.Update error:%v", err)
	}
	if err := myDB.Close(); err != nil {
		return errors.Errorf("db close error:%v", err)
	}
	return nil
}

func TestBalance(t *testing.T) {
	myDB, err := initBalanceDB()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer cleanBalanceDB(myDB)
	conn, err := getConn()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer conn.Close()
	balanceClient := balance.NewBalanceClient(conn)
	stream, err := balanceClient.Deduct(context.Background())
	if err != nil {
		t.Errorf("balanceClient.Deduct error: %v", err)
		return
	}
	deduct := &balance.DeductRequest{
		UserID:    UserID,
		ProductID: "234",
		Price:     "10",
	}
	if err := stream.Send(deduct); err != nil {
		t.Errorf("Failed to send a deduct: %v", err)
		return
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		t.Errorf("stream.CloseAndRecv error:%v", err)
		return
	}
	if reply.Result.Code != 1000 {
		t.Errorf("reply.Result.Code return code:%v, msg:%v", reply.Result.Code, reply.Result.Msg)
		return
	}
	if reply.Balance != "10" {
		t.Errorf("dedute %v from %v get %v", 10, 20, reply.Balance)
		return
	}
}
