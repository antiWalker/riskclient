module bigrisk

go 1.15

replace golang.org/x/sys v0.3.0 => github.com/golang/sys v0.3.0

require (
	github.com/Shopify/sarama v1.27.2
	github.com/astaxie/beego v1.12.3
	github.com/go-redis/redis v6.14.2+incompatible
	github.com/sirupsen/logrus v1.7.0
	gitlaball.nicetuan.net/wangjingnan/golib v1.1.3
	go.uber.org/zap v1.16.0
)
