FROM scratch

MAINTAINER donghongchen docker riskclient "donghongchen@shihuituan.com"

##将代码复制到容器中
COPY ./riskclient .
## 将我们的代码编译成二进制可以执行的文件，可执行文件名为 riskclient
COPY ./conf  .

EXPOSE 3351

CMD ["./riskclient"]