#!/bin/bash

#把pod IP写到配置文件的Host中，并生成新的配置文件
sed  "s/0.0.0.0/$MY_POD_IP/g" /etc/config/config.toml > /go/src/micro_framework/config.toml 
micro_framework run /go/src/micro_framework/config.toml
