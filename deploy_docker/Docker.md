# Docker CMDs

## For dockerize an executable, pls refer to [golang prometheus example](../../code-note/observability_prometheus/pusherGoPing/dockerize/dockerize.md)

## Login

```sh
docker login -u [username]
```

## List Docker

```sh
docker image ls
```

## Run a Container

```sh
docker run 
  [-d = detach to background] # -dit 一開始就要加，不然以後也不能 exec 進去
  [-i = interactive]
  [-t = by terminal]
  [--name = add name to the container]
  [--mount source=[volume],target=[path]] # less use
  [-v <host path>:<container path>:ro=readOnly rw=readWrite]
  [-p <outside port>:<inside port>] <- can be multiple -p tag
  [--rm = remove the container once it exits/stops]
  [-c = execute container command]  # docker run -it --rm --net host busybox sh -c "ifconfig"
  [--network = ?]
[imagename]
```

## enter container

```sh
docker exec -it containerName bash

# enter container with root user
docker exec -u 0 -it containerName bash
```

## execute container command

```sh
docker exec -it containerName bash -c "ls -al /directory"
```

## Run hangging node and docker with on bash

```sh
docker run -dit --name node-dev node bash -c "tail -f /dev/null"
```

## apline container with sshpass

```sh
docker run -dit --name alpine-dev alpine sh -c "tail -f /dev/null"
docker exec -it alpine-dev sh
apk add rsync
apk add --update openssh
apk add --update --no-cache openssh sshpass
```

## mount Docker Deamon

```sh
docker run -dit --name docker-dev -v $(pwd):/app -v /var/run/docker.sock:/var/run/docker.sock --privileged docker:dind sh -c "tail -f /dev/null"
```

## Container Operation

```sh
# List
docker ps -a

# common actions
docker stop    [container-name/hash]
docker start   [container-name/hash]
docker restart [container-name/hash]
docker rm -f [container-name/hash] # delete container without stopping it

# Batch actions
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker system prune  
docker volume prune
```

## Enter container

```sh
docker exec [-d = background] [container] bash
[-i = interactive] [container] bash
[-t = by terminal] [container] bash
```

## Remount volume into running container

```sh
# Mounting volume after container is created
docker commit 73e3ee82305e cache-image-name
docker run -it -v "$PWD/somedir":/somedir --name=name cache-image-name [cmd]
```

## show docker details statc

```sh
docker system df
```

## volumn commands

```sh
docker volume ls
rm
prune
```

## inspect specific volumn

```sh
docker inspect [volumn]

docker logs [container name]
docker logs --since 2013-01-02T13:23:37 parseser
```

## change image name

```sh
docker tag {imageID} {imageName}:{tagName}
```

## Copy fie in and out container

```sh
# The cp command can be used to copy files. One specific file can be copied like
docker cp foo.txt mycontainer:/foo.txt
docker cp mycontainer:/foo.txt foo.txt

# Multiple files contained by the folder src can be copied into the target folder using
docker cp src/. mycontainer:/target
docker cp mycontainer:/src/. target
```
