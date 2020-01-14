package domo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

// StreamsService handles communication with the streams
// related methods of the Domo API.
//
// Domo API Docs: https://developer.domo.com/docs/streams-api-reference/streams
type StreamsService service

// List the streams. Limit should be between 1 and 500.
func (s *StreamsService) List(ctx context.Context, limit, offset int) ([]*StreamDataset, *http.Response, error) {
	if limit < 1 {
		return nil, nil, fmt.Errorf("limit must be above 0, but %d is not", limit)
	}
	if limit > 500 {
		return nil, nil, fmt.Errorf("limit must be 500 or below, but %d is not", limit)
	}
	u := fmt.Sprintf("v1/streams?limit=%d&offset=%d", limit, offset)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var streams []*StreamDataset
	resp, err := s.client.Do(ctx, req, &streams)
	if err != nil {
		return nil, resp, err
	}

	return streams, resp, nil
}

// Info for the stream for the given stream id.
func (s *StreamsService) Info(ctx context.Context, streamID int) (*StreamDataset, *http.Response, error) {
	u := fmt.Sprintf("v1/streams/%d", streamID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *StreamDataset
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// CreateStream creates a stream to use to create executions and upload data with streams to a dataset.
func (s *StreamsService) CreateStream(ctx context.Context, schema StreamDatasetSchema) (*StreamDataset, *http.Response, error) {
	u := "v1/streams"
	req, err := s.client.NewRequest("POST", u, schema)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *StreamDataset
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil

}

// DeleteStream deletes a domo stream with the given stream id. It does not delete the dataset associated with the stream.
func (s *StreamsService) DeleteStream(ctx context.Context, streamID int) (*http.Response, error) {
	u := fmt.Sprintf("v1/streams/%d", streamID)
	req, err := s.client.NewRequest("DELETE", u, nil)
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

// ModifyStreamUpdateMethod updates whether the stream's update strategy is appending new data or replacing the dataset.
func (s *StreamsService) ModifyStreamUpdateMethod(ctx context.Context, streamID int, isAppending bool) (*StreamDataset, *http.Response, error) {
	u := fmt.Sprintf("v1/streams/%d", streamID)
	var m string
	if isAppending {
		m = "APPEND"
	} else {
		m = "REPLACE"
	}
	ds := struct {
		UpdateMethod string `json:"updateMethod"`
	}{UpdateMethod: m}
	req, err := s.client.NewRequest("PATCH", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *StreamDataset
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// CreateExecution creates a new execution for a given stream to upload dataparts to.
func (s *StreamsService) CreateExecution(ctx context.Context, streamID int) (*StreamExecution, *http.Response, error) {
	u := fmt.Sprintf("v1/streams/%d/executions", streamID)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var sExecution *StreamExecution
	resp, err := s.client.Do(ctx, req, &sExecution)
	if err != nil {
		return nil, resp, err
	}

	return sExecution, resp, nil

}

// ListExecutions lists Domo stream executions for a given stream ID, limit, and offset.
func (s *StreamsService) ListExecutions(ctx context.Context, streamID, limit, offset int) ([]*StreamExecution, *http.Response, error) {
	if limit < 1 {
		return nil, nil, fmt.Errorf("limit must be above 0, but %d is not", limit)
	}
	if limit > 500 {
		return nil, nil, fmt.Errorf("limit must be 500 or below, but %d is not", limit)
	}
	u := fmt.Sprintf("v1/streams/%d/executions?limit=%d&offset=%d", streamID, limit, offset)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var streamExecutions []*StreamExecution
	resp, err := s.client.Do(ctx, req, &streamExecutions)
	if err != nil {
		return nil, resp, err
	}

	return streamExecutions, resp, nil
}

// CommitExecution finalizes a stream execution and inserts data parts into the dataset for the stream.
func (s *StreamsService) CommitExecution(ctx context.Context, streamID, executionID int) (*StreamExecution, *http.Response, error) {
	u := fmt.Sprintf("v1/streams/%d/executions/%d/commit", streamID, executionID)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var sExecution *StreamExecution
	resp, err := s.client.Do(ctx, req, &sExecution)
	if err != nil {
		return nil, resp, err
	}

	return sExecution, resp, nil
}

// AbortExecution aborts the execution and abandons any uploaded data parts for that execution.
func (s *StreamsService) AbortExecution(ctx context.Context, streamID, executionID int) (*http.Response, error) {
	u := fmt.Sprintf("v1/streams/%d/executions/%d/abort", streamID, executionID)
	req, err := s.client.NewRequest("PUT", u, nil)
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

// UploadDataPartStr uploads a csv given as a string to an active stream execution.
func (s *StreamsService) UploadDataPartStr(ctx context.Context, streamID, executionID, part int, csvData string) (*StreamFragment, *http.Response, error) {
	u := fmt.Sprintf("v1/streams/%d/executions/%d/part/%d", streamID, executionID, part)
	buf := new(bytes.Buffer)
	buf.WriteString(csvData)
	req, err := s.client.NewRequest("PUT", u, buf)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "text/csv")

	var sFragment *StreamFragment
	resp, err := s.client.Do(ctx, req, &sFragment)
	if err != nil {
		return nil, resp, err
	}
	return sFragment, resp, nil
}

// UploadDataPart uploads an array of structs serialized to csv to an active stream execution.
func (s *StreamsService) UploadDataPart(ctx context.Context, streamID, executionID, part int, data []interface{}) (*StreamFragment, *http.Response, error) {
	return nil, nil, fmt.Errorf("UploadDataPart not implemented")
}
