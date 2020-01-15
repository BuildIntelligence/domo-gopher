package domo

import (
	"context"
	"net/http"
	"testing"
)

func TestUsersService_Info(t *testing.T) {
	client, server := testClientFileV2(http.StatusOK, "../test_data/users/retrieve_user.json")
	defer server.Close()
	ctx := context.Background()

	userInfo, _, err := client.Users.Info(ctx, 871428330)
	if err != nil {
		t.Fatal(err)
	}
	if userInfo == nil {
		t.Fatal("Got nil User Details")
	}
	if userInfo.ID != 871428330 {
		t.Error("Got wrong User")
	}
	if userInfo.Name != "Leonard Euler" {
		t.Error("Got wrong Name")
	}
}

func TestUsersService_InfoBadID(t *testing.T) {
	client, server := testClientStringV2(http.StatusNotFound, `{"error": { "status": 404, "message": "domo err msg"}}`)
	defer server.Close()
	ctx := context.Background()

	userInfo, _, err := client.Users.Info(ctx, 0)
	if userInfo != nil {
		t.Fatal("Expected nil user, got", userInfo.ID)
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 404 {
		t.Errorf("Expected HTTP 404, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func TestUsersService_Create(t *testing.T) {

	filename := "../test_data/users/create_user.json"
	client, server := testClientFileV2(http.StatusOK, filename)
	defer server.Close()
	ctx := context.Background()

	user := User{Name: "Leonhard Euler"}
	res, _, err := client.Users.Create(ctx, user, false)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Got nil user")
	}
	if res.Name != "Leonhard Euler" {
		t.Error("Got Wrong Name")
	}

}


func TestUsersService_Delete(t *testing.T) {
	client, server := testClientStringV2(http.StatusNoContent, "")
	defer server.Close()
	ctx := context.Background()
	_, err := client.Users.Delete(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
}
