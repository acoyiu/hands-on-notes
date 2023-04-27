## Statuful Set for MongoDB Example

### 1 Create K8s Resources

**Headless server act as a manager to help provide stateful pod static endpoint**

```yaml
# mongo-headless-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mongo-svc
  namespace: mongo
spec:
  ports:
    - port: 27017
      targetPort: 27017
  clusterIP: None
  selector:
    role: mongo

---
# statefulset-mongo.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mongo
  namespace: mongo
spec:
  selector:
    matchLabels:
      role: mongo
  serviceName: mongo-svc
  replicas: 3
  template:
    metadata:
      labels:
        role: mongo
    spec:
      terminationGracePeriodSeconds: 10
      containers:
        - name: mongo
          image: mongo:5
          command:
            - "mongod"
            - "--auth"
            - "--keyFile=/common/keyfile"
            - "--bind_ip_all"
            - "--replSet=MainRepSet"
          ports:
            - containerPort: 27017
          volumeMounts:
            - name: mongo-persistent-storage
              mountPath: /data/db
            - name: cm-common-key
              mountPath: /common
      volumes:
        - name: cm-common-key
          configMap:
            name: common-key
            defaultMode: 0400
  volumeClaimTemplates:
    - metadata:
        name: mongo-persistent-storage
      spec:
        accessModes: ["ReadWriteMany"]
        storageClassName: nfs-client
        resources:
          requests:
            storage: 1Gi

---
# config map is for enable mongo auth mode

# openssl rand -base64 756 > <path-to-keyfile>
# chmod 400 <path-to-keyfile>

apiVersion: v1
kind: ConfigMap
metadata:
  name: common-key
  namespace: mongo
data:
  keyfile: |-
    UIKFx5f3y1a6qr/AgnC1cHNkZOLPCdLiIFcE5M9pjYJEjB1r3mgcK/k3HSIhVj0U
    C/7oXAOicxzSqwzHWcWIGxg6SGuuYgmslEXdatQboui2ZF+gTtIYnRx3b9H/1sUb
    n0elzLPIHvM2NEgj8jPOvPBaX9Kvl07A/CEZd+oLB5niTSRBfW3Cp9iAaTM8Nqpg
    303fSL8cUshJLx2OBd1VinhFFCzBEpUqknmbqjKI2ECYClWKrK21N0DZyfsC1tqF
    qRQR9KKIXoNz/KTUVjeRyT/hiUeRz436AXHJyj2zMkfM9CVpLGIaxewm9D22g7/l
    8wQAqrxgD0ztGjdNW1wz8eK8mkXPjAXSZq9jUUpx9PvtzMjJo1D2HUIKBebR7x8T
    d7F2BWwy0iSwFh08+NTwzAlPL086MPQm8GVtsiumgzdSxmMIV0V71bFZAjO7N8MX
    0qeC1EqlNlqIjieX/XAG+7D8Yp2E6SerRSWvbzTMNzJclGDUzRizfCHyVgB1l1Tf
    7DF8fris75FcDvMs5ZQ06CtLHedL4ReNH05jaAXpwP7ZbkvCSs+cTaRoLZEem6Qj
    DZb5HYkA9dU2Ta9xOLUq7X8dgUefPTadNB6phJqgu9GM7xrv/2mGIRXDwSnyJuZB
    qfaBTE+T5eMAW5vA0Ej/VqiWmwxFRpzEhCcXnbHVP8UshF94jN9fBC04JqbxQnZq
    Riho081kcqHh76qsH6psdIyhgyuj5U8Krs+zs3UK7KdOUSj2JZ2+Pr9g53BAhvv9
    02Lr2MlcxoEZIRxWgepQp5sUkAayWRT7s8Kru0cyCcpGlE7Ew31ci4wiI5EZUOYi
    ag6d9Dc0/ydOZdBeyO7v2OFtUFSW0EPW1qe49jC6K7DZPuHL2BSyw6SsjGjSjknb
    oQrsD6LR31qyGc4SpthCSHPusdcWRGL/bQ3b/BDVgyUfuzIjxmRJiU+72kFCmVBy
    eG2eccUZJwQKH/8on7X7rAA1JIBZIZ8lrCdbTvg3YWXMiiMn
```

### 2 Initiate ReplicaSet in Mongo

```sh
# enter mongo primary (may be in mongo-1 or mongo-2)
kubectl -n mongo exec -it mongo-0 -- bash

# get hostname for making connection
hostname -f

# enter mongo shell
mongosh
```

### 3 Init mongo RS in mongo shell

```javascript
// init function
rs.initiate({
    _id: "MainRepSet",
    version: 1,
    members: [
        { _id: 0, host: "mongo-0.mongo-svc.mongo.svc.cluster.local" },
        { _id: 1, host: "mongo-1.mongo-svc.mongo.svc.cluster.local" },
        { _id: 2, host: "mongo-2.mongo-svc.mongo.svc.cluster.local" },
    ]
});
```

[Can create db amin user now](../db_mongo/mongo.md#dbcreateuser-admin)

then can connect with uri like below WITH-IN the cluster:

```sh
mongodb://username:password@mongo-2:27017,mongo-1:27017,mongo-0:27017/?authSource=db
# which will always connect to the primary
```

### Then set fixed pod to always become primary if need

```sh
// Run on Primary

cfg = rs.conf()

cfg.members[0].priority = 1
cfg.members[1].priority = 0.5
cfg.members[2].priority = 0.1

rs.reconfig(cfg)

// step down in-case required
// rs.stepDown() 
```

### 4 Expose service

#### expose fix pod

```yaml
# the nodeport service
apiVersion: v1
kind: Service
metadata:
  name: mongo-0
  namespace: mongo
spec:
  type: NodePort
  ports:
    - protocol: TCP
      nodePort: 32100
      port: 27017
      targetPort: 27017
  selector:
    role: mongo
    statefulset.kubernetes.io/pod-name: mongo-0
```

#### port forward svc (if needed)

```sh
kubectl -n mongo port-forward svc/mongo-0 27017:27017 --address=0.0.0.0
```