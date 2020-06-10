#!/bin/bash

function init_ingress(){
    kubectl create -f script/ingress/mandatory.yaml
    kubectl create -f script/ingress/cloud-generic.yaml
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
}

function destory_ingress(){
    kubectl delete -f script/ingress/mandatory.yaml
    kubectl delete -f script/ingress/cloud-generic.yaml
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
}

function start(){
    kubectl create configmap micro-config --from-file=config.toml
    kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v1/ --docker-username=gongfupanda2 --docker-password=090405docker --docker-email=15623492306@163.com
	helm install script/mychart --name mychart-micro-framework 
}
function stop(){
    kubectl delete configmap micro-config
    kubectl delete secret regcred
	helm delete --purge mychart-micro-framework
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

