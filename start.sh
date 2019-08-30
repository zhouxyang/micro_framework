#!/bin/bash

function init_ingress(){
    kubectl create -f script/ingress/mandatory.yaml
    kubectl create -f script/ingress/cloud-generic.yaml
    kubectl create -f script/ingress/ingress.yaml
    kubectl create secret tls tls-secret --key script/ingress/server.key --cert script/ingress/server.crt  -n ingress-nginx
}

function init_mysql(){
    kubectl create -f script/mysql/mysql-rc.yaml
	kubectl create -f script/mysql/mysql-svc.yaml
}

function init_etcd(){
    kubectl create -f script/etcd/etcd-deployment.yaml 
    kubectl create -f script/etcd/etcd-client-service-lb.yaml
    kubectl create -f script/etcd/etcd-cluster.yaml
	kubectl create -f script/etcd/hazelcast-rbac.yaml
}

function destory_ingress(){
    kubectl delete -f script/ingress/mandatory.yaml
    kubectl delete -f script/ingress/cloud-generic.yaml
    kubectl delete -f script/ingress/ingress.yaml
    kubectl delete secret tls-secret -n ingress-nginx
}

function destory_mysql(){
    kubectl delete -f script/mysql/mysql-rc.yaml
    kubectl delete -f script/mysql/mysql-svc.yaml
}

function destory_etcd(){
    kubectl delete -f script/etcd/etcd-cluster.yaml
    kubectl delete -f script/etcd/etcd-deployment.yaml 
    kubectl delete -f script/etcd/etcd-client-service-lb.yaml
	kubectl delete -f script/etcd/hazelcast-rbac.yaml
}

function start(){
    kubectl create configmap micro-config --from-file=config.toml
    kubectl create -f script/micro-framework/deployment.yaml
    kubectl create -f script/micro-framework/service.yaml
}
function stop(){
    kubectl delete configmap micro-config
    kubectl delete -f script/micro-framework/deployment.yaml
    kubectl delete -f script/micro-framework/service.yaml
}

function help() {
    echo "USAGE: $0 init|destory|start|stop|check|restart"
    exit 1
}


if [ "$1" == "" ]; then
    help
elif [ "$1" == "stop" ];then
    stop
elif [ "$1" == "start" ];then
    start
elif [ "$1" == "init" ];then
    init_mysql
    init_etcd
	init_ingress
elif [ "$1" == "destory" ];then
    destory_mysql
    destory_etcd
	destory_ingress
elif [ "$1" == "restart" ];then
    stop
	start
else
    help
fi

