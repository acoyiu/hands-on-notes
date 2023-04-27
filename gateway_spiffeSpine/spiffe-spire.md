# Spiffe / Spine

- ### Access Control = Who can do What = Policy
- ### Spine Server/Agent 都有個 API Socket（Unix）file

Indexes:

- [Installation](https://github.com/spiffe/spire/releases)
- command: [spire-server](#spire-server)
- command: [spire-agent](#spire-agent)
- Tutor: [first](#first-tutor)
- [List of available method to register workload](https://spiffe.io/docs/latest/deploying/registering/)

<br/><br/>

---

## Installation

```sh
# get the executable:
curl -s -N -L https://github.com/spiffe/spire/releases/download/v1.4.4/spire-1.4.4-linux-x86_64-glibc.tar.gz | tar xz
```

<br/><br/>

---

## First Tutor

```sh
# First get the project
git clone --single-branch --branch v1.3.3 https://github.com/spiffe/spire.git
cd spire

#  show the configs
cat ./conf/server/server.conf
```

### Example config for server

Config = server (x509 subjects) + plugin config

```config
server {
    bind_address = "127.0.0.1"
    bind_port = "8081"
    socket_path = "/tmp/spire-server/private/api.sock"
    trust_domain = "example.org"
    data_dir = "./.data"
    log_level = "DEBUG"
    ca_subject {
        country = ["US"]
        organization = ["SPIFFE"]
        common_name = ""
    }
}

plugins {
    DataStore "sql" {
        plugin_data {
            database_type = "sqlite3"
            connection_string = "./.data/datastore.sqlite3"
        }
    }

    NodeAttestor "join_token" {
        plugin_data {}
    }

    KeyManager "memory" {
        plugin_data = {}
    }

    UpstreamAuthority "disk" {
        plugin_data {
            key_file_path = "./conf/server/dummy_upstream_ca.key"
            cert_file_path = "./conf/server/dummy_upstream_ca.crt"
        }
    }
}
```

```sh
# start server
spire-server run -config ./conf/server/server.conf &

# check health
spire-server healthcheck
```

### 生成一次性令牌以用於證明代理

```sh
spire-server token generate -spiffeID spiffe://example.org/myagent
# Token: 847127b4-daa5-4bbe-b5b1-1d18150193b4
```

### 啟動 SPIRE 代理

```sh
# Agent config
cat ./conf/agent/agent.conf
```

```config
agent {
    data_dir = "./.data"
    log_level = "DEBUG"
    server_address = "127.0.0.1"
    server_port = "8081"
    socket_path ="/tmp/spire-agent/public/api.sock"
    trust_bundle_path = "./conf/agent/dummy_root_ca.crt"
    trust_domain = "example.org"
}

plugins {
    NodeAttestor "join_token" {
        plugin_data {
        }
    }
    KeyManager "disk" {
        plugin_data {
            directory = "./.data"
        }
    }
    WorkloadAttestor "k8s" {
        plugin_data {
            kubelet_read_only_port = "10255"
        }
    }
    WorkloadAttestor "unix" {
        plugin_data {
        }
    }
    WorkloadAttestor "docker" {
        plugin_data {
        }
    }
}
```

```sh
# start agent
spire-agent run -config ./conf/agent/agent.conf -joinToken <token_string> &

# Check agent health
spire-agent healthcheck
```

### Create a registration policy for your workload

為您的工作負載註冊

為了讓 SPIRE 識別工作負載，您必須通過註冊條目向 SPIRE 服務器註冊工作負載。 工作負載註冊告訴 SPIRE 如何識別工作負載以及為其提供哪個 SPIFFE ID。

```sh
# 此命令正在根據當前用戶的 UID ($(id -u)) 創建一個註冊條目
# 如有必要，請隨意調整
spire-server entry create \
  -parentID spiffe://example.org/myagent \
  -spiffeID spiffe://example.org/myservice \
  -selector unix:uid:$(id -u)
```

```sh
# Example Response
Entry ID         : f9ce9ed1-d680-4a33-bbb0-7843b3727503
SPIFFE ID        : spiffe://example.org/myservice
Parent ID        : spiffe://example.org/myagent
Revision         : 0
TTL              : default
Selector         : unix:uid:1000
```

### Retrieve and view a x509-SVID

檢索和查看 x509-SVID

```sh
spire-agent api fetch x509 -write /tmp/

openssl x509 -in /tmp/svid.0.pem -text -noout
```

<br/><br/>

---

## spire-server

| Common in server and agent | Usage                     |
| :------------------------- | :------------------------ |
| **healthcheck**            | 確定服務器健康狀態        |
| **run**                    | 運行服務器                |
| **validate**               | 驗證 SPIRE 服務器配置文件 |

| Sub-command    | Usage                                     |
| :------------- | :---------------------------------------- |
| **federation** |                                           |
| - create       | 與外部信任域創建動態聯合關係              |
| - delete       | 刪除動態聯合關係                          |
| - list         | 列出所有動態聯合關係                      |
| - refresh      | 刷新來自指定聯合信任域的捆綁包            |
| - show         | 顯示動態聯合關係                          |
| - update       | 更新與外部信任域的動態聯合關係            |
| **agent**      |                                           |
| - ban          | 根據 SPIFFE ID 禁止經過認證的代理         |
| - count        | 獲取總數認證代理                          |
| - evict        | 根據 SPIFFE ID 驅逐經過認證的代理         |
| - list         | 列出經過認證的代理及其 SPIFFE ID          |
| - show         | 顯示給定 SPIFFE ID 的已證明代理的詳細信息 |
| **bundle**     |                                           |
| - count        | 計數捆綁                                  |
| - delete       | 刪除捆綁數據                              |
| - list         | 列出聯合捆綁數據                          |
| - set          | 創建或更新捆綁數據                        |
| - show         | 將服務器 CA 包打印到標準輸出              |
| **entry**      |                                           |
| - count        | 統計註冊條目                              |
| - create       | 創建註冊條目                              |
| - delete       | 刪除註冊條目                              |
| - show         | 顯示配置的註冊條目                        |
| - update       | 更新註冊條目                              |

| Asset-command  | Usage                  |
| :------------- | :--------------------- |
| jwt            | JWT Util               |
| jwt mint       | 鑄造一個 JWT-SVID      |
| token          | token Util             |
| token generate | Generates a join token |
| x509           | x509 Util              |
| x509 mint      | 鑄造一個 X509-SVID     |

<br/><br/>

---

## spire-agent

| Sub-command     | Usage                   |
| :-------------- | :---------------------- |
| **healthcheck** | 確定代理健康狀態        |
| **run**         | 運行代理                |
| **validate**    | 驗證 SPIRE 代理配置文件 |
| api             | 未知                    |

<br/><br/>

---

# Spire MTLS

- https://spiffe.io/docs/latest/spire-about/use-cases/#authenticating-workloads-in-untrusted-networks-using-mtls
- https://spiffe.io/docs/latest/microservices/envoy/