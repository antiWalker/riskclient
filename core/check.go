package core

import (
	"context"
	"errors"
)

// Check the rule
func Check(ctx context.Context,rule []byte) (bool, error) {

	var validRule = false

	if rule == nil {
		return validRule, errors.New("riskEngine: rule is empty")
	}

	_, ok := constructNodeFromString(ctx,rule)

	if ok != nil {
		return validRule, errors.New("riskEngine: construct cooked rule failed" + ok.Error())
	}

	return true, nil
}
