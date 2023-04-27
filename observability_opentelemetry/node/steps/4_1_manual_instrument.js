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
        [SemanticResourceAttributes.SERVICE_NAME]: 'aco-nodejs-manual-instru-seperate-A',
    }),
});

// 4: add span processor to provider
provider.addSpanProcessor(consoleLogSpanProcessor);
provider.addSpanProcessor(jaegerSpanProcessor);

// 5: register provider to global
provider.register();

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// Create a Trace Type
const tracer = opentelemetry.trace.getTracer('aco-sep-manual-root-route-A');

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

const express = require('express');
const app = express();
const port = 3000;
const axios = require('axios').default;

app.get('/', async (_, res) => {

    const parentSpan = tracer.startSpan('aco-manual-root-calling');

    // inject parent context into header which will send to "sub_action"
    const headers = { aco: 'yes' };
    propagation.inject(
        { getValue: () => parentSpan },
        headers,
    );

    // Make a request to the sub path
    try {
        await axios.get('http://localhost:3001/sub_action', { headers });
        res.send('root traced');
    }
    catch (error) {
        // console.log(error);

        // add exception to log, but success overall
        parentSpan.recordException(error);
        res.send('root traced Errored');
    }

    // Be sure to end the span. 
    parentSpan.end();
});

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`);
});

