## 镜像: scratch
```
FROM scratch

MAINTAINER donghongchen docker riskclient "donghongchen@shihuituan.com"

##将代码复制到容器中
COPY ./riskclient .
## 将我们的代码编译成二进制可以执行的文件，可执行文件名为 riskclient
COPY ./conf  .

EXPOSE 3351

CMD ["./riskclient"]
```

打tag
```docker
 docker build --rm -t donghongchen/riskclient:v5 .
```

推送到仓库
```docker
docker push donghongchen/riskclient:v5
```

获取镜像
```docker
docker pull donghongchen/riskclient:v4
```

启动
```docker
docker run -itd --name riskclient --restart always  -p 3353:3351  -v /data/riskclient/conf:/conf donghongchen/riskclient:v4

```


## 镜像: golang:alpine

docker run -itd --name riskclient --restart always  -p 3354:3351  -v /data/riskclient/conf:/go/conf donghongchen/riskclient:v3


启动
```docker
docker run -itd --name riskclient --restart always  -p 3353:3351  -v /data/riskclient/conf:/go/conf donghongchen/riskclient:v4

```