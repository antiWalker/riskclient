package operator

import "bigrisk/core"

// return nil on error
func OpInfo(op core.AOT) *core.OpInfo {
	return core.OperatorInfo(op)
}
