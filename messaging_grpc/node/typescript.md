# Typescript gRPC

```sh
# linux
apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+

# mac
brew install protobuf
protoc --version
```

## Inside the node project

```sh
# npm package needed
npm install ts-protoc-gen

# Path to this plugin
PROTOC_GEN_TS_PATH="./node_modules/.bin/protoc-gen-ts"

# Protobuf directory
INPUT_DIR="./proto"
PROTO_FILE="email.proto"

# Build
protoc --plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
    --js_out="import_style=commonjs,binary:${INPUT_DIR}" \
    --ts_out="service=grpc-node,mode=grpc-js:${INPUT_DIR}" \
    "${INPUT_DIR}/${PROTO_FILE}"
```