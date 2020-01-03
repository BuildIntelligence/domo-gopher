package domo

import (
	"reflect"
	"testing"
)

func TestNewAuthenticator(t *testing.T) {
	type args struct {
		scopes []string
	}
	tests := []struct {
		name string
		args args
		want Authenticator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthenticator(tt.args.scopes...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthenticator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticator_SetAuthInfo(t *testing.T) {
	type args struct {
		clientID     string
		clientSecret string
	}
	tests := []struct {
		name string
		a    *Authenticator
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.SetAuthInfo(tt.args.clientID, tt.args.clientSecret)
		})
	}
}

func TestAuthenticator_NewClient(t *testing.T) {
	tests := []struct {
		name string
		a    Authenticator
		want Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.NewClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authenticator.NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
