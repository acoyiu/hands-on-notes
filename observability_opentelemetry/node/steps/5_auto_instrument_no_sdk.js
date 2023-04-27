/**
 * OpenTelemetry HTTP Instrumentation allows the user
 * to automatically collect trace data for http and https.
 * 
 * because "/" and "sub_action" are using same nodejs(host) environment,
 * so that the module {HttpInstrumentation & registerInstrumentations}
 * can pipe up two request context automatically
 */

// basic
const { NodeTracerProvider } = require('@opentelemetry/sdk-trace-node');
const { ConsoleSpanExporter, SimpleSpanProcessor } = require('@opentelemetry/sdk-trace-base');

// exporter
const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');

// lib used
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');

// auto add http socket in out monitoring
const { HttpInstrumentation } = require('@opentelemetry/instrumentation-http');
const { registerInstrumentations } = require('@opentelemetry/instrumentation');

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
        [SemanticResourceAttributes.SERVICE_NAME]: 'aco-nodejs-auto-instru-no-sdk',
    }),
});

// 4: add span processor to provider
provider.addSpanProcessor(consoleLogSpanProcessor);
provider.addSpanProcessor(jaegerSpanProcessor);

// 5: register provider to global
provider.register();

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// Auto listen to http socket activities
registerInstrumentations({
    instrumentations: [new HttpInstrumentation()],
});

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

const express = require('express');
const app = express();
const port = 3000;
const axios = require('axios').default;


app.get('/', async (req, res) => {
    try {
        console.log('"/" called');
        await axios.get('http://localhost:3000/sub_action');
        res.send('root traced');
    }
    catch (error) {
        res.send('root traced with fetch error');
    }
});

app.get('/sub_action', async (_, res) => {
    if (Math.random() > 0.5) {
        res.status(200);
        res.send('sub_action traced')
    } else {
        res.status(400);
        res.send('none shall pass');
    }
});

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`);
});

