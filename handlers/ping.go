package handlers

import (
	"context"
	"net/http"
)

func IndexHandler(ctx context.Context, w http.ResponseWriter, args interface{}) error {
	_ = args

	_, _ = w.Write([]byte("ok"))

	return nil
}

func PingHandler(ctx context.Context, w http.ResponseWriter, args interface{}) error {
	_ = args // intended

	_, _ = w.Write([]byte("pong"))

	return nil
}
