const express = require('express');
const app = express();
const port = 3000;

app.get('/', async (_, res) => {
    try {
        res.send('Hello World!, message from go-env: ' + await (await fetch('http://go-env:3000/')).text());
    } catch (error) {
        res.send('Not Hello World!');
    }
});

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`);
});