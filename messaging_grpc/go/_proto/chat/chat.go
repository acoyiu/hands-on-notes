package chat

import (
	context "context"
	"fmt"
	"log"
	"time"
)

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
