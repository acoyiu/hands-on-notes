package basic

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"com.aco.go.otel/initOtel"
	"go.opentelemetry.io/otel/attribute"
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.

func Start() {

	initOtel.Engage("go-basic")

	// 1: get Provider
	traceProvider := initOtel.TraceProvider

	// 2: create tracer
	tracer := traceProvider.Tracer("go-test-tracer-parent")

	// 3: create Go's context and trace's span By tracer, can use Gin's context as well
	ctx, span := tracer.Start(context.Background(), "span-name-parent")

	// **: handle shutdown properly so nothing leaks, not the case if with Gin
	// **: need to add if in instant finish situation, because span is not send out
	// **: at once but batch like, so shutdown will send out all the span not yet sent
	// **: right a way
	defer func(ctx context.Context) {
		traceProvider.Shutdown(ctx)
	}(ctx)

	// 4: defer end trace at the end
	defer span.End()

	// await time
	time.Sleep(250 * time.Millisecond)

	print(childFunction(ctx))

	// await time
	time.Sleep(250 * time.Millisecond)
}

func childFunction(ctx context.Context) bool {

	// 1: get Provider
	traceProvider := initOtel.TraceProvider

	// 2: create tracer
	tracer := traceProvider.Tracer("go-test-tracer-child")

	// 2: Can use global TracerProvider to not using scope var "traceProvider"
	// tracer := otel.Tracer("go-test-tracer-child")

	// 3: create Go's context and trace's span By tracer
	_, span := tracer.Start(ctx, "span-name-child")

	// 4: end trace at the end
	defer span.End()

	// add custom tag
	span.SetAttributes(attribute.Key("tag-key").String("tag-value"))

	// add trace event point in "Logs"
	span.AddEvent("Before wait")
	time.Sleep(250 * time.Millisecond)
	span.AddEvent("After wait")

	// random add error tag
	shouldSucc := (rand.New(rand.NewSource(time.Now().UnixNano()))).Int()%2 > 0
	if !shouldSucc {

		// show as errored in "Tags"
		span.SetAttributes(attribute.Key("error").String("true"))

		// exception.message in "Logs"
		span.RecordError(errors.New("shouldSucc is smaller than 1"))
		println("ERROR ed")
	}

	return shouldSucc
}
