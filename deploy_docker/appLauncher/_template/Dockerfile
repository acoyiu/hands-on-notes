# Golang Multi-Stage Build

FROM golang:1.18.1 AS builder

COPY ./ /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

FROM alpine:latest  

COPY --from=builder /app /production

WORKDIR /production

EXPOSE 3000

CMD ["./app"]  