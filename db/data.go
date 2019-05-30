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
	OrderID    string          `db:"orderid"`
	OrderPrice decimal.Decimal `db:"order_price"`
	CreateTime mysql.NullTime  `db:"create_time"`
	ProductID  string          `db:"productid"`
}

// Balance info
type Balance struct {
	UserID  string          `db:"userid"`
	Balance decimal.Decimal `db:"balance"`
}

/*
  CREATE TABLE `User` (
    `userid` varchar(255) NOT NULL DEFAULT '' COMMENT 'user唯一id',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
    PRIMARY KEY (`userid`)
  ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='user信息表';


  CREATE TABLE `Product` (
    `productid` varchar(255) NOT NULL DEFAULT '' COMMENT 'product唯一id',
    `product_name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
    `product_price` decimal(14,4) NOT NULL DEFAULT '0.0000' COMMENT '订单总价(格式：0.0000)',
    PRIMARY KEY (`productid`)
  ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='product信息表';

  CREATE TABLE `Balance` (
    `userid` varchar(255) NOT NULL DEFAULT '' COMMENT 'user唯一id',
    `balance` decimal(14,4) NOT NULL DEFAULT '0.0000' COMMENT '用户余额',
    PRIMARY KEY (`userid`)
  ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='余额信息表';

  CREATE TABLE `Orders` (
    `orderid` varchar(255) NOT NULL DEFAULT '' COMMENT 'order唯一id',
    `product_price` decimal(14,4) NOT NULL DEFAULT '0.0000' COMMENT '订单总价(格式：0.0000)',
	`create_time` datetime DEFAULT NULL COMMENT '创建订单时间',
    `productid` varchar(255) NOT NULL DEFAULT '' COMMENT 'product唯一id',
    PRIMARY KEY (`productid`, `orderid`)
  ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='order信息表';

*/
