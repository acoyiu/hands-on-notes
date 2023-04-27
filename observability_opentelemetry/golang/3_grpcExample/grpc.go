package grpc_example

//
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"com.aco.go.otel/3_grpcExample/chat"
	"com.aco.go.otel/initOtel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Start() {

	initOtel.Engage("go-grpc")

	// start gRPC server
	go startGrpcServer()

	r := gin.Default()
	r.GET("/", g01)
	println(":: curl localhost:6060/g01")
	r.Run(":6060")
}

//

func g01(c *gin.Context) {

	// tracing =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	traceProvider := initOtel.TraceProvider

	tracer := traceProvider.Tracer("go-grpc-parent-trace")

	ginTraceCtx, span := tracer.Start(c, "go-grpc-gin-span")
	defer span.End()

	// grpc =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	var conn *grpc.ClientConn

	// use grpc.DialContext if wants non-one time connection!!
	conn, err := grpc.Dial("localhost:6061", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	clientService := chat.NewChatServiceClient(conn)

	// propagation trace extraction =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

	theMetadata := metadata.New(map[string]string{})

	otelgrpc.Inject(ginTraceCtx, &theMetadata)

	println(theMetadata["traceparent"][0])

	grpcSendCtx := metadata.AppendToOutgoingContext(ginTraceCtx, []string{"traceparent", theMetadata["traceparent"][0]}...)

	// grpc call =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	time.Sleep(100 * time.Millisecond)
	response, err := clientService.SayHello(grpcSendCtx, &chat.Message{Body: "Hello From Client!"}, grpc.Header(&theMetadata))
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	time.Sleep(100 * time.Millisecond)

	log.Printf("Response from server: %s", response.Body)

	// Gin =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	c.JSON(200, gin.H{
		"message": response.Body,
	})
}

//

type ChatServer struct{}

func (s *ChatServer) SayHello(ctx context.Context, inComingMsg *chat.Message) (*chat.Message, error) {

	// propagation trace extraction =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &chat.Message{Body: "Get metadata error"}, nil
	}

	println("gRPC server md")
	PrettyPrint(md)

	_, parentSpanContext := otelgrpc.Extract(ctx, &md)

	// start child span =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

	traceProvider := initOtel.TraceProvider

	tracer := traceProvider.Tracer("go-grpc-child-trace")

	grpcTraceCtx := trace.ContextWithSpanContext(ctx, parentSpanContext)

	_, span := tracer.Start(
		grpcTraceCtx,
		"go-grpc-grpc-span",
	)

	defer span.End()

	// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

	time.Sleep(100 * time.Millisecond)

	log.Printf("Receive message body from client: %s", inComingMsg.Body)
	return &chat.Message{Body: "Hello From the Server!"}, nil
}

func startGrpcServer() {

	lis, err := net.Listen("tcp", ":6061")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serviceWrapper := ChatServer{}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &serviceWrapper)

	println("gRPC server started:..")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

//

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
