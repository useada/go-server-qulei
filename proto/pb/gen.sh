#!/bin/bash
echo $1
if [ $# == 1 ]; then 
    protoc -I ./$1 ./$1/*.proto --go_out=plugins=grpc:./$1/
else
    echo "请输入需要生成代码的文件夹"
fi
