module bigrisk

go 1.15

replace (
	golang.org/x/sys v0.3.0 => github.com/golang/sys v0.3.0
	riskengine/common v0.0.0 => ./common
	riskengine/control v0.0.0 => ./control
	riskengine/core v0.0.0 => ./core
	riskengine/handlers v0.0.0 => ./handlers
	riskengine/models v0.0.0 => ./models
)

require (
	github.com/Shopify/sarama v1.27.2
	github.com/astaxie/beego v1.12.3
	github.com/sirupsen/logrus v1.7.0
	gitlaball.nicetuan.net/wangjingnan/golib v0.0.0
	riskengine/common v0.0.0
	riskengine/handlers v0.0.0
	riskengine/models v0.0.0
)
