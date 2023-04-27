# fluentd

## 1: Local file's log --> fluentd's file log

### delete if previous ran

```sh
rm -R ./1_local_files_to_logs/output
```

### Start sample 1

```sh
docker-compose -f ./1_local_files_to_logs/docker-compose.yml up -d
```

### Add below json line into [example.log](./1_local_files_to_logs/input/example.log) file

```json
// line break is Important for fluentd to know it's a complete log
{ "app": "file-app" }
```

Check real-time log message from [export-to-file](./1_local_files_to_logs/output/) directory

<br/>
<br/>

---

<br/>
<br/>

## 2: Listen Port --> fluentd's file log

### delete if previous ran

```sh
rm -R ./2_listen_port_to_logs/output
```

### Start sample 2

```sh
docker-compose -f ./1_local_files_to_logs/docker-compose.yml up -d
```

### Add log to fluentd agent by http request

```sh
curl -X POST -d 'json={"foo":"bar"}' http://localhost:9880/http-myapp.log
```

view logs in [.2_listen_port_to_logs/output/http.log](./2_listen_port_to_logs/output/http.log) file

<br/>

---

<br/>

## 3: Listen Port --> fluentd's file log + Elastic

### 1 start elastic v8.1.2 seperately

```powershell
# If using wsl
wsl -d docker-desktop
sysctl -w vm.max_map_count=262144
```

```sh
docker-compose -f ./3_logs_to_elastic/docker-compose.yml up -d elasticsearch kibana

# Check is elastic and kibana up
docker logs -f flu-3-elastic
docker logs -f flu-3-kibana
```

view Kibana at [localhost:5601](http://localhost:5601)
 
### 2 start fluentd server

```sh
# start fluentd
docker-compose -f ./3_logs_to_elastic/docker-compose.yml up -d --build fluentd

# check is fluentd working
docker logs -f flu-3-http-request
```

### Add log to fluentd agent by http request

```sh
curl -X POST -d 'json={"foo":"bar"}' http://localhost:9880/http-myapp.log
```

<br/>

---

<br/>

Delete after finish:

```sh
docker-compose -f ./1_local_files_to_logs/docker-compose.yml down
docker-compose -f ./2_listen_port_to_logs/docker-compose.yml down
docker-compose -f ./3_logs_to_elastic/docker-compose.yml down
```
