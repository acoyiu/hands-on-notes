import opentelemetry from '@opentelemetry/api';
import { propagation, ROOT_CONTEXT } from '@opentelemetry/api';
import { ConsoleSpanExporter, SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';

import { Resource } from '@opentelemetry/resources';
import { ZipkinExporter } from '@opentelemetry/exporter-zipkin';

const consoleExporter = new ConsoleSpanExporter();
const zipkinExporter = new ZipkinExporter({
    headers: { 'Content-Type': 'application/json;', },
    url: 'http://localhost:9411/api/v2/spans',
});

const consoleLogSpanProcessor = new SimpleSpanProcessor(consoleExporter);
const zipkinSpanProcessor = new SimpleSpanProcessor(zipkinExporter);

const provider = new WebTracerProvider({
    resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'service-web',
    }),
});

provider.addSpanProcessor(consoleLogSpanProcessor);
provider.addSpanProcessor(zipkinSpanProcessor);

provider.register();

const tracer = provider.getTracer('service-web-calling');

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

const but = document.getElementById('button1');
but.addEventListener('click', async () => {

    const parent_span = tracer.startSpan('main-but-func');
    console.log('Span started.');
    but.setAttribute('disabled', true);

    {
        await new Promise(res => setTimeout(res, 100));

        await childProcess(parent_span);

        await new Promise(res => setTimeout(res, 100));

        parent_span.end();
    }

    console.log('Parent span ended.');
    but.removeAttribute('disabled');
    alert('check localhost on port 9411');
});



const childProcess = async parentSpan => {

    const ctx = opentelemetry.trace.setSpan(opentelemetry.context.active(), parentSpan);
    const childSpan = tracer.startSpan('child-func', undefined, ctx);

    // inject context into header
    const headers = { aco: 'the-aco-example' };
    propagation.inject(opentelemetry.trace.setSpan(ROOT_CONTEXT, childSpan), headers);

    // Make a request to the sub path
    try {
        const { data } = await axios.get('http://localhost:8086/', { headers });
        console.log(data);
    }
    catch (error) {
        console.log('errored');
        childSpan.recordException(error);
    }

    await new Promise(res => setTimeout(res, 100));
    childSpan.end();

    return;
};