FROM golang:1.19 AS BUILDER

COPY . /go/src/github.com/wu8685/demo
WORKDIR /go/src/github.com/wu8685/demo

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server echo/cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /client echo/cmd/client/main.go

FROM centos:centos7

RUN yum install -y telnet bind-utils

COPY --from=BUILDER /server /server
COPY --from=BUILDER /client /client

EXPOSE 8080
EXPOSE 2112

