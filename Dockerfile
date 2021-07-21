FROM alpine
MAINTAINER donghongchen docker riskclient "donghongchen@shihuituan.com"
### 安装需要的软件，解决时区问题
#RUN apk --update add tzdata && \
#    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
#    echo "Asia/Shanghai" > /etc/timezone && \
#    apk del tzdata && \
#    rm -rf /var/cache/apk/*
ENV TZ=Asia/Shanghai
RUN echo "http://mirrors.aliyun.com/alpine/v3.4/main/" > /etc/apk/repositories \
    && apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone
##将代码复制到容器中
COPY ./riskclient .
## 将我们的代码编译成二进制可以执行的文件，可执行文件名为 riskclient
COPY ./conf  .

EXPOSE 3351

CMD ["./riskclient"]