/**
 * This registers a tracer provider with the OpenTelemetry API as the global tracer provider,
 * and exports a tracer instance that you can use to create spans.
 * 
 * If you do not register a global tracer provider, 
 * any instrumentation calls will be a no-op, so this is important to do!
 */

// register tracer
const { BasicTracerProvider, ConsoleSpanExporter, SimpleSpanProcessor } = require('@opentelemetry/sdk-trace-base');
const opentelemetry = require('@opentelemetry/api');

const provider = new BasicTracerProvider();

// Configure span processor to send spans to the exporter
provider.addSpanProcessor(new SimpleSpanProcessor(new ConsoleSpanExporter()));
provider.register();

// This is what we'll access in all instrumentation code
const tracer = opentelemetry.trace.getTracer('aco-basic-tracer-node');


// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


// Create a span. A span must be closed.
const parent_span = tracer.startSpan('mainTask');

const create_child_span = trace_name => {
    const opentelemetry = require('@opentelemetry/api');

    // Start another span. In this example, the main function already started a
    // span, so that'll be the parent span, and this will be a child span.

    // create context that inherit all info from parent
    const ctx = opentelemetry.trace.setSpan(opentelemetry.context.active(), parent_span);

    // inject parent's context into child process
    const span = tracer.startSpan(trace_name, undefined, ctx);

    return span;
}

// loop of nested
for (let i = 0; i < 10; i += 1) {

    const child_span = create_child_span('subTask');

    // simulate some random work.
    for (let i = 0; i <= Math.floor(Math.random() * 40000000); i += 1) { }

    // Make sure to end this child child_span! If you don't,
    // it will continue to track work beyond 'doWork'!
    child_span.end();
}


// Be sure to end the parent span.
parent_span.end();