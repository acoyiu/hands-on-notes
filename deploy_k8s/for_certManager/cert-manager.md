# Cert-manager.io

![xerr](./_image/what%20is%20cert-manager.png)

cert-manager 會攜帶 Subject(如 domain 等) 資料向 "Issuer"（如 Let's encrypt， Valve 等）申請證書，然後生成 K8s 的 "Certificate" resource, 然後在以 k8s 的 certificate resource，生成 k8s 的 secret resource

- [Installation 1](https://cert-manager.io/docs/installation/)
- [Installation 2](https://github.com/cert-manager/cert-manager/releases)
- [Installation 3](https://help.aliyun.com/document_detail/409430.html), Becareful of the version
- [ACME](#acme---automatic-certificate-management-environment)
- [Type of Issuer](#type-of-issuer)
  - [ClusterIssuer (cluster-wide scoped)](#clusterissuer)
  - Issuer (namespace scoped)
- [Certificate](#certificate)
- [Ingress](#ingress-implmentation)

<br/><br/>

---

# ACME - Automatic Certificate Management Environment

是一種通信協定 (Protocol), send json 去 ACME server，server 會以

- HTTP01（通過提供計算密鑰來完成） 或
- DNS01（TXT）

方法校証 domain，然後返回證書

<br/><br/>

---

# Type of Issuer

| cert-manager 有支援幾種的 issuer type: | 描述                                                      |
| :------------------------------------- | :-------------------------------------------------------- |
| CA                                     | 使用 x509 keypair 產生certificate，存在 kubernetes secret |
| Self signed                            | 自簽 certificate                                          |
| ACME                                   | 從 ACME (ex. Let's Encrypt) server 取得 ceritificate      |
| Vault                                  | 從 Vault PKI backend 頒發 certificate                     |
| Venafi                                 | Venafi Cloud                                              |

<br/><br/>

---

## ClusterIssuer / Issuer (namespace scoped)

當您創建一個新的 ACMEIssuer時，cert-manager 將生成一個私鑰，用於在 ACME 服務器中識別您。

Creating a Basic ACME Issuer：

- Staging Cert: https://acme-staging-v02.api.letsencrypt.org/directory
- Production Cert: https://acme-v02.api.letsencrypt.org/directory

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    # 您必須將此電子郵件地址替換為您自己的, Let's Encrypt 將使用它與您聯繫有關過期的問題證書，以及與您的帳戶相關的問題。
    email: acoyiu@gmail.com

    # Staging or Production
    server: https://acme-staging-v02.api.letsencrypt.org/directory  # <- this for stainging (testing)
    # server: https://acme-v02.api.letsencrypt.org/directory        # <- this for production

    # 為 true 時，不會每次都生成俾 ACME 的身份 key, 第一次因爲沒有生成過的 key，所以多數都係 false
    disableAccountKeyGeneration: false

    # 生成俾 ACME 的身份 key 的名字, 用於被 acme server 記錄你是誰的 用途
    privateKeySecretRef:
      name: trial-issuer-account-key

    # 添加單個挑戰求解器，HTTP01 使用 nginx（ingress）
    solvers:
      - http01:
          ingress:
            class: nginx # <-- Remeber to update the ingress class name
```

```sh
# check is READY = true
kubectl get ClusterIssuer
# or
kubectl get issuer

# 檢查新創建的 k8s 的秘密
kubectl -n cert-manager get secret
# or
kubectl -n <namespaces> get issuer,secret

# 目前只有一個 tls.key 在裡面
kubectl -n cert-manager get secrets/trial-issuer-account-key -o yaml
```

<br/><br/>

---

## Certificate

- https://cert-manager.io/docs/usage/certificate/
- https://ithelp.ithome.com.tw/articles/10227274

Certificate 資源指定用於生成證書籤名請求的字段 (subject)，然後由您引用的頒發者類型完成。

**Certificate 通過指定 "certificate.spec.issuerRef" 字段來指定他們想要從哪個頒發者獲取證書。**

**Certificate K8s Obj 一旦生成，就會馬上經 issuer 找到 CA 嘗試取得頒發的證書**

```sh
kubectl create ns sandbox
```

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: trial-certificate
  namespace: sandbox
spec:
  # 秘密名稱必須
  secretName: trial-cert

  duration: 2160h # 90d
  renewBefore: 360h # 15d
  
  # SSL 裏面的 Subject field
  subject:
    organizations:
      - acokong

  # common name 字段的使用自 2000 年以來已被棄用，並且不鼓勵使用。
  commonName: yschan.info

  # 自設 CA 為 true
  isCA: false

  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 2048

  # x509 extension 入邊的 用途
  usages:
    - server auth
    - client auth
    
  # 三者中至少需要 DNS 名稱、URI 或 IP 地址其中之一，用途同 SubjectAlternativName 一樣
  dnsNames: [yschan.info]
  # uris: [spiffe://cluster.local/ns/sandbox/sa/example]
  # ipAddresses: [192.168.0.5]

  # 指定發行人（Issuer）
  issuerRef:
    kind: ClusterIssuer    # ClusterIssuer || Issuer
    name: letsencrypt-prod
```

```sh
# 檢查證書是否 READY = true
kubectl -n sandbox get certificate

# 檢查證書 K8s 對象
kubectl -n sandbox get certificate/trial-certificate -o yaml

# 檢查簽名證書, 有 tls.key + tls.crt 在裡面
kubectl -n sandbox get secret/trial-cert -o yaml

# tls.key = server's private key
# tls.crt = server's public cert
```

<br/><br/>

---

## Ingress implmentation

**需要 Ingress 的 80 和 443 能直接接收公網 request**

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    cert-manager.io/issuer: letsencrypt-prod
  name: sandbox-ingress
  namespace: sandbox
spec:
  ingressClassName: public # <------------------- Remember to update this value
  tls:
    - hosts: [yschan.info] # <------------------- Update this
      secretName: trial-cert
  rules:
    - host: yschan.info # <---------------------- Update this
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: nnginx # <----------------- Update this
                port:
                  number: 8080 # <--------------- Update this
```
