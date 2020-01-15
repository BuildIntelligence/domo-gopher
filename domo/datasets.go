package domo

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

// DatasetsService handles communication with the dataset
// related methods of the Domo API.
//
// Domo API Docs: https://developer.domo.com/docs/dataset-api-reference/dataset
type DatasetsService service

// List the datasets. Limit should be between 1 and 50.
func (s *DatasetsService) List(ctx context.Context, limit, offset int) ([]*DatasetDetails, *http.Response, error) {
	if limit < 1 {
		return nil, nil, fmt.Errorf("limit must be above 0, but %d is not", limit)
	}
	if limit > 50 {
		return nil, nil, fmt.Errorf("limit must be 50 or below, but %d is not", limit)
	}
	u := fmt.Sprintf("v1/datasets?limit=%d&offset=%d", limit, offset)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var datasets []*DatasetDetails
	resp, err := s.client.Do(ctx, req, &datasets)
	if err != nil {
		return nil, resp, err
	}

	return datasets, resp, nil
}

// Info for the dataset for the given dataset id.
func (s *DatasetsService) Info(ctx context.Context, id string) (*DatasetDetails, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *DatasetDetails
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// Create a new Domo Dataset.
func (s *DatasetsService) Create(ctx context.Context, ds DatasetDetails) (*DatasetDetails, *http.Response, error) {
	u := "v1/datasets"
	req, err := s.client.NewRequest("POST", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *DatasetDetails
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// UpdateSchema updates the Dataset Schema for the Dataset ID provided.
func (s *DatasetsService) UpdateSchema(ctx context.Context, id string, schema Schema) (*DatasetDetails, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
	ds := struct {
		Schema Schema `json:"schema"`
	}{Schema: schema}
	req, err := s.client.NewRequest("PUT", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *DatasetDetails
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// UpdateName updates the Dataset Name for the Dataset ID provided.
func (s *DatasetsService) UpdateName(ctx context.Context, id, name string) (*DatasetDetails, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
	ds := struct {
		Name string `json:"name"`
	}{Name: name}
	req, err := s.client.NewRequest("PUT", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *DatasetDetails
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// UpdateDescription updates the Dataset Description for the Dataset ID provided.
func (s *DatasetsService) UpdateDescription(ctx context.Context, id, description string) (*DatasetDetails, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
	ds := struct {
		Description string `json:"description"`
	}{Description: description}
	req, err := s.client.NewRequest("PUT", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *DatasetDetails
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// Delete a specified Domo Dataset by Dataset ID.
func (s *DatasetsService) Delete(ctx context.Context, id string) (*http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
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

// UploadDataStr Uploads a string CSV to the given dataset. If the dataset is set to append it will append the CSV otherwise it will replace.
func (s *DatasetsService) UploadDataStr(ctx context.Context, id string, dataCSV string) (*http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s/data", id)
	buf := new(bytes.Buffer)
	buf.WriteString(dataCSV)
	req, err := s.client.NewRequest("POST", u, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/csv")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// UploadData serializes an array of structs to CSV and then uploads them to the Domo Dataset.
func (s *DatasetsService) UploadData(ctx context.Context, id string, data []interface{}) (*http.Response, error) {
	return nil, fmt.Errorf("Not Implemented")
}

// DownloadDatasetCSV retrieves the datasets data as a string CSV.
func (s *DatasetsService) DownloadDatasetCSV(ctx context.Context, id string, includeHeader bool) (string, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s/data?includeHeader=%t", id, includeHeader)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Accept", "text/csv")

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return "", resp, err
	}

	io.Copy(buf, resp.Body)
	csv := buf.String()
	return csv, resp, nil
}

// QueryData takes a sql query and uses it to return a string csv of the query table results for the dataset.
func (s *DatasetsService) QueryData(ctx context.Context, id, sqlQuery string) (string, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/query/execute/%s", id)
	b := struct {
		SQL string `json:"sql"`
	}{SQL: sqlQuery}
	req, err := s.client.NewRequest("POST", u, b)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Accept", "text/csv")

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return "", resp, err
	}

	io.Copy(buf, resp.Body)
	csv := buf.String()
	return csv, resp, nil
}
