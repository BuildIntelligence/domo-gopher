package domo_gopher

import (
	"github.com/bmizerany/assert"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var clientID string
var clientSecret string
var domo Domo

// Create out api variables for easy access
func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	clientID = os.Getenv("DOMO_CLIENT_ID")
	clientSecret = os.Getenv("DOMO_SECRET")
	domo = New(clientID, clientSecret)
	domo.Authorize()
	runTests := m.Run()
	os.Exit(runTests)
}

// Should create a new Domo object.
func TestNew(t *testing.T) {
	assert.T(t, domo.clientID == clientID, "clientID should be the same")
	assert.T(t, domo.clientSecret == clientSecret, "clientSecret should be the same")
}

// Should create a new Domo obj.
func TestDomo_Authorize(t *testing.T) {
	result, err := domo.Authorize()
	assert.T(t, result, "should be true")
	assert.T(t, len(err) == 0, "should be nil")
	assert.T(t, len(domo.accessToken) > 0, "should not be nil")
}

// Should create a new Domo obj.
func TestDomo_Request(t *testing.T) {
	result, err := domo.Request("GET", "datasets/%s", nil, "77faea51-68ab-4dd3-ae1a-8992bc1b58a8")
	assert.T(t, result != nil, "Shouldn't be null")
	assert.T(t, err == nil, "Should be null")
}

// Should create a new Domo obj.
func TestDomo_Get(t *testing.T) {
	result, err := domo.Get("datasets/%s", nil, "")
	assert.T(t, result != nil, "Shouldn't be null")
	assert.T(t, err == nil, "Should be null")
}

// Should create a new Domo obj.
func TestGetEncodedKeys(t *testing.T) {
	result := domo.getEncodedKeys()
	assert.T(t, len(result) > 0, "shouldn't be null")
}

// Should create a new Domo obj.
func TestUnauthorizedResponse(t *testing.T) {
	result := unauthorizedResonse([]byte(`"{error: {	status: 401, message: "Not Authorized"}}"`))
	assert.T(t, result, "should be true")
}

// Should create a new Domo obj.
func TestCreateURL(t *testing.T) {
	result := domo.createTargetURL("datasets")
	assert.T(t, result == "https://api.domo.com/v1/datasets", "should be same URL")
}
