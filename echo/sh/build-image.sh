#! /bin/bash

if [ -z $GOPATH ]; then
    echo "no GOPATH"
    exit 1
fi

cd $GOPATH/src/demo/echo
 
godep go build -o echo ./cmd/main.go

sudo docker build -t echo:v1 .
