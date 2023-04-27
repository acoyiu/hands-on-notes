# Start Jaeger

```sh
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.36

# or

docker run -d -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest

# or

# If using Zipkin, endpoint would be: http://localhost:9411/api/v2/spans
docker run -d -p 9411:9411 --name zipkin openzipkin/zipkin
```

<br/>

---

<br/>

## Start of example project in 3 different terminals

```sh
# 1: start golang service with new terminal
cd ./service_go && go mod tidy && go run main.go

# 2: start node service
cd ./service_node && npm install && node ./index.js

# 3: start web service
cd ./service_web && npm install && npm run dev
```

- open [localhost:8085](http://localhost:8085) with debug console before sending initial request

- then open [localhost:16686](http://localhost:16686/search) to view the trace

<br/>

---

<br/>

## Filter by Tag

- Service => service-node
- Tag => internal.span.format=jaeger

<br/>

---

<br/>

```sh
# delete all container
docker kill jaeger-backend && docker rm jaeger-backend
```
