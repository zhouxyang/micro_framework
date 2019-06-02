package db

import (
	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
)

//User Info
type User struct {
	UserID string `db:"userid"`
	Name   string `db:"name"`
}

// Product  Info
type Product struct {
	ProductID    string          `db:"productid"`
	ProductName  string          `db:"product_name"`
	ProductPrice decimal.Decimal `db:"product_price"`
}

// Order Info
type Order struct {
	OrderID      string          `db:"orderid"`
	ProductPrice decimal.Decimal `db:"product_price"`
	CreateTime   mysql.NullTime  `db:"create_time"`
	ProductID    string          `db:"productid"`
}

// Balance info
type Balance struct {
	UserID  string          `db:"userid"`
	Balance decimal.Decimal `db:"balance"`
}
