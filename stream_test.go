package domo

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_GetStreamDetails(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/streams/get_stream_details.json")
	defer server.Close()

	streamInfo, err := client.GetStreamDetails(42)
	if err != nil {
		t.Fatal(err)
	}
	if streamInfo == nil {
		t.Fatal("Got nil Stream Details")
	}
	if streamInfo.ID != 42 {
		t.Error("Got wrong stream")
	}

}

func Test_GetStreamDetailsBadID(t *testing.T) {
	client, server := testClientString(http.StatusNotFound, `{"error": { "status": 404, "message": "domo err msg"}}`)
	defer server.Close()

	streamInfo, err := client.GetStreamDetails(0)
	if streamInfo != nil {
		t.Fatal("Expected nil stream, got", streamInfo.ID)
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 404 {
		t.Errorf("Expected HTTP 404, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func TestClient_CreateNewStream(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		schema StreamDatasetSchema
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *StreamDataset
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			got, err := c.CreateNewStream(tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateNewStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.CreateNewStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetStreamDetails(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		streamID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *StreamDataset
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			got, err := c.GetStreamDetails(tt.args.streamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetStreamDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetStreamDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListStreams(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*StreamDataset
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			got, err := c.ListStreams(tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListStreams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListStreams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UpdateStreamMeta(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		streamID      int
		streamDataset StreamDataset
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *StreamDataset
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			got, err := c.UpdateStreamMeta(tt.args.streamID, tt.args.streamDataset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateStreamMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UpdateStreamMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteStream(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		streamID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			if err := c.DeleteStream(tt.args.streamID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteStream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateStreamExecution(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		streamID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *StreamExecution
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			got, err := c.CreateStreamExecution(tt.args.streamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateStreamExecution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.CreateStreamExecution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListStreamExecutions(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		streamID int
		limit    int
		offset   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*StreamExecution
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			got, err := c.ListStreamExecutions(tt.args.streamID, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListStreamExecutions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListStreamExecutions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UploadDataPart(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		streamID    int
		executionID int
		partNumber  int
		csvData     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			if err := c.UploadDataPart(tt.args.streamID, tt.args.executionID, tt.args.partNumber, tt.args.csvData); (err != nil) != tt.wantErr {
				t.Errorf("Client.UploadDataPart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CommitExecution(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		streamID    int
		executionID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *StreamExecution
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			got, err := c.CommitExecution(tt.args.streamID, tt.args.executionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CommitExecution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.CommitExecution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_AbortStreamExecution(t *testing.T) {
	type fields struct {
		http      *http.Client
		baseURL   string
		AutoRetry bool
	}
	type args struct {
		streamID    int
		executionID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:      tt.fields.http,
				baseURL:   tt.fields.baseURL,
				AutoRetry: tt.fields.AutoRetry,
			}
			if err := c.AbortStreamExecution(tt.args.streamID, tt.args.executionID); (err != nil) != tt.wantErr {
				t.Errorf("Client.AbortStreamExecution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
