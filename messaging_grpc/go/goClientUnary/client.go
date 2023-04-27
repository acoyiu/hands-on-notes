package main

import (
	"log"

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

	// Step 3: 使用 gRPC connection 呼叫 remote function
	response, err := c.UnaryRequest(context.Background(), &chat.Message{Body: "Hello From Unary!"})
	if err != nil {
		log.Fatalf("Error when calling UnaryRequest: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)
}
