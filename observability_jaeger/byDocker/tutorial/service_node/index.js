const opentelemetry = require('@opentelemetry/api');
const { propagation, ROOT_CONTEXT } = require('@opentelemetry/api');
const { NodeTracerProvider } = require('@opentelemetry/node');
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');
const { ConsoleSpanExporter, SimpleSpanProcessor } = require('@opentelemetry/tracing');
const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');

const consoleExporter = new ConsoleSpanExporter();
const jaegerTheExporter = new JaegerExporter({ endpoint: 'http://localhost:14268/api/traces' });

const consoleLogSpanProcessor = new SimpleSpanProcessor(consoleExporter);
const jaegerSpanProcessor = new SimpleSpanProcessor(jaegerTheExporter);

const provider = new NodeTracerProvider({
    resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'service-node',
    }),
});

provider.addSpanProcessor(consoleLogSpanProcessor);
provider.addSpanProcessor(jaegerSpanProcessor);

provider.register();

const tracer = opentelemetry.trace.getTracer('aco-manual-root-route');

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

const express = require('express');
const cors = require('cors');
const app = express();
const port = 8086;
const axios = require('axios').default;

app.use(cors())

app.get('/', async (req, res) => {

    // extract the base context from header
    const fromContext /* :BaseContext */ = propagation.extract(ROOT_CONTEXT, req.headers);

    const parentSpan = tracer.startSpan('aco-manual-root-calling', undefined, fromContext);

    const headers = { aco: 'node' };
    propagation.inject(opentelemetry.trace.setSpan(ROOT_CONTEXT, parentSpan), headers);
    console.log('Injected headers', headers);

    await new Promise(res => setTimeout(res, 100));

    try {
        const { data } = await axios.get('http://localhost:8087', { headers });
        res.json({ data });
    }
    catch (error) {
        parentSpan.recordException(error);
        res.send('root traced Errored');
    }

    await new Promise(res => setTimeout(res, 100));

    parentSpan.end();
});

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`);
});

