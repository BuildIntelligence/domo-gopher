// Package domo provides utilities for easily interfacing
// with Domo's APIs.
package domo

import (
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
)

var (
	domod      = flag.Bool("domo", false, "Check for getting Domo'd, Run Domo integration tests to check if Domo's API even works how they have it documented")
	domogopher = flag.Bool("domoGopher", false, "Run Domo integration tests to check if Domo Gopher works correctly with Domo's API")
)

func testClientV2(code int, body io.Reader, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range validators {
			v(r)
		}
		w.WriteHeader(code)
		io.Copy(w, body)
		r.Body.Close()
		if closer, ok := body.(io.Closer); ok {
			closer.Close()
		}
	}))

	client := NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u

	return client, server
}

// Client whose reqs will always return a specified status code and body.
func testClientStringV2(code int, body string, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	return testClientV2(code, strings.NewReader(body))
}

// Client whose reqs will always return a specified status code and return a body read from a file.
func testClientFileV2(code int, filename string, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return testClientV2(code, f)
}
