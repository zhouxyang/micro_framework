package test

import (
	"context"
	"fmt"
	"micro_framework/db"
	"micro_framework/pb/order"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func initOrderDB(t *testing.T) (*gorm.DB, error) {
	myDB, err := initDB()
	if err != nil {
		return nil, errors.Errorf("%v", err)
	}
	//insert user
	u := db.User{
		UserID: UserID,
		Name:   "zy",
	}
	if err := myDB.Save(&u).Error; err != nil {
		return nil, errors.Errorf("s.Update error:%v", err)
	}
	//insert balance
	b := db.Balance{
		UserID:  UserID,
		Balance: 20,
	}
	if err := myDB.Save(&b).Error; err != nil {
		return nil, errors.Errorf("s.Update error:%v", err)
	}
	//insert product
	prd := []db.Product{
		{
			ProductID:    ProductID,
			ProductName:  "product1",
			ProductPrice: 1.11,
		},
		{
			ProductID:    ProductID2,
			ProductName:  "product2",
			ProductPrice: 2.22,
		},
		{
			ProductID:    ProductID3,
			ProductName:  "product3",
			ProductPrice: 3.33,
		},
	}
	for _, p := range prd {
		if err := myDB.Save(&p).Error; err != nil {
			return nil, errors.Errorf("s.Update error:%v", err)
		}
	}
	return myDB, nil
}

func cleanOrderDB(myDB *gorm.DB) error {
	if err := myDB.Delete(db.Balance{}, "userid=?", UserID).Error; err != nil {
		return errors.Errorf("s.Update error:%v", err)
	}
	if err := myDB.Delete(db.Product{}, "productid=? or productid=? or productid=?", ProductID, ProductID2, ProductID3).Error; err != nil {
		return errors.Errorf("s.Update error:%v", err)
	}
	if err := myDB.Delete(db.User{}, "userid=?", UserID).Error; err != nil {
		return errors.Errorf("s.Update error:%v", err)
	}
	if err := myDB.Delete(db.Order{}, "productid=? or productid=? or productid=?", ProductID, ProductID2, ProductID3).Error; err != nil {
		return errors.Errorf("s.Update error:%v", err)
	}
	if err := myDB.Close(); err != nil {
		return errors.Errorf("db close error:%v", err)
	}
	return nil
}

func TestOrder(t *testing.T) {
	myDB, err := initOrderDB(t)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer cleanOrderDB(myDB)
	conn, err := getConn()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer conn.Close()
	client := order.NewOrderClient(conn)

	createReq := &order.OrderCreateRequest{
		UserID:    UserID,
		ProductID: []string{ProductID, ProductID2, ProductID3},
	}
	createResp, err := client.CreateOrder(context.Background(), createReq)
	if err != nil {
		t.Errorf("fail to CreateOrder:%v", err)
		return
	}
	if createResp.Result.Code != 1000 {
		t.Errorf("reply.Result.Code return code:%v, msg:%v", createResp.Result.Code, createResp.Result.Msg)
		return
	}
	if createResp.Balance != fmt.Sprintf("%v", 20-1.11-2.22-3.33) {
		t.Errorf("create order balance incorrect:%v", createResp.Balance)
		return
	}
	if len(createResp.OrderMsg.ProductID) != 3 {
		t.Errorf("len of reateResp.OrderMsg.ProductID:%v incorrect", len(createResp.OrderMsg.ProductID))
		return
	}

	//for query
	queryReq := &order.OrderQueryRequest{
		OrderID: createResp.OrderMsg.OrderID,
	}
	queryResp, err := client.QueryOrder(context.Background(), queryReq)
	if err != nil {
		t.Errorf("fail to QueryOrder:%v", err)
		return
	}

	if queryResp.Result.Code != 1000 {
		t.Errorf("reply.Result.Code return code:%v, msg:%v", queryResp.Result.Code, queryResp.Result.Msg)
		return
	}
	if queryResp.OrderMsg.TotolMoney != fmt.Sprintf("%v", 1.11+2.22+3.33) {
		t.Errorf("query order balance incorrect:%v", queryResp.Balance)
		return
	}
	if len(queryResp.OrderMsg.ProductID) != 3 {
		t.Errorf("len of reateResp.OrderMsg.ProductID:%v incorrect", len(queryResp.OrderMsg.ProductID))
		return
	}
}
