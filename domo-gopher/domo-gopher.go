// Package domo_gopher:
// domo-gopher provides an easy to use
// API to interact with Domo's API
package domo_gopher

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
)

// Domo struct to wrap our req ops.
type Domo struct {
	clientID     string
	clientSecret string
	accessToken  string
}

const (
	BASE_URL           = "https://api.domo.com"
	AUTHENTICATION_URL = "https://api.domo.com/oauth/token" // TODO: scope
	API_VERSION        = "v1"
)

// Creates a new Domo API object with the
// clientID and ClientSecret
// Usage:
//		domo.New("abc","xyz")
func New(clientID, clientSecret string) Domo {
	return initialize(clientID, clientSecret)
}

func initialize(clientID, clientSecret string) Domo {
	domo := Domo{clientID: clientID, clientSecret: clientSecret}
	return domo
}

// Authorizes app against Domo
func (domo *Domo) Authorize() (bool, []error) {
	result := false

	// Get Encoded Access KEys for Auth
	auth := fmt.Sprintf("Basic %s", domo.getEncodedKeys())

	// create a new request to get our access token
	// and send our Keys on Authorization Header
	request := gorequest.New()
	request.Post(AUTHENTICATION_URL)
	request.Set("Authorization", auth)
	request.Send("grant_type=client_credentials&scope=data")

	_, body, errs := request.End()

	// Parse res to simplejson obj
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		fmt.Println("[Authorize] Error parsing Json!")
		errs = []error{err}
	}

	// check whether we got the access token or not
	jsToken, exists := js.CheckGet("access_token")
	if exists {
		// If we got it then assign it to the domo object.
		domo.accessToken, err = jsToken.String()
		if err != nil {
			fmt.Println("[Authorize] Error Getting Access Token from Json!")
		}
		result = true
	}

	return result, errs
}

// Creates a new GET req to Domo and returns
// the response as a map[string]interface{}.
//
// format: target enpoind format like "datasets/%s" - string
//
// data: content to be sent with req - map[string]interface{}
//
// args: Arguments to be used based on format
//
// Usage:
//		domo.Get("datasets/%s",nil,0absdfsdfesfiljk)
func (domo *Domo) Get(format string, data map[string]interface{}, args ...interface{}) ([]byte, []error) {
	return domo.Request("GET", format, data, args...)
}

// Creates a new POST Request to Domo and returns
// the res as a map[string]interface{}
//
// format: target enpoint format like "datasets/%s" - string
//
// data: content ot be sent with the req - map[string]interface{}
//
// args: Arguments to be used based on format
//
// Usage:
//		domo.Post("datasets/%s",map[string]interface{},"what")
func (domo *Domo) Post(format string, data map[string]interface{}, args ...interface{}) ([]byte, []error) {
	return domo.Request("POST", format, data, args...)
}

// Creates a new PUT Request to Domo and returns
// the res as a map[string]interface{}
//
// format: target enpoint format like "datasets/%s" - string
//
// data: content ot be sent with the req - map[string]interface{}
//
// args: Arguments to be used based on format
//
// Usage:
//		domo.Post("datasets/%s",map[string]interface{},"what")
func (domo *Domo) Put(format string, data map[string]interface{}, args ...interface{}) ([]byte, []error) {
	return domo.Request("PUT", format, data, args...)
}

// Creates a new DELETE Request to Domo and returns
// the res as a map[string]interface{}
//
// format: target enpoint format like "datasets/%s" - string
//
// data: content ot be sent with the req - map[string]interface{}
//
// args: Arguments to be used based on format
//
// Usage:
//		domo.Post("datasets/%s",map[string]interface{},"what")
func (domo *Domo) Delete(format string, data map[string]interface{}, args ...interface{}) ([]byte, []error) {
	return domo.Request("DELETE", format, data, args...)
}

// Creates a new Request to Domo and returns
// the res as a map[string]interface{}.
//
// method: GET/POST/PUT/DELETE - string
//
// format: target enpoint format - string
//
// data: contetnt to be sent with the req - map[string]interface{}
//
// args: Arguments to be used based on format
//
// Usage:
//		domo.request("GET","datasets/%s",nil,"0asdfasdf")
func (domo *Domo) Request(method, format string, data map[string]interface{}, args ...interface{}) ([]byte, []error) {

	// create endpoint based on passed format
	endpoint := fmt.Sprintf(format, args...)

	targetURL := domo.createTargetURL(endpoint)

	request := gorequest.New()

	// Check method type to call corresponding
	// go-request method
	if method == "GET" {
		request.Get(targetURL)
	}
	if method == "POST" {
		request.Post(targetURL)
	}
	if method == "PUT" {
		request.Put(targetURL)
	}
	if method == "DELETE" {
		request.Delete(targetURL)
	}

	request.Set("Authorization", fmt.Sprintf("Bearer %s", domo.accessToken))

	// Add the data to the request if it
	// isn't null
	if data != nil {
		jsonData, _ := getJsonBytesFromMap(data)
		if jsonData != nil {
			request.Send(string(jsonData))
		}
	}

	_, body, errs := request.End()

	result := []byte(body)
	if unauthorizedResonse(result) {
		result = nil
		errs = []error{
			errors.New("Authorization Error. Make sure you called Domo.Authorize() method!"),
			errors.New(body)}
	}

	return result, errs
}

// Checks for the res content to see if we
// received a not authorized err.
func unauthorizedResonse(body []byte) bool {

	// Parse response to simplejson obj
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println("[unauthorizedResponse] Error parsing Json!")
		return true
	}

	// check whether we got an error or not.
	_, exists := js.CheckGet("error")
	if exists {
		return true
	}

	return false
}

// Creates target URL for making a Domo Request
// to a given endpoint
func (domo *Domo) createTargetURL(endpoint string) string {
	result := fmt.Sprintf("%s/%s/%s", BASE_URL, API_VERSION, endpoint)
	return result
}

// returns base64 encoded authorization
// keys for Domo.
func (domo *Domo) getEncodedKeys() string {
	data := fmt.Sprintf("%v:%v", domo.clientID, domo.clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	return encoded
}

// Extracts Json Bytes from map[string]interface
func getJsonBytesFromMap(data map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Invalid data object, can't parse to json:")
		fmt.Println("Error:", err)
		fmt.Println("Data:", data)
		return nil, err
	}
	return jsonData, nil
}
