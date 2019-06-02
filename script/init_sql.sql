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


insert into `Balance` values ("123",20);
insert into  `User` values('123','zy');
insert into `Product` values('1','p1',1.11);
insert into `Product` values('2','p2',2.22);
insert into `Product` values('3','p3',3.33);
