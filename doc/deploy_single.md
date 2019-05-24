@startuml
left to right direction
actor client

node nginx

cloud server{
node server1
node server2
node server3
node etcd
server1 -- etcd
server2 -- etcd
server3 -- etcd
server1 -- server2
server1 -- server3
server2 -- server3
}

client -- nginx

nginx -- server


@enduml

