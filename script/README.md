# 脚本使用说明
## 创建nginx-ingress-controller以及service
1. kubectl create -f mandatory.yaml
2. kubectl create -f cloud-generic.yaml
3. kubectl create -f ingress.yaml
4. kubectl create secret tls tls-secret --key script/server.key --cert script/server.crt  -n ingress-nginx

* 附证书创建过程:  openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout script/server.key -out script/server.crt -subj "/CN=micro-framework.mydomain.com/O=micro-framework.mydomain.com"

## 创建etcd集群
1. kubectl create -f etcd-cluster.yaml
2. kubectl crreate -f etcd-deployment.yaml 
3. 如果要在集群外访问etcd，还需要 kubectl create -f etcd-client-service-lb.yaml

## 创建micro_framework服务
* premise 创建镜像 docker build --no-cache -t micro-framework ./
1. 创建配置文件的configmap对象 kubectl create configmap micro-config --from-file=config.toml
2. 创建micro_framework的deploy对象 kubectl create -f deployment.yaml
3. 创建service对象 kubectl create -f service.yaml
