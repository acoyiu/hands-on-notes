import { HTTP, CloudEvent } from "cloudevents";

async function call() {
    const ce = new CloudEvent({
        source: "/some/js/event",
        type: "example",
        data: { test: "string" },
    });

    console.log(
        await toUpstream(ce)
    );
}

async function toUpstream(ce) {
    console.log(ce);

    switch (ce.source) {
        case '/some/js/event':
            const ceGluedMsg = HTTP.binary(ce); // Or HTTP.structured(ce)
            return await fetch(
                "http://localhost:8080/",
                {
                    method: "POST",
                    body: ceGluedMsg.body,
                    headers: ceGluedMsg.headers,
                }
            );
        case '/some/kafka/event':
            return 'kafka binding'
    }
}