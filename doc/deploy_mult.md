@startuml
left to right direction
actor client

node nginx1
node nginx2

cloud server1{
node server1.1
node server1.2
node server1.3
node etcd1
server1.1 -- etcd1
server1.2 -- etcd1
server1.3 -- etcd1
server1.1 -- server1.2
server1.1 -- server1.3
server1.2 -- server1.3
}

cloud server2{
node server2.1
node server2.2
node server2.3
node etcd2
server2.1 -- etcd2
server2.2 -- etcd2
server2.3 -- etcd2
server2.1 -- server2.2
server2.1 -- server2.3
server2.2 -- server2.3
}

client -- nginx1
client -- nginx2

nginx1 -- server1
nginx2 -- server2


@enduml

