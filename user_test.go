package domo

import (
	"net/http"
	"testing"
)

func TestClient_GetUserInfo(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/users/retrieve_user.json")
	defer server.Close()

	userInfo, err := client.GetUserInfo(871428330)
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

func TestClient_GetUserInfoBadID(t *testing.T) {
	client, server := testClientString(http.StatusNotFound, `{"error": { "status": 404, "message": "domo err msg"}}`)
	defer server.Close()

	userInfo, err := client.GetUserInfo(0)
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

func TestClient_CreateNewUser(t *testing.T) {

	filename := "test_data/users/create_user.json"
	client, server := testClientFile(http.StatusOK, filename)
	defer server.Close()

	user := User{Name: "Leonhard Euler"}
	res, err := client.CreateNewUser(user, false)
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

func TestClient_UpdateUser(t *testing.T) {

}

func Test_ListUsers(t *testing.T) {

}

func TestClient_DeleteUser(t *testing.T) {
	client, server := testClientString(http.StatusNoContent, "")
	defer server.Close()

	err := client.DeleteUser(1)
	if err != nil {
		t.Fatal(err)
	}
}
