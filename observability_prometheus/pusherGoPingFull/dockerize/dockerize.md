# Dockerize

## Build Image

```sh
# Build executable
cd ../ && CGO_ENABLED=0 GOOS=linux go build

# copy etable if need
cp ./ping_metrics_to_gateway ./dockerize/etable
cd dockerize

# build docker img
docker build -t ping-pusher .

# try push
docker run -it --rm --net host --name ping-pusher ping-pusher http://172.30.166.200:9091 instance_name job_name http://123.57.136.251 https://google.com
```

## Push to Registry

```sh
DIRECTORY="./"
REGISTRY=registry.greatics.net
IMGNAME=ping-pusher
TAGNAME=

./imageBuildPusher_noLogin.sh $DIRECTORY $REGISTRY $IMGNAME $TAGNAME
```
