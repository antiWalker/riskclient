FROM alpine
MAINTAINER donghongchen docker riskclient "donghongchen@shihuituan.com"
### 安装需要的软件，解决时区问题
#RUN apk --update add curl bash tzdata && \ rm -rf /var/cache/apk/*
###修改镜像为东八区时间
#ENV TZ Asia/Shanghai
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*
##将代码复制到容器中
COPY ./riskclient .
## 将我们的代码编译成二进制可以执行的文件，可执行文件名为 riskclient
COPY ./conf  .

EXPOSE 3351

CMD ["./riskclient"]