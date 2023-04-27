const { httpTransport, emitterFor, CloudEvent } = require("cloudevents");

// Create an emitter to send events to a reciever
const emit = emitterFor(httpTransport("https://my.receiver.com/endpoint"));

// Create a new CloudEvent
const ce = new CloudEvent({ type, source, data });

// Send it to the endpoint - encoded as HTTP binary by default
emit(ce);