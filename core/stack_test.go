package core

import (
	"testing"
)

func TestStack_Len(t *testing.T) {
	var myStack Stack
	myStack.Push(1)
	myStack.Push("test")
	if myStack.Len() == 2 {
		t.Log("Pass Stack.Len")
	} else {
		t.Error("Failed Stack.Len")
	}
}
