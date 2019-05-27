# 脚本使用说明
## 创建nginx-ingress-controller以及service
1. kubectl create -f mandatory.yaml
2. kubectl create -f cloud-generic.yaml
3. kubectl create -f ingress.yaml

## 创建etcd集群
1. kubectl create -f etcd-cluster.yaml
2. kubectl crreate -f etcd-deployment.yaml 
3. 如果要在集群外访问etcd，还需要 kubectl create -f etcd-client-service-lb.yaml

## 创建micro_framework服务
1. 创建配置文件的configmap对象 kubectl create configmap micro-config --from-file=config.toml
2. 创建micro_framework的deploy对象 kubectl create -f deployment.yaml
3. 创建service对象 kubectl create -f service.yaml
