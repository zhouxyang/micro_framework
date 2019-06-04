package mysql

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"micro_framework/db"
)

// ProductDao access to Product
type ProductDao struct {
	*db.MyDB
}

func (dao *ProductDao) InsertProduct(log *logrus.Entry, product *db.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	query := `INSERT IGNORE INTO Product(productid, product_name, product_price)
          VALUES (?, ?, ?)`

	_, err := dao.Exec(log, query, product.ProductID, product.ProductName, product.ProductPrice)
	if err != nil {
		return errors.Wrap(err, "insert Product")
	}
	return nil
}

func (dao *ProductDao) GetProductByProductID(log *logrus.Entry, productid string) (*db.Product, error) {
	query := `SELECT productid, product_name, product_price FROM Product WHERE productid = ?`
	product := &db.Product{}
	err := dao.Get(log, product, query, productid)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, errors.Wrap(err, "db GetProductByProductID")
	}
	return product, nil
}
