package handlers

import (
	"net/http"
	"testing"
)

func TestUrlsHandler_Handle(t *testing.T) {
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name    string
		handler *UrlsHandler
		args    args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.handler.Handle(tt.args.res, tt.args.req)
		})
	}
}
