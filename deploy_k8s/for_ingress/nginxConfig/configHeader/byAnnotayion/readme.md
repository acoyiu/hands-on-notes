### 所有 Nginx 的 config 都可以用 ingress yaml 裏面的 annotation 的 nginx.ingress.kubernetes.io/*** 去設定的！

---

## Available annotation
https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/

---


## extra info
https://www.cxymm.net/article/catoop/114013172



## attention
要注意官方文档中开篇的一句话：
!!! note The annotation prefix can be changed using the --annotations-prefix command line argument,
but the default is nginx.ingress.kubernetes.io, as described in the table below.



## Way to insert config
nginx.ingress.kubernetes.io/server-snippet （扩展配置到 server 块中的代码段）
nginx.ingress.kubernetes.io/configuration-snippet （扩展配置到 location 块代码段）



## example
```
annotations:

    nginx.ingress.kubernetes.io/server-snippet: |
        add_header aco-custom-header customheadervalue;
        more_set_headers "aco-id: acsff";

    nginx.ingress.kubernetes.io/configuration-snippet: |
        proxy_set_header X-Aco-Real-IP $remote_addr;

    nginx.ingress.kubernetes.io/rewrite-target: /$2 
```
### into
```
server {
    # server-snippet 配置 :: Add header to Response's header
    add_header aco-custom-header customheadervalue;
    more_set_headers "aco-id: acsff";
    
    location / {
        # configuration-snippet配置 :: Adding header to proxy request
        proxy_set_header X-Aco-Real-IP $remote_addr;
    }
}
```