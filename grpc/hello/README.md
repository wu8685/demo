# update .pb.go file after modifying .proto file

```
protoc -I <path-to-proto> <proto-files> --go_out=plugins=grpc:<path-to-output>
```