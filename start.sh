#!/bin/bash

#把pod IP写到配置文件的Host中，并生成新的配置文件
printenv
echo $MY_POD_IP
sed  "s/0.0.0.0/$MY_POD_IP/g" /etc/config/config.toml > /go/src/route_guide/config.toml 
micro_framework run /go/src/route_guide/config.toml