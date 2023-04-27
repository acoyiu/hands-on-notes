package initOtel

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var TraceProvider *tracesdk.TracerProvider

func Engage(nameOfService string) {

	/**
	 * 1: create exporter: Jaeger exporter
	 */
	traceExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint("http://localhost:14268/api/traces"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	/**
	 * 2: create batcher => span processor in nodejs, Always be sure to batch in production.
	 */
	processor := tracesdk.WithBatcher(traceExporter)

	/**
	 * 3: set service name when create provider
	 */
	TraceProvider = tracesdk.NewTracerProvider(

		/**
		 * 4: add span processor to provider like nodejs
		 */
		processor,

		tracesdk.WithResource( // Record information about this application in a Resource.
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(nameOfService),
				attribute.String("environment", "production"),
				attribute.Int64("ID", 1),
			),
		),
	)

	/**
	 * 5: Register our TracerProvider as the global
	 */
	// so any imported instrumentation in the future will default to using it.
	// Same as nodeJS register
	otel.SetTracerProvider(TraceProvider)

	/**
	 * 6: Must Set "SetTextMapPropagator" to use "propagation"
	 */
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
