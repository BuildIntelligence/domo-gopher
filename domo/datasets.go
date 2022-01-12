package domo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Data types for Domo Columns
const (
	ColumnTypeString   = "STRING"
	ColumnTypeLong     = "LONG"
	ColumnTypeDate     = "DATE"
	ColumnTypeDatetime = "DATETIME"
	ColumnTypeDouble   = "DOUBLE"
	ColumnTypeDecimal  = "DECIMAL"
)

// Dataset contains basic data about a domo Dataset.
type Dataset struct {
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
func (s *DatasetsService) List(ctx context.Context, limit, offset int) ([]*Dataset, *http.Response, error) {
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

	var datasets []*Dataset
	resp, err := s.client.Do(ctx, req, &datasets)
	if err != nil {
		return nil, resp, err
	}

	return datasets, resp, nil
}

// Info for the dataset for the given dataset id.
func (s *DatasetsService) Info(ctx context.Context, id string) (*Dataset, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *Dataset
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// Create a new Domo Dataset.
func (s *DatasetsService) Create(ctx context.Context, ds Dataset) (*Dataset, *http.Response, error) {
	u := "v1/datasets"
	req, err := s.client.NewRequest("POST", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *Dataset
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// UpdateSchema updates the Dataset Schema for the Dataset ID provided.
func (s *DatasetsService) UpdateSchema(ctx context.Context, id string, schema Schema) (*Dataset, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
	ds := struct {
		Schema Schema `json:"schema"`
	}{Schema: schema}
	req, err := s.client.NewRequest("PUT", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *Dataset
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// UpdateName updates the Dataset Name for the Dataset ID provided.
func (s *DatasetsService) UpdateName(ctx context.Context, id, name string) (*Dataset, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
	ds := struct {
		Name string `json:"name"`
	}{Name: name}
	req, err := s.client.NewRequest("PUT", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *Dataset
	resp, err := s.client.Do(ctx, req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// UpdateDescription updates the Dataset Description for the Dataset ID provided.
func (s *DatasetsService) UpdateDescription(ctx context.Context, id, description string) (*Dataset, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/%s", id)
	ds := struct {
		Description string `json:"description"`
	}{Description: description}
	req, err := s.client.NewRequest("PUT", u, ds)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var d *Dataset
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
func (s *DatasetsService) UploadData(ctx context.Context, id string, data interface{}) (*http.Response, error) {
	// k := reflect.TypeOf(data).Kind()
	// if k != reflect.Slice && k != reflect.Array {
	// 	return nil, fmt.Errorf("Expected data to be a Slice or Array but got type %s", k)
	// }
	// schema := GenerateDataSetSchema(reflect.TypeOf(data).Elem())
	// _, _, err := s.UpdateSchema(ctx, id, schema)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, fmt.Errorf("Not Implemented")
}

// HasSchemaChanged checks if a structs generated Schema differs from the Schema of a domo dataset.
// use FindSchemaChanges to retrieve a list of schema differences.
func (s *DatasetsService) HasSchemaChanged(ctx context.Context, datasetID string, rType reflect.Type) (bool, error) {
	ds, _, err := s.Info(ctx, datasetID)
	if err != nil {
		return true, err
	}
	currentSchema := GenerateDataSetSchema(rType)
	if len(currentSchema.Columns) != len(ds.Schema.Columns) {
		return true, nil
	}
	idxDiff := checkForSchemaChangeByColumnIndexComparision(currentSchema, ds.Schema)
	if idxDiff.DiffsCount() > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func checkForSchemaChangeByColumnNameMatching(local, domo Schema) SchemaDiffError {
	var diffs SchemaDiffError
	type schemaCol struct {
		ColumnType string
		Name       string
		Index      int
	}
	localSchema := make(map[string]schemaCol, len(local.Columns))
	for i, c := range local.Columns {
		localSchema[c.Name] = schemaCol{
			ColumnType: c.ColumnType,
			Name:       c.Name,
			Index:      i,
		}
	}
	domoSchema := make(map[string]schemaCol, len(domo.Columns))
	for i, c := range domo.Columns {
		domoSchema[c.Name] = schemaCol{
			ColumnType: c.ColumnType,
			Name:       c.Name,
			Index:      i,
		}
	}
	// Check domo schema for local matches and for columns to delete.
	for _, col := range domoSchema {
		c, ok := localSchema[col.Name]
		if ok {
			// check if column type is different between Domo and local.
			if c.ColumnType != col.ColumnType {
				diffs.ColumnTypeMismatch = append(diffs.ColumnTypeMismatch, SchemaMismatch{
					DomoColumnIndex:    col.Index,
					DomoColumnName:     col.Name,
					ComparedColumnName: c.Name,
					DomoColumnType:     col.ColumnType,
					ComparedColumnType: c.ColumnType,
					Message:            fmt.Sprintf("Column Type Change: Expected column %s (%d) to be type %s (current Domo Column Type), found local type %s", col.Name, col.Index, col.ColumnType, c.ColumnType),
				})
			}
		} else {
			// No name match in local schema, assuming it's a column to delete.
			diffs.ColumnsToDeleteFromDomo = append(diffs.ColumnsToDeleteFromDomo, SchemaMismatch{
				DomoColumnIndex:    col.Index,
				DomoColumnName:     col.Name,
				ComparedColumnName: "",
				DomoColumnType:     col.ColumnType,
				ComparedColumnType: "",
				Message:            fmt.Sprintf("Missing a column %d found in Domo Schema: Expected a column %s (%s). Either it's missing from the local schema or it's a column to delete from Domo Schema", col.Index, col.Name, col.ColumnType),
			})
		}
	}
	// Check to see if there's any columns added
	for _, col := range localSchema {
		_, ok := domoSchema[col.Name]
		if !ok {
			// No name match in domo, assuming it's a new column.
			diffs.ColumnsToAddToDomo = append(diffs.ColumnsToAddToDomo, SchemaMismatch{
				DomoColumnIndex:    0,
				DomoColumnName:     "",
				ComparedColumnName: col.Name,
				DomoColumnType:     "",
				ComparedColumnType: col.ColumnType,
				Message:            fmt.Sprintf("Extra Column Found: Found a column %s (%s) that's not in Domo. Either it's an extra column that should be removed from the local schema or it's a column to add to Domo Schema", col.Name, col.ColumnType),
			})
		}
	}
	return diffs
}
func checkForSchemaChangeByColumnIndexComparision(local, domo Schema) SchemaDiffError {
	var diffs SchemaDiffError

	for i, col := range local.Columns {
		if col.Name != domo.Columns[i].Name {
			diffs.NameMismatch = append(diffs.NameMismatch, SchemaMismatch{i, domo.Columns[i].Name, col.Name, domo.Columns[i].ColumnType, col.ColumnType, fmt.Sprintf("Column Name Change: Expected column %d to be named %s, currently it's called %s", i, col.Name, domo.Columns[i].Name)})
		}
		if col.ColumnType != domo.Columns[i].ColumnType {
			diffs.ColumnTypeMismatch = append(diffs.ColumnTypeMismatch,
				SchemaMismatch{
					DomoColumnIndex:    i,
					DomoColumnName:     domo.Columns[i].Name,
					ComparedColumnName: col.Name,
					DomoColumnType:     domo.Columns[i].ColumnType,
					ComparedColumnType: col.ColumnType,
					Message:            fmt.Sprintf("Column Name Change: Expected column (%d) %s to be type %s, in Domo it's type %s", i, col.Name, col.ColumnType, domo.Columns[i].ColumnType),
				})
		}
	}
	return diffs
}

// FindSchemaChanges creates a list of schema differences between the schema generated by a struct and a domo dataset
// schema. If the Generated schema and Domo schema have the same column count, it'll do a by index comparision and note
// column name changes and column data type mismatches. If the Generated schema and Domo schema have a different column
// count it'll try to match by column name for column data type mismatches and for columns missing from Domo/missing from
// generated schema. If the struct has had a column name change AND a column added/removed, it'll count the name change
// as a column to Add to domo upon a schema update and it'll show the original name as a missing column that's still in
// the Domo Schema.
func (s *DatasetsService) FindSchemaChanges(ctx context.Context, id string, rType reflect.Type) error {
	ds, _, err := s.Info(ctx, id)
	if err != nil {
		return err
	}
	currentSchema := GenerateDataSetSchema(rType)
	var diffs SchemaDiffError
	if len(currentSchema.Columns) != len(ds.Schema.Columns) {
		// Column count is different. Check for differences by column name.
		byNameDiffs := checkForSchemaChangeByColumnNameMatching(currentSchema, ds.Schema)
		if byNameDiffs.DiffsCount() > 0 {
			diffs.Merge(byNameDiffs)
		}
	} else {
		// Column count is the same. Assuming No Column Deletes/Adds and comparing column index to check for name and
		// type diffs.
		idxDiffs := checkForSchemaChangeByColumnIndexComparision(currentSchema, ds.Schema)
		if idxDiffs.DiffsCount() > 0 {
			diffs.Merge(idxDiffs)
		}
	}

	diffsCount := len(diffs.NameMismatch) + len(diffs.ColumnTypeMismatch) + len(diffs.ColumnsToAddToDomo) + len(diffs.ColumnsToDeleteFromDomo)
	if diffsCount == 0 {
		return nil
	} else {
		return diffs
	}
}

type SchemaMismatch struct {
	DomoColumnIndex    int
	DomoColumnName     string
	ComparedColumnName string
	DomoColumnType     string
	ComparedColumnType string
	Message            string
}
type SchemaDiffError struct {
	NameMismatch            []SchemaMismatch
	ColumnTypeMismatch      []SchemaMismatch
	ColumnsToDeleteFromDomo []SchemaMismatch
	ColumnsToAddToDomo      []SchemaMismatch
}

func (s *SchemaDiffError) OnlyColumnNameChanges() bool {
	if len(s.ColumnTypeMismatch) > 0 {
		return false
	}
	if len(s.ColumnsToAddToDomo) > 0 {
		return false
	}
	if len(s.ColumnsToDeleteFromDomo) > 0 {
		return false
	}
	return true
}
func (s *SchemaDiffError) DiffsCount() int {
	return len(s.NameMismatch) + len(s.ColumnTypeMismatch) + len(s.ColumnsToDeleteFromDomo) + len(s.ColumnsToAddToDomo)
}

func (s *SchemaDiffError) Merge(otherSchemaDiff SchemaDiffError) {
	if len(otherSchemaDiff.NameMismatch) > 0 {
		s.NameMismatch = append(s.NameMismatch, otherSchemaDiff.NameMismatch...)
	}
	if len(otherSchemaDiff.ColumnTypeMismatch) > 0 {
		s.ColumnTypeMismatch = append(s.ColumnTypeMismatch, otherSchemaDiff.ColumnTypeMismatch...)
	}
	if len(otherSchemaDiff.ColumnsToDeleteFromDomo) > 0 {
		s.ColumnsToDeleteFromDomo = append(s.ColumnsToDeleteFromDomo, otherSchemaDiff.ColumnsToDeleteFromDomo...)
	}
	if len(otherSchemaDiff.ColumnsToAddToDomo) > 0 {
		s.ColumnsToAddToDomo = append(s.ColumnsToAddToDomo, otherSchemaDiff.ColumnsToAddToDomo...)
	}
}

func (e SchemaDiffError) Error() string {
	diffsCount := len(e.NameMismatch) + len(e.ColumnTypeMismatch) + len(e.ColumnsToAddToDomo) + len(e.ColumnsToDeleteFromDomo)
	//diffs := []string{fmt.Sprintf("%d schema differences found:", diffsCount)}
	firstLine := fmt.Sprintf("%d schema differences found:", diffsCount)

	// Might be overkill doing it this way, but this way the array doesn't get reallocated a handful of times to
	// increase capacity. Given that the biggest column count a real-world scenario will ever get is probably 250ish,
	// and most real-world scenarios the datasets will have less than 80, it realistically will only save a dozen or so
	// allocations in a scenario where every single column has a name and/or type change vs using append.
	diffs := make([]string, diffsCount, diffsCount+5)
	diffs[0] = firstLine
	idx := 1
	for _, nameMismatch := range e.NameMismatch {
		//diffs = append(diffs, fmt.Sprintf("Column Name Change: Expected column %d to be named %s, currently it's called %s", nameMismatch.DomoColumnIndex, nameMismatch.ComparedColumnName, nameMismatch.DomoColumnName))
		diffs[idx] = fmt.Sprintf("Column Name Change: Expected column %d to be named %s, currently it's called %s", nameMismatch.DomoColumnIndex, nameMismatch.ComparedColumnName, nameMismatch.DomoColumnName)
		idx++
	}
	for _, colTypeDiff := range e.ColumnTypeMismatch {
		//diffs = append(diffs, fmt.Sprintf("Column Type Change: Expected column (%d) %s to be type %s, in Domo it's type %s", colTypeDiff.DomoColumnIndex, colTypeDiff.DomoColumnName, colTypeDiff.ComparedColumnType, colTypeDiff.DomoColumnType))
		diffs[idx] = fmt.Sprintf("Column Type Change: Expected column (%d) %s to be type %s, in Domo it's type %s", colTypeDiff.DomoColumnIndex, colTypeDiff.DomoColumnName, colTypeDiff.ComparedColumnType, colTypeDiff.DomoColumnType)
		idx++
	}
	for _, deleteCol := range e.ColumnsToDeleteFromDomo {
		diffs[idx] = fmt.Sprintf("Missing a column found in Domo Schema: Expected a column %s (%s). Either it's missing from the local schema or it's a column to delete from Domo Schema", deleteCol.DomoColumnName, deleteCol.DomoColumnType)
		idx++
	}
	for _, addCol := range e.ColumnsToAddToDomo {
		diffs[idx] = fmt.Sprintf("Extra Column Found: Found a column %s (%s) that's not in Domo. Either it's an extra column that should be removed from the local schema or it's a column to add to Domo Schema", addCol.ComparedColumnName, addCol.ComparedColumnType)
		idx++
	}
	return strings.Join(diffs, "\n")
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

// QueryData takes a sql query and uses it to return a json string of the query table results for the dataset.
// see https://developer.domo.com/docs/dataset-api-reference/dataset#Query%20a%20DataSet for an example response.
func (s *DatasetsService) QueryData(ctx context.Context, id, sqlQuery string) (string, *http.Response, error) {
	u := fmt.Sprintf("v1/datasets/query/execute/%s", id)
	b := struct {
		SQL string `json:"sql"`
	}{SQL: sqlQuery}
	req, err := s.client.NewRequest("POST", u, b)
	if err != nil {
		return "", nil, err
	}

	// This API endpoint doesn't return a CSV, but it returns a JSON object with the query results and some metadata.
	// see https://developer.domo.com/docs/dataset-api-reference/dataset#Query%20a%20DataSet for an example response.
	req.Header.Set("Accept", "application/json")

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return "", resp, err
	}

	io.Copy(buf, resp.Body)
	csv := buf.String()
	return csv, resp, nil
}
