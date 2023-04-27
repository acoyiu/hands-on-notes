CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" main.go

docker run -dit --name alpine -v $PWD:/app alpine sh -c 'sleep infinity'

