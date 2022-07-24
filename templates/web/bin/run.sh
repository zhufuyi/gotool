#!/bin/bash

serverName="userExample"
#CGO_ENABLED=0 GOOS=linux go build -o ${serverName}
go build -o ${serverName}
#echo "完成编译${serverName}"

# 运行服务
./${serverName}
