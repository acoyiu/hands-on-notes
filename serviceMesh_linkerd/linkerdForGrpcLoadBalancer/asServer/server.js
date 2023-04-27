require('dotenv').config();

const port = process.env.SERVER_HOST_PORT ?? '8080';
console.log('Using port: ', port);

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const packageDefinition =
    protoLoader.loadSync(
        __dirname + '/../proto/chat.proto',
        {
            keepCase: false,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        },
    );

const grpcObject = grpc.loadPackageDefinition(packageDefinition);
const { ChatService } = grpcObject;
const server = new grpc.Server();

server.addService(
    ChatService.service,
    {
        SayHello: (call, callback) => {
            console.log('gRPC triggered remotely');
            console.log('call.request.body: ', call.request.body);
            callback(null, { body: 'message to client' });
        }
    },
);

server.bindAsync(
    `0.0.0.0:${port}`, // 0.0.0.0 for connection from all domain/ip
    grpc.ServerCredentials.createInsecure(),
    (error, _port_) => {
        if (error) throw error;
        console.log(`Server running at http://0.0.0.0:${port}`);
        server.start();
    }
);