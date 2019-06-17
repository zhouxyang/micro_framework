package db

import (
	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
)

//User Info
type User struct {
	UserID string `gorm:"primary_key;column:userid"`
	Name   string `gorm:"column:name"`
}

// TableName of User
func (User) TableName() string {
	return "User"
}

// Product  Info
type Product struct {
	ProductID    string          `gorm:"primary_key;column:productid"`
	ProductName  string          `gorm:"column:product_name"`
	ProductPrice decimal.Decimal `gorm:"column:product_price;type:decimal(20,4)"`
}

// TableName of Product
func (Product) TableName() string {
	return "Product"
}

// Order Info
type Order struct {
	OrderID      string         `gorm:"primary_key;column:orderid"`
	ProductPrice float64        `gorm:"column:product_price"`
	CreateTime   mysql.NullTime `gorm:"column:create_time"`
	ProductID    string         `gorm:"column:productid"`
}

// TableName of Order
func (Order) TableName() string {
	return "Orders"
}

// Balance info
type Balance struct {
	UserID  string  `gorm:"primary_key;column:userid"`
	Balance float64 `gorm:"column:balance"`
}

// TableName of Balance
func (Balance) TableName() string {
	return "Balance"
}
