package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"grpc.test/_proto/chat"
)

func main() {

	// 1, Create TCP socket occupation
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2, Instantiate custom wrapper of gRPC functions
	serviceWrapper := chat.Server{}

	// 3, Create gRPC server for later bind with App's service
	grpcServer := grpc.NewServer()

	// 4, 將我們的 wrapper bind 至 gRPC server
	chat.RegisterChatServiceServer(grpcServer, &serviceWrapper)

	// 5, gRPC server bind 至 localhost 的 port
	println("gRPC server started:..")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
