package mysql

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"route_guide/db"
)

// OrderDao access to Order
type OrderDao struct {
	*db.MyDB
}

func (dao *OrderDao) InsertOrder(log *logrus.Entry, order *db.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}

	query := `INSERT IGNORE INTO Orders(orderid, product_price, create_time, productid) VALUES (?, ?, ?, ?)`

	_, err := dao.Exec(log, query, order.OrderID, order.OrderPrice.String(), order.CreateTime, order.ProductID)
	if err != nil {
		return errors.Wrap(err, "insert Order")
	}
	return nil
}

func (dao *OrderDao) GetOrderByOrderID(log *logrus.Entry, orderid string) ([]*db.Order, error) {
	query := `SELECT orderid, product_price, create_time, productid FROM Orders`
	orders := []*db.Order{}
	err := dao.Select(log, &orders, query)

	switch {
	case err != nil:
		return nil, errors.Wrap(err, "db get GetOrderByOrderID")
	}
	return orders, nil
}
