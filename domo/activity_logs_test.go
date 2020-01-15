package domo

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"testing"
)

func Test_generateAuditQueryUrlParams(t *testing.T) {
	params := AuditQueryParams{User: "123", Start: 1, End: 2, Limit: 5, Offset: 1}
	expected := "end=2&limit=5&offset=1&start=1&user=123" // It'll order params alphabetically
	actual := generateAuditQueryUrlParams(params)
	if actual != expected {
		t.Errorf("Expected URL params: %s\nFound URL params:  %s", expected, actual)
	}
}

func Test_generateAuditQueryUrlParams_no_user_set(t *testing.T) {
	params := AuditQueryParams{Start: 1, End: 2, Limit: 5, Offset: 1}
	expected := "end=2&limit=5&offset=1&start=1" // It'll order params alphabetically
	actual := generateAuditQueryUrlParams(params)
	if actual != expected {
		t.Errorf("Expected URL params: %s\nFound URL params:  %s", expected, actual)
	}
}

func Test_generateAuditQueryUrlParams_no_time_range_set(t *testing.T) {
	params := AuditQueryParams{User: "123", Limit: 5, Offset: 1}
	expected := "limit=5&offset=1&user=123" // It'll order params alphabetically
	actual := generateAuditQueryUrlParams(params)
	if actual != expected {
		t.Errorf("Expected URL params: %s\nFound URL params:  %s", expected, actual)
	}
}

func Test_generateAuditQueryUrlParams_no_limit_set(t *testing.T) {
	params := AuditQueryParams{User: "123", Start: 1, End: 2, Offset: 1}
	expected := "end=2&limit=50&offset=1&start=1&user=123" // It'll order params alphabetically, default for limit is 50
	actual := generateAuditQueryUrlParams(params)
	if actual != expected {
		t.Errorf("Expected URL params: %s\nFound URL params:  %s", expected, actual)
	}
}
func Test_generateAuditQueryUrlParams_no_parameters_set(t *testing.T) {
	params := AuditQueryParams{}
	expected := "limit=50&offset=0" // default for limit is 50, default for offset is 0. It'll order params alphabetically
	actual := generateAuditQueryUrlParams(params)
	if actual != expected {
		t.Errorf("Expected URL params: %s\nFound URL params:  %s", expected, actual)
	}
}

func Test_Entries(t *testing.T) {
	client, server := testClientFileV2(http.StatusOK, "../test_data/activity_log.json")
	defer server.Close()
	ctx := context.Background()

	params := AuditQueryParams{User: "1619916076", Start: 1513230600000, End: 1513231200000}
	logEntries, _, err := client.Logs.Entries(ctx, params)
	if err != nil {
		t.Fatal(err)
	}
	if logEntries == nil {
		t.Fatal("Got nil LogEntries")
	}
	if len(logEntries) != 1 {
		t.Errorf("Unexpected number of entries, Expected 1, Got %d", len(logEntries))
	}
}

func Test_Entries_Bad_User_ID(t *testing.T) {
	// Looks like insufficient scope is HTTP 403. created insufficient_scope_audit.json with body.
	// TODO: Get real Domo Error message response and update this... Looks like the API for this one is
	// also broken at the moment.
	// domo: couldn't decode err: (49) [An error occurred while retrieving audit messages]
	// HTTP 500
	client, server := testClientStringV2(http.StatusBadRequest, `{"error": { "status": 400, "message": "domo err msg"}}`)
	defer server.Close()
	ctx := context.Background()

	params := AuditQueryParams{User: "-123"} // Bad User ID
	logEntries, _, err := client.Logs.Entries(ctx, params)
	if logEntries != nil {
		t.Fatal("Expected nil log entries")
	}
	se, ok := err.(Error)
	if !ok {
		t.Errorf("Expected domo error, got %v", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Errorf("Unexpected error message: %s\nExpected: %s", se.Message, "domo err msg")
	}
}

func Test_Entries_no_optional_params_Integration_Test(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domod {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as this file")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeAudit)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		params := AuditQueryParams{} // No optional params
		domoURL := fmt.Sprintf("%s/v1/audit?%s", client.baseURL, generateAuditQueryParamsString(params))
		res, err := client.getRespBody(domoURL)
		if err != nil {
			t.Error(err)
		}
		if res == "" {
			t.Error("Expected response got empty string")
		}
		if res == "[]" {
			t.Error("Expected response with entries, got empty array")
		}
	}
}
func Test_LogEntries_Bad_User_ID_Integration_Test(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	// *domod = true // uncomment if you want to debug without passing -domo flag
	if *domod {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as this file")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeAudit)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		params := AuditQueryParams{User: "-123", Start: 1542211215000, End: 1542241576000} // Bad User ID
		domoURL := fmt.Sprintf("%s/v1/audit?%s", client.baseURL, generateAuditQueryParamsString(params))
		res, err := client.getRespBody(domoURL)
		// Bad User ID, so I would expect an Error, either StatusNotFound or StatusBadRequest
		if err == nil {
			t.Error(err)
		}
		if res != "" {
			t.Errorf("Expected no response got: %s", res)
		}
	}
}

func Test_LogEntries_time_range_params_only_Integration_Test(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domod {

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as this file")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeAudit)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		params := AuditQueryParams{Start: 1542211215000, End: 1542241576000} // Bad User ID
		domoURL := fmt.Sprintf("%s/v1/audit?%s", client.baseURL, generateAuditQueryParamsString(params))
		res, err := client.getRespBody(domoURL)
		if err != nil {
			t.Error(err)
		}
		if res == "" {
			t.Error("Expected response got empty string")
		}
		if res == "[]" {
			t.Error("Expected entries, got empty array")
		}
	}
}

func Test_Entries_User_ID_Param_Integration_Test(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domod {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as this file")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeAudit)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		TestUserID := "1704739518"

		params := AuditQueryParams{User: TestUserID, Start: 1542211215000, End: 1542241576000} // Ryan's User ID
		domoURL := fmt.Sprintf("%s/v1/audit?%s", client.baseURL, generateAuditQueryParamsString(params))
		res, err := client.getRespBody(domoURL)
		if err != nil {
			t.Error(err)
		}
		if res == "" {
			t.Errorf("Expected no response got: %s", res)
		}

		userActivityLog, err := client.GetActivityLogEntries(params)
		if err != nil {
			t.Error(err)
		}
		if userActivityLog == nil {
			t.Error("Expected entries, got nil")
		}
		if len(userActivityLog) == 0 {
			t.Error("Expected entries, Got none")
		}
		for _, entry := range userActivityLog {
			if entry.UserID != TestUserID {
				t.Errorf("Expected only entries with UserID %s, found one with UserID %s", TestUserID, entry.UserID)
			}
		}
	}
}
