const opentelemetry = require('@opentelemetry/api');
const { propagation, ROOT_CONTEXT } = require('@opentelemetry/api');
const { NodeTracerProvider } = require('@opentelemetry/node');
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');
const { ConsoleSpanExporter, SimpleSpanProcessor } = require('@opentelemetry/tracing');
const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// 1: create exporter for jaeger and console.log
const consoleExporter = new ConsoleSpanExporter();
const jaegerTheExporter = new JaegerExporter({ endpoint: 'http://localhost:14268/api/traces' });

// 2: create span processor
const consoleLogSpanProcessor = new SimpleSpanProcessor(consoleExporter);
const jaegerSpanProcessor = new SimpleSpanProcessor(jaegerTheExporter);

// 3: set service name when create provider
const provider = new NodeTracerProvider({
    resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'aco-nodejs-manual-instru-seperate-B',
    }),
});

// 4: add span processor to provider
provider.addSpanProcessor(consoleLogSpanProcessor);
provider.addSpanProcessor(jaegerSpanProcessor);

// 5: register provider to global
provider.register();

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// Create a Trace Type
const tracer = opentelemetry.trace.getTracer('aco-sep-manual-root-route-B');

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

const express = require('express');
const app = express();
const port = 3001;

app.get('/sub_action', async (req, res) => {

    console.log('req.headers');
    console.log(req.headers);

    // extract the base context from header
    const fromContext /* :BaseContext */ = propagation.extract(ROOT_CONTEXT, req.headers);

    // create child span for tracing child activities
    const childSpan = tracer.startSpan('aco-manual-child-call', undefined, fromContext);

    // response after 3 sec
    await new Promise(res => setTimeout(res, 500));

    // dynamic error
    if (Math.random() > 0.5) {
        res.status(200);
        res.send('sub_action called');
    }
    else {

        // add exception to log
        childSpan.recordException('num smaller exception');

        // add tags to label as error
        childSpan.setStatus({
            code: opentelemetry.SpanStatusCode.ERROR,
            message: 'exception: num is smaller than 0.5',
        });

        res.status(400);
        res.send('none shall pass');
    }

    // end the child span
    childSpan.end();
});

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`);
});

