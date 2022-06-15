## 镜像 : alpine
```dockerfile
FROM alpine
MAINTAINER rd docker riskclient "rd@dev.com"
### 安装需要的软件，解决时区问题  
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*
##将代码复制到容器中
COPY ./riskclient .
## 将我们的代码编译成二进制可以执行的文件，可执行 文件名为 riskclient
COPY ./conf  .

EXPOSE 3351

CMD ["./riskclient"]
```

## 打tag
```docker
 docker build --rm -t 10.0.44.57:5000/risk/riskclient:v1 .
```

## 推送到仓库
```docker
docker push 10.0.44.57:5000/risk/riskclient:v1
```

## 获取镜像
```docker
docker pull 10.0.44.57:5000/risk/riskclient:v1
```

## 启动
```docker
docker run -itd --name riskclient --restart always  -p 3353:3351  -v /data/riskclient/conf:/conf 10.0.44.57:5000/risk/riskclient:v1
```
主要是以client方式来消费kafka的订单信息。
nacos redis:
redis.host = 127.0.0.1
redis.port = 6379