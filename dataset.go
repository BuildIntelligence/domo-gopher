package domo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Data types for Domo Columns
const (
	ColumnTypeString   = "STRING"
	ColumnTypeLong     = "LONG"
	ColumnTypeDate     = "DATE"
	ColumnTypeDatetime = "DATETIME"
	ColumnTypeDouble   = "DOUBLE"
)

// DatasetDetails contains basic data about a domo Dataset.
type DatasetDetails struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Columns       int      `json:"columns,omitempty"`
	Rows          int      `json:"rows,omitempty"`
	Schema        Schema   `json:"schema,omitempty"`
	CreatedAt     string   `json:"createdAt,omitempty"`
	UpdatedAt     string   `json:"updatedAt,omitempty"`
	DataCurrentAt string   `json:"dataCurrentAt,omitempty"`
	PDPEnabled    bool     `json:"pdpEnabled,omitempty"`
	Owner         *Owner   `json:"owner,omitempty"`
	Policies      []Policy `json:"policies,omitempty"`
}

// DatasetSchema contains basic data about dataset and schema
type DatasetSchema struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Rows        int    `json:"rows,omitempty"`
	Schema      Schema `json:"schema,omitempty"`
}

// Schema contains the columns describing a domo dataset schema.
type Schema struct {
	Columns []Column `json:"columns,omitempty"`
}

// Column describes a Domo Dataset Column.
type Column struct {
	ColumnType string `json:"type,omitempty"`
	Name       string `json:"name,omitempty"`
}

// Owner identifies a Domo User that owns a resource.
type Owner struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Policy describes PDP policies for a dataset.
type Policy struct {
	ID         int      `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	PolicyType string   `json:"type,omitempty"`
	UserIDs    []int    `json:"users,omitempty"`
	GroupIDs   []int    `json:"groups,omitempty"`
	Filters    []Filter `json:"filters,omitempty"`
}

// Filter describes a PDP policy Data filter for a dataset.
type Filter struct {
	Column   string   `json:"column,omitempty"`
	Not      bool     `json:"not,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Values   []string `json:"values,omitempty"`
}

// GetDatasets gets Domo Datasets Lists given limit and offset.
func (c *Client) GetDatasets(limit int, offset int) ([]*DatasetDetails, error) {
	domoURL := fmt.Sprintf("%s/v1/datasets?limit=%d&offset=%d", c.baseURL, limit, offset)

	var d []*DatasetDetails

	// Note that this is using lowercase `get` which is a wrapper aroung `Get` defined in domo.go
	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// GetDatasetInfo gets Domo Dataset Details for the given dataset id.
func (c *Client) GetDatasetInfo(id string) (*DatasetDetails, error) {
	domoURL := fmt.Sprintf("%s/v1/datasets/%s", c.baseURL, id)

	var d *DatasetDetails

	err := c.get(domoURL, &d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CreateDataset creates a new domo Dataset with the given DatasetDetails for Name, Description, Schema, etc.
func (c *Client) CreateDataset(dataset DatasetDetails) (*DatasetDetails, error) {
	domoURL := fmt.Sprintf("%s/v1/datasets", c.baseURL)

	buf := new(bytes.Buffer) // I think this could be buf, err := json.Marshal(dataset) instead
	err := json.NewEncoder(buf).Encode(dataset)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", domoURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result DatasetDetails
	err = c.execute(req, &result, 201)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateDatasetMeta updates a dataset with the given fields passed in dataset param. You can omit values that don't change.
func (c *Client) UpdateDatasetMeta(id string, dataset DatasetDetails) (*DatasetDetails, error) {
	domoURL := fmt.Sprintf("%s/v1/datasets/%s", c.baseURL, id)

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(dataset)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", domoURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var result DatasetDetails
	err = c.execute(req, &result, 200)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteDataset deletes the Domo Dataset for which the id was provided.
func (c *Client) DeleteDataset(id string) error {
	domoURL := fmt.Sprintf("%s/v1/datasets/%s", c.baseURL, id)
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

// ReplaceData replaces the Domo Dataset with the provided CSV String.
func (c *Client) ReplaceData(id string, dataRows string) error {
	domoURL := fmt.Sprintf("%s/v1/datasets/%s/data", c.baseURL, id)
	buf := new(bytes.Buffer)
	buf.WriteString(dataRows)
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

// ExportData returns the data from a given dataset id as a csv with or without header row.
func (c *Client) ExportData(id string, includeHeader bool) (string, error) {
	domoURL := fmt.Sprintf("%s/v1/datasets/%s/data?includeHeader=%t", c.baseURL, id, includeHeader)

	s, err := c.getCSV(domoURL)
	if err != nil {
		return "", err
	}

	return s, nil
}
