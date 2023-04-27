const ws = new WebSocket("ws://localhost:8080");

ws.on("message", function incoming(message) {
    const event = new CloudEvent(JSON.parse(message));
    if (event.type === "weather.error") {
        console.error(`Error: ${event.data}`);
    } else {
        print(event.data);
    }
    ask();
});

function send() {
    const event = new CloudEvent({
        type: "weather.query",
        source: "/weather.client",
        data: { zip },
    });
    ws.send(event.toString());
}