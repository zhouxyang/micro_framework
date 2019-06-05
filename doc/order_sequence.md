@startuml

客户端 --> order_service : 进行下单
order_service --> user_service : 进行用户校验
user_service --> order_service : 返回用户校验结果
order_service --> product_service : 校验并查询商品信息
product_service --> order_service : 返回商品价格
order_service --> balance_service : 请求变动余额
balance_service --> order_service : 返回变动后的余额
order_service --> order_service : 记录本次订单信息
order_service --> 客户端  : 返回本次下单信息以及余额

@enduml
