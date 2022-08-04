#!/bin/bash

serverName="user"

if [ -f "${serverName}" ] ;then
     rm "${serverName}"
fi
if [ -f "${serverName}" ] ;then
     rm ${serverName}
fi

sleep 0.2

#CGO_ENABLED=0 GOOS=linux go build -o ${serverName}
go build -o ${serverName}

# 运行服务
./${serverName}
