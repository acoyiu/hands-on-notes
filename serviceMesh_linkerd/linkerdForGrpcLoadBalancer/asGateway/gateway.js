require('dotenv').config();

const grpcServerLocation = process.env.SERVER_URL_LOCATION ?? 'localhost:8080';
console.log('Target gRPC server at: ', grpcServerLocation);

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// start gRPC server

const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const packageDefinition =
    protoLoader.loadSync(
        __dirname + '/../proto/chat.proto',
        {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        },
    );

const { ChatService } = grpc.loadPackageDefinition(packageDefinition);
const client = new ChatService(grpcServerLocation, grpc.credentials.createInsecure());

const MakeGrpcCall = () => new Promise((res, rej) => {
    console.log('gRPC calling to: ', grpcServerLocation);
    client.SayHello({ body: 'message to server' }, (error, inComingMsg) => {
        if (error) return rej(error);
        console.log('inComingMsg', inComingMsg);
        res(inComingMsg);
    });
});
// MakeGrpcCall(); // <- for testing

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// start expressjs server

const express = require('express');
const app = express();
const cors = require('cors');
const port = 7070;

app.use(express.json());
app.use(cors());

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

app.get('/', async (req, res) => {
    try {
        const responseFromGrpc = await MakeGrpcCall();
        res.send(responseFromGrpc);
    }
    catch (error) {
        res.status(500).send('Something broke!');
    }
});

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

app.listen(port, () => {
    console.log(`gateway listening on port ${port}`);
});

