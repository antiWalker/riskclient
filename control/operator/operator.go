package operator

import "riskengine/core"

// return nil on error
func OpInfo(op core.AOT) *core.OpInfo {
	return core.OperatorInfo(op)
}
