package main

import (
	"time"

	"com.aco.go.otel/initOtel"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {

	initOtel.Engage("service-go")

	r := gin.Default()
	r.GET("/", root)
	r.Run(":8087")
}

func root(c *gin.Context) {

	println("root called")

	traceProvider := initOtel.TraceProvider

	tracer := traceProvider.Tracer("go-test-tracer-child")

	// extract when there is context injection
	ctx := c.Request.Context()
	_, span := tracer.Start(
		otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header)),
		"span-name-child",
	)

	defer span.End()

	time.Sleep(250 * time.Millisecond)

	c.JSON(200, gin.H{
		"message": "service go return",
	})
}
