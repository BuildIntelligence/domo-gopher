package domo

import (
	"context"
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

const (
	// TokenURL is the Authentication token URL for Client Credential OAuth2
	TokenURL = "https://api.domo.com/oauth/token"
	// Base Address for Domo API
	baseAddress = "https://api.domo.com"
)

const (
	// ScopeData marks OAuth2 Scope for dataset and stream APIs.
	ScopeData = "data"
	// ScopeUser marks OAuth2 Scope for user/group APIs.
	ScopeUser = "user"
	// ScopeAudit marks OAuth2 Scope for Activity Logs.
	ScopeAudit = "audit"
	// ScopeDashboard marks OAuth2 Scope for Pages API.
	ScopeDashboard = "dashboard"
	ScopeAccount = "account"
	ScopeBuzz = "buzz"
	ScopeWorkflow = "workflow"
)

// Authenticator makes it easy to configure authentication
// for Domo's client credential OAuth2 flow.
// You should use `NewAuthenticator` to make them.
//
// Example:
//
// a := domo.NewAuthenticator(domo.ScopeData, domo.ScopeUser)
// client := a.NewClient()
//
type Authenticator struct {
	config  *clientcredentials.Config
	context context.Context
}

// NewAuthenticator creates an authenticator to handle client crediential
// OAuth2 flow. Scopes should match the permissions on you credentials and the
// scopes for the API endpoints you're interacting with.
//
// By default, NewAuthenticator pulls client ID and secret from the DOMO_CLIENT_ID
// and DOMO_SECRET environment variables. If you'd like to provide them from a different
// source, call `SetAuthInfo(id, secret)` on the returned authenticator.
func NewAuthenticator(scopes ...string) Authenticator {
	cfg := &clientcredentials.Config{
		ClientID:     os.Getenv("DOMO_CLIENT_ID"),
		ClientSecret: os.Getenv("DOMO_SECRET"),
		TokenURL:     TokenURL,
		Scopes:       scopes,
	}

	ctx := context.Background()
	return Authenticator{
		config:  cfg,
		context: ctx,
	}
}

// SetAuthInfo overwrites the ClientID and ClientSecret used by the authenticator.
// Use this to provide them from a source other than the default.
func (a *Authenticator) SetAuthInfo(clientID, clientSecret string) {
	a.config.ClientID = clientID
	a.config.ClientSecret = clientSecret
}

// NewClient creates a Client that will be used for Domo API Requests.
func (a Authenticator) NewClient() *Client {
	client := a.config.Client(a.context)
	return NewClient(client)
}
