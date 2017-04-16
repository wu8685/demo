#! /bin/sh

if [ -z $GOPATH ]; then
    echo "no GOPATH"
    exit 1
fi

cd $GOPATH/src/demo/failstart

godep go build -o failstart ./main.go

sudo docker build -t fail_start:v1 .