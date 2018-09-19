// Package domo provides utilities for easily interfacing
// with Domo's APIs.
package domo

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    Error
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_decodeError(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.decodeError(tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Client.decodeError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_shouldRetry(t *testing.T) {
	type args struct {
		status int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRetry(tt.args.status); got != tt.want {
				t.Errorf("shouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isFailure(t *testing.T) {
	type args struct {
		code       int
		validCodes []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFailure(tt.args.code, tt.args.validCodes); got != tt.want {
				t.Errorf("isFailure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_execute(t *testing.T) {
	type args struct {
		req         *http.Request
		result      interface{}
		needsStatus []int
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.execute(tt.args.req, tt.args.result, tt.args.needsStatus...); (err != nil) != tt.wantErr {
				t.Errorf("Client.execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_retryDuration(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := retryDuration(tt.args.resp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("retryDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_get(t *testing.T) {
	type args struct {
		url    string
		result interface{}
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.get(tt.args.url, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("Client.get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}