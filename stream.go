package domo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateMethodAppend and UpdateMethodReplace for StreamDataset UpdateMethod.
const (
	UpdateMethodAppend  = "APPEND"
	UpdateMethodReplace = "REPLACE"
)

// StreamDataset describes the Domo Dataset for a Domo Stream.
type StreamDataset struct {
	ID            int              `json:"id,omitempty"`
	Dataset       *DatasetDetails  `json:"dataset,omitempty"`
	UpdateMethod  string           `json:"updateMethod,omitempty"`
	CreatedAt     string           `json:"createdAt,omitempty"`
	ModifiedAt    string           `json:"modifiedAt,omitempty"`
	LastExecution *StreamExecution `json:"lastExecution,omitempty"`
}

// StreamExecution describes the Execution for a given Domo Stream.
type StreamExecution struct {
	ID           int    `json:"id,omitempty"`
	StartedAt    string `json:"startedAt,omitempty"`
	EndedAt      string `json:"endedAt,omitempty"`
	CurrentState string `json:"currentState,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	ModifiedAt   string `json:"modifiedAt,omitempty"`
}

// StreamDatasetSchema describes the schema for a StreamDataset.
type StreamDatasetSchema struct {
	DatasetSchema *DatasetSchema `json:"schema,omitempty"`
	UpdateMethod  string         `json:"updateMethod,omitempty"`
}

// CreateNewStream creates a stream to use to create executions and upload data with streams to a dataset.
func (c *Client) CreateNewStream(schema StreamDatasetSchema) (*StreamDataset, error) {
	//POST
	domoURL := fmt.Sprintf("%s/v1/streams", c.baseURL)
	buf := new(bytes.Buffer) // I think this could be buf, err := json.Marshal(dataset) instead
	err := json.NewEncoder(buf).Encode(schema)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", domoURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result StreamDataset
	err = c.execute(req, &result, 201)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStreamDetails returns the details for a given stream ID dataset.
func (c *Client) GetStreamDetails(streamID int) (*StreamDataset, error) {
	domoURL := fmt.Sprintf("%s/v1/streams/%d", c.baseURL, streamID)

	var d *StreamDataset

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// ListStreams lists domo streams sliced with given limit and offset.
func (c *Client) ListStreams(limit int, offset int) ([]*StreamDataset, error) {

	domoURL := fmt.Sprintf("%s/v1/streams?limit=%d&offset=%d", c.baseURL, limit, offset)

	var d []*StreamDataset

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// UpdateStreamMeta updates the stream with the given StreamDataset values. Values that aren't changing can be omitted.
func (c *Client) UpdateStreamMeta(streamID int, streamDataset StreamDataset) (*StreamDataset, error) {
	domoURL := fmt.Sprintf("%s/v1/streams/%d", c.baseURL, streamID)

	buf := new(bytes.Buffer) // I think this could be buf, err := json.Marshal(dataset) instead
	err := json.NewEncoder(buf).Encode(streamDataset)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", domoURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result StreamDataset
	err = c.execute(req, &result, 204)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteStream deletes a domo stream with the given stream id. It does not delete the dataset associated with the stream.
func (c *Client) DeleteStream(streamID int) error {

	//DELETE
	domoURL := fmt.Sprintf("%s/v1/streams/%d", c.baseURL, streamID)
	req, err := http.NewRequest("DELETE", domoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	if err != nil {
		return err
	}

	return nil
}

// CreateStreamExecution creates a new execution for a given stream to upload dataparts to.
func (c *Client) CreateStreamExecution(streamID int) (*StreamExecution, error) {
	domoURL := fmt.Sprintf("%s/v1/streams/%d/executions", c.baseURL, streamID)
	req, err := http.NewRequest("POST", domoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result StreamExecution
	err = c.execute(req, &result, 201)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListStreamExecutions lists Domo stream executions for a given stream ID, limit, and offset.
func (c *Client) ListStreamExecutions(streamID int, limit int, offset int) ([]*StreamExecution, error) {
	domoURL := fmt.Sprintf("%s/v1/streams/%d/executions?limit=%d&offset=%d", c.baseURL, streamID, limit, offset)

	var d []*StreamExecution

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// UploadDataPart uploads a csv given as a string to an active stream execution.
func (c *Client) UploadDataPart(streamID int, executionID int, partNumber int, csvData string) error {
	domoURL := fmt.Sprintf("%s/v1/streams/%d/executions/%d/part/%d", c.baseURL, streamID, executionID, partNumber)
	buf := new(bytes.Buffer)
	buf.WriteString(csvData)
	req, err := http.NewRequest("POST", domoURL, buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "text/csv")

	err = c.execute(req, nil, 204)
	if err != nil {
		return err
	}

	return nil
}

// CommitExecution finalizes a stream execution and inserts data parts into the dataset for the stream.
func (c *Client) CommitExecution(streamID int, executionID int) (*StreamExecution, error) {

	domoURL := fmt.Sprintf("%s/v1/streams/%d/executions/%d/commit", c.baseURL, streamID, executionID)
	req, err := http.NewRequest("PUT", domoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result StreamExecution
	err = c.execute(req, &result, 200)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// AbortStreamExecution aborts the execution and abandons any uploaded data parts for that execution.
func (c *Client) AbortStreamExecution(streamID int, executionID int) error {
	domoURL := fmt.Sprintf("%s/v1/streams/%d/executions/%d/abort", c.baseURL, streamID, executionID)
	req, err := http.NewRequest("PUT", domoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.execute(req, nil, 204)
	if err != nil {
		return err
	}

	return nil
}
