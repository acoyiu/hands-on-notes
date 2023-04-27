# Docker Example Script

## BusyBox shell

```sh
# getting gateway IP, get docker ip
docker run -it --rm --net host busybox sh -c "ifconfig"

# run hangging node and docker with on bash
docker run -dit --name node-dev node bash -c "tail -f /dev/null" 
docker run -dit --name docker-dev docker sh -c "tail -f /dev/null"
docker run -dit --name busy-dev busybox sh -c "tail -f /dev/null"
docker run -dit --name ubun-dev ubuntu sh -c "tail -f /dev/null"

# apline container with sshpass
docker run -dit --name alpine-dev alpine sh -c "tail -f /dev/null"
docker exec -it alpine-dev sh
apk add rsync
apk add --update openssh
apk add --update --no-cache openssh sshpass

# mount Docker Deamon
docker run -dit --name docker-dev -v $(pwd):/app -v /var/run/docker.sock:/var/run/docker.sock --privileged docker:dind sh -c "tail -f /dev/null"
```

---

## Parse

要用 -v 將 cloud code 同 config.json 接入 container 裡面 （folder remapping）

```sh
# 創建 by using config.json !! *****
docker run --name parseser -dit -p 1337:1337 -v $PWD/cloud:/parse-server/cloud -v $PWD/config:/parse-server/config parseplatform/parse-server /parse-server/config/config.json

# 創建 by using appid, masterkey ( no good )
docker run --name parseser -dit -p 1337:1337 -v hostPath:/parse-server/cloud -v hostPath:/parse-server/config parseplatform/parse-server --appId myAppId --masterKey myMasterKey --databaseURI mongodb://172.17.0.2:27017/dev

# default Parse
npm start -- --appId myAppId --masterKey myMasterKey --databaseURI mongodb://localhost:27017/dev
npm start -- path/to/config.json

# default dashboard
docker run -dit -p 4040:4040 --name dashser -v $PWD/config/dashboard.json:/src/Parse-Dashboard/parse-dashboard-config.json parseplatform/parse-dashboard --config /src/Parse-Dashboard/parse-dashboard-config.json --allowInsecureHTTP
docker run -dit -p 4040:4040 --name dashser -v $PWD/dashboard.json:/dashboard.json parseplatform/parse-dashboard --config /dashboard.json --allowInsecureHTTP

# docker run -dit --name dashser -p 4040:4040 -v $PWD/parse-dashboard-config.json:/json/parse-dashboard-config.json node
# docker exec -it dashser bash
# npm install -g parse-dashboard
# parse-dashboard --config /json/parse-dashboard-config.json --allowInsecureHTTP
```

---

## gitea

```sh
docker run -d --name=gitea -p 10022:22 -p 10080:3000 -v /var/lib/gitea:/data gitea/gitea:latest
docker run -dit --name gitser -p 10022:22 -p 3000:3000 -v $PWD/gitdata:/data gitea/gitea:latest  
```

---

## apache

```sh
docker run -p 7000:80 -dit --name apacheser -v $PWD/htdocs:/usr/local/apache2/htdocs -v $PWD/httpd.conf:/usr/local/apache2/conf/httpd.conf httpd
# /usr/local/apache2/conf/httpd.conf

# php with apache
docker run -p 7000:80 -dit --name phpser -v $PWD/html:/var/www/html php:apache
docker run -p 7000:80 -dit --name phpser -v $PWD/html:/var/www/html -v $PWD/apache2:/etc/apache2 php:apache
docker run -p 7000:80 -dit --name phpser -v $PWD/html:/var/www/html -v $PWD/apache2.conf:/etc/apache2/apache2.conf php:apache

# mysqli extension
docker exec -ti <your-php-container> sh
>> docker-php-ext-install mysqli
>> docker-php-ext-enable mysqli
>> apachectl restart
```

---

## nginx

```sh
docker run --name n7 -p 8080:80 -v $PWD:/usr/share/nginx/html:ro -d --restart unless-stopped nginx
docker run --name n7 -p 8080:80 -v $PWD:/usr/share/nginx/html:ro -v $PWD/default.conf:/etc/nginx/conf.d/default.conf:ro -d --restart unless-stopped nginx

# allow directory listing
docker run --name n7 -p 8080:80 -v $PWD:/usr/share/nginx/html:ro -v $PWD/listindex.conf:/etc/nginx/conf.d/default.conf:ro -d nginx

## REverse router
docker run --name n7 -p 80:80 -p 443:443 -v $PWD/default.conf:/etc/nginx/conf.d/default.conf:ro -v $PWD/ssl:/ssl --restart unless-stopped -dit nginx
```

---

## mysql

```sh
docker run -p 3306:3306 --name sqlser -e MYSQL_ROOT_PASSWORD=root -v $PWD/mysql/conf:/etc/mysql/conf.d -v $PWD/mysql/data:/var/lib/mysql -dit mysql:5.7
docker exec -it sqlser mysql -u root -p

# docker exec -it sqlser mysql(唔係bash，係mysql指令) -u（用戶） root（root用戶） -p（password）
# sql admin ser
docker run -p 6060:80 --name adminser -e PMA_HOST=192.168.43.37 -e PMA_PORT=3306 --link sqlser -dit phpmyadmin/phpmyadmin
```

---

## node installer docker & runner

```sh
docker run -it --name node-installer -v $PWD:/install node bash -c "cd /install && npm i" && docker rm node-installer
docker run -dit --name node-runner -v $PWD:/project -p 9990:8080 node bash -c "cd /project && npm start"
```

---

# Redis

```sh
docker run --name some-redis -dit -v $PWD/redisData:/data -p 6379:6379 --restart unless-stopped redis redis-server --appendonly yes
docker exec -it some-redis bash
redis-cli
SET asd "wrgwe"
GET asd

## external call
brew install redis
redis-cli -h localhost -p 6379 ping

### -> PONG
```

---

## swagger

```sh
docker run -d \
-p 5000:8080 \
--platform linux/amd64 \
--name some-swagger \
swaggerapi/swagger-editor
```
