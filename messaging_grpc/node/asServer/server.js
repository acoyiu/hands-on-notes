const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

/** 
 * 1: Load Protobuf in, As "Definition"
 * 2: Generate gRPC object by definition -- 將 protobuf 變成 Js object
 * 3: Extract protobuf content by KeyName in gRPC Obj
 * 4: Create new grpc.Server() with function implementation, just like express()
 * 5: start gRPC server, like "app.listen" in Express
 */

// 1 - load the protobuf file with param as 'Definition'
const packageDefinition =
    protoLoader.loadSync(
        __dirname + '/../chat.proto', // <- need use absolute path if not under
        {
            keepCase: false,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        },
    );

// 2 - generate gRPC object by definition
const grpcObject = grpc.loadPackageDefinition(packageDefinition);

// 3 - GrpcObject contains all the data struct and service declared in proto
const { ChatService, Dialogue } = grpcObject;
console.log('Dialogue', Dialogue.type.field);


const server = new grpc.Server();


// 4 - Create new grpc.Server(), just like express()
// with Implementation of service's rpc call
server.addService(

    // Definition
    ChatService.service,
    /*
    ChatService.service = {
        SendSms: { ... }
    }
    */

    // Implementation of each RPC call
    {
        SayHello: (call, callback) => {
            console.log('call: ', call);

            const dialogue = call.request; // call.request = Dialogue in Protobuf
            const { body } = dialogue;

            const returnDialogue = { body: `Returned: message from client: ${body}` };
            
            callback(null, returnDialogue);
        }
    },
);


// 5 - start gRPC server
server.bindAsync(
    '0.0.0.0:8080', // 0.0.0.0 for connection from all domain/ip
    grpc.ServerCredentials.createInsecure(),
    (error, port) => {
        if (error) throw error;
        console.log(`Server running at http://0.0.0.0:${port}`);
        server.start();
    }
);