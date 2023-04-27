## Compile protobuf
```sh
mkdir ./3_grpcExample/chat
protoc --go_out=plugins=grpc:3_grpcExample ./3_grpcExample/chat.proto
touch ./3_grpcExample/chat/chat.go
```

```sh
# if needed
go mod tidy
```