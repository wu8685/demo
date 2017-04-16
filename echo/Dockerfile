FROM golang:1.8

CMD go get github.com/tools/godep && \
    godep go build -o echo cmd/main.go

EXPOSE 8080

ENTRYPOINT [./echo, --port, 8080]

