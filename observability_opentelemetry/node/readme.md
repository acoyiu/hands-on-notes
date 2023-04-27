## 1: Basic trace generation

```sh
node ./steps/1_basic_create_span.js
```

## 2: Nested trace(span) generation

```sh
node ./steps/2_basic_create_nested_span.js
```

## 3: Basic export to Jaeger

```sh
node ./steps/3_export_to_jaeger.js
```

## 3: Pass context between expressjs route

```sh
node ./steps/4_0_manual_instrument.js
curl localhost:3000
```

## 4: Pass context between different expressjs service

```sh
node ./steps/4_1_manual_instrument.js
node ./steps/4_2_manual_instrument.js
curl localhost:3000
```

## 5: Auto instrumentation (no SDK)

```sh
node ./steps/5_auto_instrument_no_sdk.js
curl localhost:3000
```

## 6: Auto instrumentation (with SDK)

```sh
node ./steps/6_auto_instrument_by_sdk.js
curl localhost:3000
```

<br/>

---

<br/>

## JavaScript SDK

```javascript
const opentelemetry = require("@opentelemetry/api");
const provider = new BasicTracerProvider();
```

### register tracer

```javascript
const tracer_registration = (tracer_get_name) => {
  const {
    BasicTracerProvider,
    ConsoleSpanExporter,
    SimpleSpanProcessor,
  } = require("@opentelemetry/sdk-trace-base");

  // Configure span processor to send spans to the exporter
  provider.addSpanProcessor(new SimpleSpanProcessor(new ConsoleSpanExporter()));
  provider.register();

  // This is what we'll access in all instrumentation code
  return opentelemetry.trace.getTracer(tracer_get_name ?? "basic-tracer");
};
```

### Create root span

```javascript
const parent_span = tracer.startSpan("main");
```

### Create child span

```javascript
const create_child_span = (
  trace_name,
  attribute_object = undefined,
  parent_span
) => {
  // Start another span. In this example, the main function already started a
  // span, so that'll be the parent span, and this will be a child span.

  // create context that inherit all info from parent
  const ctx = opentelemetry.trace.setSpan(
    opentelemetry.context.active(),
    parent_span
  );

  // inject parent's context into child process
  const span = tracer.startSpan(trace_name, attribute_object, ctx);

  return span;
};
```

### Get current span

```javascript
const span = opentelemetry.trace.getSpan(opentelemetry.context.active());
// do something with the current span, optionally ending it if that is appropriate for your use case.
```

### Add attribute to span

```javascript
// way 1 - add attribute when create span
const span = tracer.startSpan(
  "trace_name",
  { attributes: { attribute1: "value1" } },
  ctx
);

// way 2 - add anytime an attribute to the same span later on
span.setAttribute("attribute2", "value2");

// use sematic attribute
// npm install --save @opentelemetry/semantic-conventions
const { SemanticAttributes } = require("@opentelemetry/semantic-conventions");
span.setAttribute(SemanticAttributes.CODE_FILEPATH, __filename);
```

### Span events - An event is a human-readable message attached to a span

```javascript
span.addEvent("Doing something");
const result = doWork();
span.addEvent("Did something");

// You can also add an object with more data to go along with the message:
span.addEvent("some log", {
  "log.severity": "error",
  "log.message": "Data not found",
  "request.id": requestId,
});
```

### Span Record Exception and Status (to show Error)

```javascript
try {
  doWork();
} catch (err) {
  span.recordException(err);
  // will shwo as log

  // will show as tag
  span.setStatus({
    code: opentelemetry.SpanStatusCode.ERROR,
    message: err,
  });
}
```

### Get context inherit from other service if there is one (like proxy)

```javascript
// get inherit span
const ctx = opentelemetry.trace.setSpan(opentelemetry.context.active(), undefined);

// apply to currently working one
const mySpan = tracer.startSpan('nodejs-root-calling', undefined, ctx);

console.log(ctx); ->
BaseContext {
  _currentContext: Map(0) {},
  getValue: [Function (anonymous)],
  setValue: [Function (anonymous)],
  deleteValue: [Function (anonymous)]
}
```
