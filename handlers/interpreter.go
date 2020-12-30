package handlers

import (
	"context"
	"net/http"
)

type QueryForm struct {
	Rule string `json:"rule"`
}

/// todo impl
func InterpretHandler(ctx context.Context, w http.ResponseWriter, args interface{}) error {
	_ = w
	_ = args.(*QueryForm)

	return nil
}
