const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');


const packageDefinition =
    protoLoader.loadSync(
        __dirname + '/../chat.proto',
        {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        },
    );


const { ChatService } = grpc.loadPackageDefinition(packageDefinition);


// 6 - new 一個 protobuf 裡面聲明過的 service
const client = new ChatService(
    'localhost:8080',
    grpc.credentials.createInsecure()
);

// 7 - 用 client 直接 call remote function
client.SayHello({ body: 'message to server' }, (error, inComingMsg) => {
    if (error) throw error;
    console.log('inComingMsg', inComingMsg);
});