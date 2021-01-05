module handlers

go 1.15

replace (
	riskengine/control v0.0.0 => ../control
	riskengine/core v0.0.0 => ../core
	riskengine/models v0.0.0 => ../models
	riskengine/common v0.0.0 => ../common
)

require (
	gitlaball.nicetuan.net/wangjingnan/golib v0.0.0
	riskengine/control v0.0.0
	riskengine/core v0.0.0
	riskengine/common v0.0.0
)
