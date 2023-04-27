## for making project

```sh
go mod init <name>

go get \
  go.opentelemetry.io/otel \
  go.opentelemetry.io/otel/attribute \
  go.opentelemetry.io/otel/trace \
  go.opentelemetry.io/otel/sdk/resource \
  go.opentelemetry.io/otel/sdk/trace \
  go.opentelemetry.io/otel/exporters/jaeger \
  github.com/gin-contrib/cors \
  github.com/gin-gonic/gin
```

## for using project

```sh
# npm install
go mod tidy

go run main.go
# run multipls times, then check jaeger ui

go run main.go --propa
# then curl localhost:6060/c01

go run main.go --grpc
# then curl localhost:6060
```
