const opentelemetry = require('@opentelemetry/api');
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');
const { BasicTracerProvider, ConsoleSpanExporter, SimpleSpanProcessor } = require('@opentelemetry/sdk-trace-base');
const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');


// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// 1: Create and configure span processor to send spans to the exporter
const jaegerTheExporter = new JaegerExporter({ endpoint: 'http://localhost:14268/api/traces' });
const consoleExporter = new ConsoleSpanExporter();

// 2: Create "SpanProcessor" to handle the way to export and where to export
const jaegerSpanProcessor = new SimpleSpanProcessor(jaegerTheExporter);
const consoleLogSpanProcessor = new SimpleSpanProcessor(consoleExporter);

// 3: create "Provider" with "Resource"
const provider = new BasicTracerProvider({
    resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'aco-service-name',
        // if not set, will be as 'unknown_service'
    }),
});

// 4: add "Processor" to "Provider" by provider.addSpanProcessor()
provider.addSpanProcessor(jaegerSpanProcessor);
provider.addSpanProcessor(consoleLogSpanProcessor);

// 5: Register global "Provider"
provider.register();


// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// Create trace by using "opentelemetry.trace.getTracer" with global "Provider"

const tracer = opentelemetry.trace.getTracer('aco-insideService-trace-name');


// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// 7: Tracing in progress

// Create a span. A span must be closed.
const parentSpan = tracer.startSpan('aco-insideTrace-span-name');

// simulate works
for (let i = 0; i < 10; i += 1) {
    doWork(parentSpan);
}

// Be sure to end the span.
parentSpan.end();

// flush and close the connection.
jaegerTheExporter.shutdown();


// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// Sub-Trace

function doWork(parent) {
    // Start another span. In this example, the main method already started a
    // span, so that'll be the parent span, and this will be a child span.
    const ctx = opentelemetry.trace.setSpan(opentelemetry.context.active(), parent);
    const span = tracer.startSpan('doWork', undefined, ctx);

    // simulate some random work.
    for (let i = 0; i <= Math.floor(Math.random() * 40000000); i += 1) {
        // empty
    }

    // Set attributes to the span.
    span.setAttribute('key', 'value');

    // Annotate our span to capture metadata about our operation
    span.addEvent('invoking doWork');

    span.end();
}
