## 1: install protoc (executable)

[Download](https://github.com/protocolbuffers/protobuf/tags) <br/>

```sh
MacOS 可以放在：
/usr/local/bin        放bin資料夾底下的東西
/usr/local/include    放include資料夾底下的東西
```

```sh
Windows 可以放：
C:\Users\USER\go\bin         放bin資料夾底下的東西
C:\Users\USER\go\include     放include資料夾底下的東西
```

```sh
# install 兩個 compile protobuf 既 binary executable
# will in user's root, need export path
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

<br/>

---

<br/>

## 2: Generate Protobuf

protoc --go_out=[ExportDirectory] --go-grpc_out=[ExportDirectory][protobuffile]

```sh
rm -r _proto
mkdir _proto
protoc --go_out=_proto --go-grpc_out=_proto ./chat.proto
```

.pb.go => Data Type <br/>
\_grpc.pb.go => Service Interface <br/>

<br/>

---

<br/>

## 3: Add below to /protp/chat/chat.go

```go
package chat

// will have auto import

type Server struct {
	// need to extend/inherite from RGPC Server Interface
	ChatServiceServer // mustEmbedUnimplementedChatServiceServer will be nil
}

func (s *Server) UnaryRequest(ctx context.Context, incomingMsg *Message) (*Message, error) {
	log.Printf("UnaryRequest receive message body from client: %s", incomingMsg.Body)
	return &Message{Body: "Hello From the Server!"}, nil
}

func (s *Server) ServerStreamRequest(incomingMsg *Message, stream ChatService_ServerStreamRequestServer) error {
	log.Printf("ServerStreamRequest receive message body from client: %s", incomingMsg.Body)

	stream.Send(&Message{
		Body: "1",
	})

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
		stream.Send(&Message{
			Body: fmt.Sprintf("%v", i+2),
		})
	}

	// Hang up connection for later reuse
	time.Sleep(time.Second * 3)
	stream.Send(&Message{
		Body: "Before Close",
	})
	time.Sleep(time.Second * 3)

	return nil
}

```

<br/>

---

<br/>

## 4: Run gRPC server

```sh
# install package
go mod tidy

go run goServer/server.go
```

## 5: Run gRPC client

```sh
go run goClientUnary/client.go
go run goClientServerStream/client.go
```

<br/>

---

<br/>

注意⚠️：

- Stream 的話，Send/Close 都是是由負責串流的一方決定的
- Receive 調用 Recv()，會取得這段時間內收到的 Event 的 Queue（先進先出）