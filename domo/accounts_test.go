package domo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestAccountsService_List(t *testing.T) {
	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData,ScopeAccount)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	ctx := context.Background()

	resp, err := client.Accounts.List(ctx)
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}

func TestAccountsService_Info(t *testing.T) {
	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData,ScopeAccount)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	ctx := context.Background()

	resp, err := client.Accounts.Info(ctx, "71")
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}

func TestAccountsService_ListAccountTypes(t *testing.T) {
	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData,ScopeAccount)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	ctx := context.Background()

	resp, err := client.Accounts.ListAccountTypes(ctx, 0, 500)
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}

func TestAccountsService_AccountTypeInfo(t *testing.T) {
	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData,ScopeAccount)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	ctx := context.Background()

	resp, err := client.Accounts.AccountTypeInfo(ctx, "asana")
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
