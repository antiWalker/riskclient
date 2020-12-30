module bigrisk

go 1.15

replace (
	golang.org/x/sys v0.3.0 => github.com/golang/sys v0.3.0
	riskengine/control v0.0.0 => ./control
	riskengine/core v0.0.0 => ./core
	riskengine/handlers v0.0.0 => ./handlers
	riskengine/models v0.0.0 => ./models
)

require (
	github.com/sirupsen/logrus v1.7.0
	gitlaball.nicetuan.net/wangjingnan/golib v0.0.0
	riskengine/handlers v0.0.0
)
