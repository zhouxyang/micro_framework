# Description
1. 使用etcd服务发现，etcd采用短连接
2. 服务的创建采用注册方式，减少业务代码与框架的耦合
3. 优雅关闭，热重启 
4. 使用logrus，结构化日志; 统一requestid，分布式追踪; 日志支持filename，行号
5. 添加Dockerfile，支持容器化
6. 支持容器编排工具k8s,添加k8s部署的脚本文件
