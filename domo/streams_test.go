package domo

import (
	"context"
	"net/http"
	"testing"
)

func Test_GetStreamDetails(t *testing.T) {
	client, server := testClientFileV2(http.StatusOK, "../test_data/streams/get_stream_details.json")
	ctx := context.Background()
	defer server.Close()

	streamInfo, _, err := client.Streams.Info(ctx, 42)
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
	client, server := testClientStringV2(http.StatusNotFound, `{"error": { "status": 404, "message": "domo err msg"}}`)
	ctx := context.Background()
	defer server.Close()

	streamInfo, _, err := client.Streams.Info(ctx, 0)
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

func Test_ListStreams(t *testing.T) {
	type fields struct {
		code     int
		filename string
	}
	type args struct {
		limit  int
		offset int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		// want    []StreamDataset
		wantErr bool
	}{
		{name: "Test List Streams", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_streams.json"}, args: args{limit: 1, offset: 0}, wantErr: false},
		{name: "Test List Streams", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_streams.json"}, args: args{limit: 3, offset: 0}, wantErr: false},
		{name: "Test List Streams", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_streams.json"}, args: args{limit: 3, offset: 1}, wantErr: false},
		{name: "Test List Streams max limit", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_streams.json"}, args: args{limit: 50, offset: 1}, wantErr: false},
		{name: "Test List Streams over max limit", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_streams.json"}, args: args{limit: 950, offset: 1}, wantErr: true},
		{name: "Test List Streams Fails", fields: fields{code: http.StatusBadRequest, filename: "../test_data/streams/bad_req_list_streams.txt"}, args: args{limit: 3, offset: 1}, wantErr: true},
		{name: "Test List Streams Fails offset out of bounds", fields: fields{code: http.StatusBadRequest, filename: "../test_data/streams/bad_req_list_streams.txt"}, args: args{limit: 3, offset: 99999}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			maxListSize := 50
			client, server := testClientFileV2(tt.fields.code, tt.fields.filename)
			ctx := context.Background()
			defer server.Close()

			streamList, _, err := client.Streams.List(ctx, tt.args.limit, tt.args.offset)
			// Not expecting err
			if (err != nil) != tt.wantErr {
				t.Fatal(err)
			}

			// Expect err
			if err != nil && tt.wantErr {
				se, ok := err.(Error)
				if ok {
					if se.Status != tt.fields.code {
						t.Errorf("Expected HTTP %d, got %d", tt.fields.code, se.Status)
					}
					if se.Message != "domo err msg" {
						t.Error("Unexpected error message: ", se.Message)
					}
				}
			}
			// if streamList == nil {
			// 	t.Fatal("Got nil Streams")
			// }

			// Over max limit doesn't return more than max limit.
			if tt.args.limit > maxListSize && len(streamList) > maxListSize {
				t.Errorf("Expected list returned to be lte %d, go list size %d", maxListSize, len(streamList))
			}
			if len(streamList) > tt.args.limit {
				t.Errorf("expected lte streams than limit of %d, got %d ", tt.args.limit, len(streamList))
			}
		})
	}
}

func Test_ListSExecutions(t *testing.T) {
	type fields struct {
		code     int
		filename string
	}
	type args struct {
		streamID int
		limit    int
		offset   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		// want    []StreamDataset
		wantErr bool
	}{
		{name: "Test List Stream Executions", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_stream_executions.json"}, args: args{streamID: 1, limit: 1, offset: 0}, wantErr: false},
		{name: "Test List Stream Executions", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_stream_executions.json"}, args: args{streamID: 1, limit: 3, offset: 0}, wantErr: false},
		{name: "Test List Stream Executions", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_stream_executions.json"}, args: args{streamID: 1, limit: 3, offset: 1}, wantErr: false},
		{name: "Test List Stream Executions max limit", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_stream_executions.json"}, args: args{streamID: 1, limit: 50, offset: 1}, wantErr: false},
		{name: "Test List Stream Executions over max limit", fields: fields{code: http.StatusOK, filename: "../test_data/streams/list_stream_executions.json"}, args: args{streamID: 1, limit: 950, offset: 1}, wantErr: true},
		{name: "Test List Stream Executions Fails", fields: fields{code: http.StatusBadRequest, filename: "../test_data/streams/bad_req_list_streams.txt"}, args: args{streamID: 0, limit: 3, offset: 1}, wantErr: true},
		{name: "Test List Stream Executions Fails offset out of bounds", fields: fields{code: http.StatusBadRequest, filename: "../test_data/streams/bad_req_list_streams.txt"}, args: args{streamID: 1, limit: 3, offset: 99999}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// maxListSize := 50
			client, server := testClientFileV2(tt.fields.code, tt.fields.filename)
			ctx := context.Background()
			defer server.Close()

			streamList, _, err := client.Streams.ListExecutions(ctx, tt.args.streamID, tt.args.limit, tt.args.offset)
			// Not expecting err
			if (err != nil) != tt.wantErr {
				t.Fatal(err)
			}

			// Expect err
			if err != nil && tt.wantErr {
				se, ok := err.(Error)
				if ok {
					if se.Status != tt.fields.code {
						t.Errorf("Expected HTTP %d, got %d", tt.fields.code, se.Status)
					}
					if se.Message != "domo err msg" {
						t.Error("Unexpected error message: ", se.Message)
					}
				}
			}
			// if streamList == nil {
			// 	t.Fatal("Got nil Streams")
			// }

			// Over max limit doesn't return more than max limit.
			// if tt.args.limit > maxListSize && len(streamList) > maxListSize {
			// 	t.Errorf("Expected list returned to be lte %d, go list size %d", maxListSize, len(streamList))
			// }
			if len(streamList) > tt.args.limit {
				t.Errorf("expected lte streams than limit of %d, got %d ", tt.args.limit, len(streamList))
			}
		})
	}
}
func Test_DeleteStream(t *testing.T) {
	client, server := testClientStringV2(http.StatusOK, "")
	ctx := context.Background()
	defer server.Close()

	_, err := client.Streams.DeleteStream(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_DeleteStreamBadID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	_, err := client.Streams.DeleteStream(ctx, 0)
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}
func Test_CreateExecution(t *testing.T) {

	filename := "../test_data/streams/create_stream_execution.json"
	client, server := testClientFileV2(http.StatusOK, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.CreateExecution(ctx, 42)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Got nil Stream Details")
	}
	if res.ID != 1 {
		t.Error("Got wrong stream execution id")
	}
}
func Test_CreateExecutionBadID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.CreateExecution(ctx, 0)
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func Test_CommitExecution(t *testing.T) {

	filename := "../test_data/streams/commit_stream_execution.json"
	client, server := testClientFileV2(http.StatusOK, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.CommitExecution(ctx, 42, 1)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Got nil Stream Details")
	}
	if res.ID != 1 {
		t.Error("Got wrong execution stream id")
	}
}
func Test_CommitExecutionBadStreamID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.CommitExecution(ctx, 0, 0)
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}
func Test_CommitExecutionBadExecutionID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.CommitExecution(ctx, 0, 0)
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func Test_CommitExecutionBadIDs(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.CommitExecution(ctx, 0, 0)
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}
func Test_AbortExecution(t *testing.T) {

	client, server := testClientStringV2(http.StatusOK, "")
	ctx := context.Background()
	defer server.Close()

	_, err := client.Streams.AbortExecution(ctx, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_AbortStreamExecutionBadStreamID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	_, err := client.Streams.AbortExecution(ctx, 0, 0)
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}
func Test_AbortStreamExecutionBadExecutionID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	_, err := client.Streams.AbortExecution(ctx, 0, 0)
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func Test_AbortStreamExecutionBadIDs(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	_, err := client.Streams.AbortExecution(ctx, 0, 0)
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func Test_UploadDataPartStr(t *testing.T) {

	filename := "../test_data/streams/upload_data_part.json"
	client, server := testClientFileV2(http.StatusOK, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.UploadDataPartStr(ctx, 42, 1, 1, "csvData string")
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Got nil Stream Details")
	}
	if res.ID != 1 {
		t.Error("Got wrong stream")
	}
}
func Test_UploadDataPartBadStreamID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.UploadDataPartStr(ctx, 0, 0, 0, "csvData string")
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}
func Test_UploadDataPartBadExecutionID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.UploadDataPartStr(ctx, 0, 0, 0, "csvData string")
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}
func Test_UploadDataPartBadPartNumber(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.UploadDataPartStr(ctx, 0, 0, 0, "csvData string")
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func Test_ModifyStreamUpdateMethod(t *testing.T) {

	filename := "../test_data/streams/update_stream.json"
	client, server := testClientFileV2(http.StatusOK, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.ModifyStreamUpdateMethod(ctx, 42, true)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Got nil Stream Details")
	}
	if res.ID != 42 {
		t.Error("Got wrong stream")
	}
}

func Test_UpdateStreamMetaBadStreamID(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.ModifyStreamUpdateMethod(ctx, 0, true)
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func Test_CreateStreamBadSchema(t *testing.T) {

	filename := "../test_data/streams/bad_req_list_streams.txt"
	client, server := testClientFileV2(http.StatusBadRequest, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.CreateStream(ctx, StreamDatasetSchema{})
	if res != nil {
		t.Error("Unexpected Stream Execution returned, expected nil")
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected domo error, got", err)
	}
	if se.Status != 400 {
		t.Errorf("Expected HTTP 400, got %d. ", se.Status)
	}
	if se.Message != "domo err msg" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func Test_CreateStream(t *testing.T) {

	filename := "../test_data/streams/create_new_stream.json"
	client, server := testClientFileV2(http.StatusOK, filename)
	ctx := context.Background()
	defer server.Close()

	res, _, err := client.Streams.CreateStream(ctx, StreamDatasetSchema{})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Got nil Stream Details")
	}
	if res.ID != 42 {
		t.Error("Got wrong stream")
	}
}

func Test_UploadDataPartBadCsvColumnTypes(t *testing.T) {
	// bad LONG
	// bad DATE
	// bad DATETIME
	// bad DOUBLE
}
func Test_UploadDataPartBadCsvColumnNumbers(t *testing.T) {
	// more columns than schema
	// less columns than schema
}
