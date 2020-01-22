// Package domo provides utilities for easily interfacing
// with Domo's APIs.
package domo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Version is the version of this lib.
const Version = "2.0.1"

const (
	defaultBaseURL = "https://api.domo.com/"
	userAgent      = "domo-gopher"
	// DomoDateFormat can be used with time.Parse to create time.Time values
	// from domo date strings.
	DomoDateFormat = "2018-09-17"
	// DomoTimestampFormat can be used with time.Parse to create time.Time
	// values from domo timestamp strings. ISO 8601 UTC timestamp 0 offset
	DomoTimestampFormat = "2018-09-17T15:04:05Z"
)

// Client is a client for working with the Domo API.
type Client struct {
	http     *http.Client
	clientMu sync.Mutex   // clientMu protects the client during calls that modify the CheckRedirect func.
	client   *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to public Domo API. BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL
	// User agent used when communicating with the Domo API.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Domo API.
	Datasets *DatasetsService
	Streams  *StreamsService
	Users    *UsersService
	Groups   *GroupsService
	Pages    *PagesService
	Logs     *ActivityLogsService
}

type service struct {
	client *Client
}

// NewClient returns a new Domo API client. If a nil httpClient is
// provided, a new http.Client will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library, or the domo-gopher auth).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.Datasets = (*DatasetsService)(&c.common)
	c.Streams = (*StreamsService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	c.Groups = (*GroupsService)(&c.common)
	c.Pages = (*PagesService)(&c.common)
	c.Logs = (*ActivityLogsService)(&c.common)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as the
// request body. If the body is an io.ReadWriter it's used as the body instead of JSON encoding it.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		if w, ok := body.(io.ReadWriter); ok {
			buf = w
		} else {
			buf = new(bytes.Buffer)
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			err := enc.Encode(body)
			if err != nil {
				return nil, err
			}
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "applicatin/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and teh context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if e, ok := err.(*url.Error); ok {
			return nil, e
		}

		return nil, err
	}
	defer resp.Body.Close()
	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors cause by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

// CheckResponse checks teh API response for errors, and returns then if present.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	// errorResp := &Error{}
	var e struct {
		E Error `json:"error"`
	}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, &e)
	}
	return e.E
}

// Error represents an error returned by the Domo API.
type Error struct {
	// Short desc of the error.
	Message string `json:"message"`
	// Http Status Code
	Status int `json:"status"`
}

func (e Error) Error() string {
	return e.Message
}
