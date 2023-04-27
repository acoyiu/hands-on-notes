package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"grpc.test/_proto/chat"
)

func main() {

	// Step 1: dial to gRPC server
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	// Step 2: gRPC 生成的 code 裡面的 NewChatServiceClient()
	c := chat.NewChatServiceClient(conn)

	streamResponse, err := c.ServerStreamRequest(context.Background(), &chat.Message{Body: "Hello From Server Stream!"})
	if err != nil {
		log.Fatalf("Error when calling ServerStreamRequest: %s", err)
	}

	for i := 0; i < 12; i++ {
		time.Sleep(time.Second * 1)

		msg, err := streamResponse.Recv()
		if err != nil {
			log.Fatalf("Error when calling streamResponse.Recv(): %s", err)
		}

		if len(msg.Body) > 0 {
			log.Println(msg.Body)
		} else {
			log.Println("no stream message")
		}
	}
}
