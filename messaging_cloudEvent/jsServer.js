const express = require('express');
const app = express();
const port = 8080;
const { HTTP } = require("cloudevents");

const cors = require('cors');
app.use(cors());

app.use(require('body-parser').json());

app.post("/", (req, res) => {
    const receivedEvent = HTTP.toEvent({ headers: req.headers, body: req.body });
    console.log(receivedEvent);
    console.log(typeof receivedEvent);
    res.json(true);
});

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`);
});