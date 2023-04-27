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


// create span

// Create a span. A span must be closed.
const parent_span = tracer.startSpan('main');
console.log('Span started.');

// simulate real works
for (let i = 0; i < 5; i += 1)
    for (let i = 0; i <= Math.floor(Math.random() * 2); i += 1)
        console.log('work...');

// Be sure to end the span.
parent_span.end();
console.log('Span ended.');