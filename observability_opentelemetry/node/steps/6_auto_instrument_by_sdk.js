const opentelemetry = require('@opentelemetry/sdk-node');

// exporter
const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');
const { ConsoleSpanExporter } = require('@opentelemetry/sdk-trace-base');

// lib used
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');


// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


/* Different Auto Instrumentations */

// Very Detail
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');

// HTTP level
const { HttpInstrumentation } = require('@opentelemetry/instrumentation-http');


// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


const sdk = new opentelemetry.NodeSDK({

    resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'aco-nodejs-auto-instru-with-sdk',
    }),

    traceExporter:
        // new ConsoleSpanExporter(),
        new JaegerExporter({ endpoint: 'http://localhost:14268/api/traces' }),

    instrumentations: [
        // getNodeAutoInstrumentations(),
        new HttpInstrumentation(),
    ]
});

sdk.start().then(() => {
    // Resources have been detected and SDK is started
    startExpress();
});


// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-


const startExpress = () => {

    /* below same as 5_, beside added shut down SDK */

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

    app.get('/shut_down', async (_, res) => {
        await sdk.shutdown();
        res.send('shuted down');
        process.exit(0);
    });

    app.listen(port, () => {
        console.log(`Example app listening on port ${port}`);
    });
};