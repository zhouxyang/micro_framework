@startuml
title 单集群下单服务架构图

left to right direction
actor client

node nginx

cloud server{
node order
node user
node product
node balance
node etcd
order --> user
order --> product
order --> balance
order --> etcd
}

client -- nginx

nginx -- server


@enduml
