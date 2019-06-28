package test

import (
	"context"
	"micro_framework/db"
	"micro_framework/pb/product"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func initProductDB() (*gorm.DB, error) {
	myDB, err := initDB()
	if err != nil {
		return nil, errors.Errorf("%v", err)
	}
	b := db.Product{
		ProductID:    ProductID,
		ProductName:  "product1",
		ProductPrice: 10,
	}
	if err := myDB.Save(&b).Error; err != nil {
		return nil, errors.Errorf("s.Update error:%v", err)
	}
	return myDB, nil
}

func cleanProductDB(myDB *gorm.DB) error {
	if err := myDB.Delete(db.Product{}, "productid=?", ProductID).Error; err != nil {
		return errors.Errorf("s.Update error:%v", err)
	}
	if err := myDB.Close(); err != nil {
		return errors.Errorf("db close error:%v", err)
	}
	return nil
}

func TestProduct(t *testing.T) {
	myDB, err := initProductDB()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer cleanProductDB(myDB)
	conn, err := getConn()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer conn.Close()
	productClient := product.NewProductClient(conn)
	productParam := &product.GetProductRequest{ProductID: ProductID}
	productRes, err := productClient.GetProduct(context.Background(), productParam)
	if err != nil {
		t.Errorf("productClient.GetProduct error: %v", err)
		return
	}
	if productRes.Result.Code != 1000 {
		t.Errorf("reply.Result.Code return code:%v, msg:%v", productRes.Result.Code, productRes.Result.Msg)
		return
	}
	if productRes.ProductName != "product1" {
		t.Errorf("get productname:%v from:%v", productRes.ProductName, ProductID)
		return
	}
	if productRes.ProductPrice != "10" {
		t.Errorf("order price:%v != %v", productRes.ProductPrice, 10)
		return
	}
}
