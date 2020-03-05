package domo

import (
	"context"
	"fmt"
	"net/http"
)

// AccountssService handles communication with the Accounts
// related methods of the Domo API.
//
// Domo API Docs: https://developer.domo.com/docs/accounts-api-reference/account-api-reference
type AccountsService service

type Account struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name:omitempty"`
	AccountType DomoAccountType
}

type DomoAccountType struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name, omitempty"`

}

func (s *AccountsService) List(ctx context.Context) (*http.Response, error) {
	u := "v1/accounts"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AccountsService) Info(ctx context.Context, accountID string) (*http.Response, error) {
	u := fmt.Sprintf("v1/accounts/%s", accountID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AccountsService) ListAccountTypes(ctx context.Context, offset, limit int) (*http.Response, error) {
	u := fmt.Sprintf("v1/account-types?offset=%d&limit=%d", offset, limit)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AccountsService) AccountTypeInfo(ctx context.Context, accountTypeID string) (*http.Response, error) {

	u := fmt.Sprintf("v1/account-types/%s", accountTypeID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
