import opentelemetry from '@opentelemetry/api';
import { ConsoleSpanExporter, SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';

import { Resource } from '@opentelemetry/resources';
import { ZipkinExporter } from '@opentelemetry/exporter-zipkin';

// 1: Create Exporter
const consoleExporter = new ConsoleSpanExporter();
const zipkinExporter = new ZipkinExporter({
    headers: { 'Content-Type': 'application/json;', }, // <- !!!! Must add for export to Jaeger backend !!!!
    url: 'http://localhost:9411/api/v2/spans',
});

// 2: Create processor
const consoleLogSpanProcessor = new SimpleSpanProcessor(consoleExporter);
const zipkinSpanProcessor = new SimpleSpanProcessor(zipkinExporter);

// 3: Create provider with name and details
const provider = new WebTracerProvider({
    resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'aco-frontend-web',
    }),
});

// 4: Add processor to provider
provider.addSpanProcessor(consoleLogSpanProcessor);
provider.addSpanProcessor(zipkinSpanProcessor);

// 5: Register as global provider for "@opentelemetry/api" use
provider.register();

const tracer = provider.getTracer('aco-example-tracer-web');

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

const but = document.getElementById('button1');
but.addEventListener('click', async () => {

    const parent_span = tracer.startSpan('main');
    console.log('Span started.');
    but.setAttribute('disabled', true);

    {
        await new Promise(res => setTimeout(res, 250));

        await childProcess(parent_span);

        await new Promise(res => setTimeout(res, 250));

        parent_span.end();
    }

    console.log('Parent span ended.');
    but.removeAttribute('disabled');
    alert('check localhost on port 9411');
});


const childProcess = async parentSpan => {

    const ctx = opentelemetry.trace.setSpan(opentelemetry.context.active(), parentSpan);
    const childSpan = tracer.startSpan('sub', undefined, ctx);

    // // inject by implement context interface
    // const childSpan = tracer.startSpan('sub', undefined, { getValue: () => parentSpan });

    await new Promise(res => setTimeout(res, 250));
    childSpan.end();
    console.log(childSpan);

    return;
};