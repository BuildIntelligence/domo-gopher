// Package domo provides utilities for easily interfacing
// with Domo's APIs.
package domo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Version is the version of this lib.
const Version = "1.0.0"

const (
	// DomoDateFormat can be used with time.Parse to create time.Time values
	// from domo date strings.
	DomoDateFormat = "2018-09-17"
	// DomoTimestampFormat can be used with time.Parse to create time.Time
	// values from domo timestamp strings. ISO 8601 UTC timestamp 0 offset
	DomoTimestampFormat = "2018-09-17T15:04:05Z"

	defaultRetryDuration = time.Second * 5

	// rateLimitExceededStatus Code is the HTTP code the server returns when request freq. is too high.
	rateLimitExceededStatusCode = 429
)

// Client is a client for working with the Domo API.
type Client struct {
	http    *http.Client
	baseURL string

	AutoRetry bool
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

// decode an Error from an io.Reader.
func (c *Client) decodeError(resp *http.Response) error {
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBody) == 0 {
		return fmt.Errorf("domo: HTTP %d: %s (body empty)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	buf := bytes.NewBuffer(responseBody)

	var e struct {
		E Error `json:"error"`
	}
	err = json.NewDecoder(buf).Decode(&e)
	if err != nil {
		return fmt.Errorf("domo: couldn't decode err: (%d) [%s]", len(responseBody), responseBody)
	}

	if e.E.Message == "" {
		e.E.Message = fmt.Sprintf("domo: unexpected HTTP %d: %s (empty error)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return e.E
}

func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

func isFailure(code int, validCodes []int) bool {
	for _, item := range validCodes {
		if item == code {
			return false
		}
	}
	return true
}

// execute a non-GET request
func (c *Client) execute(req *http.Request, result interface{}, needsStatus ...int) error {
	for {
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if c.AutoRetry && shouldRetry(resp.StatusCode) {
			time.Sleep(retryDuration(resp))
			continue
		}
		if resp.StatusCode != http.StatusOK && isFailure(resp.StatusCode, needsStatus) {
			return c.decodeError(resp)
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return err
			}
		}
		break
	}
	return nil
}

func retryDuration(resp *http.Response) time.Duration {
	raw := resp.Header.Get("Retry-After")
	if raw == "" {
		return defaultRetryDuration
	}
	seconds, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return defaultRetryDuration
	}
	return time.Duration(seconds) * time.Second
}

func (c *Client) get(url string, result interface{}) error {
	for {
		resp, err := c.http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode && c.AutoRetry {
			time.Sleep(retryDuration(resp))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return c.decodeError(resp)
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}

	return nil
}

func (c *Client) getCSV(url string) (string, error) {
	var s string
	for {
		resp, err := c.http.Get(url)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode && c.AutoRetry {
			time.Sleep(retryDuration(resp))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return "", c.decodeError(resp)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		s = string(b)

		break
	}

	return s, nil
}
