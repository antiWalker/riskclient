package core

import (
	"bigrisk/common"
	"errors"
)

// Stack Type
type Stack []interface{}

// Len of Stack
func (stack Stack) Len() int {
	return len(stack)
}

// Cap of Stack
func (stack Stack) Cap() int {
	return cap(stack)
}

// Push new element to Stack
func (stack *Stack) Push(value interface{}) {
	//log.Debug("Push To Stack", value,"toStack")
	*stack = append(*stack, value)
}

// Pop element from Stack
func (stack *Stack) Pop() (interface{}, error) {
	theStack := *stack

	if len(theStack) == 0 {
		return nil, errors.New("Pop of stack Error: Out of index, len is 0")
	}

	value := theStack[len(theStack)-1]

	//log.Debug("Pop Of Stack", value,"ofStack")

	*stack = theStack[:len(theStack)-1]

	return value, nil
}

// Top of the Stack
func (stack *Stack) Top() (interface{}, error) {
	if len(*stack) == 0 {
		return nil, errors.New("Get top of Stack Out of index, len is 0")
	}

	common.InfoLogger.Info("Top Of Stack", (*stack)[len(*stack)-1])

	return (*stack)[len(*stack)-1], nil
}
