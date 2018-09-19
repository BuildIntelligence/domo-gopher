package domo

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

const (
	DSName = "Test DomoGopher"
)

func TestGetDatasets(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
	}

	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()

	data, err := client.GetDatasets(5, 0)
	if err != nil {
		t.Errorf("Unexpected Error Retrieving Datasets: %s", err)
	}
	if len(data) != 5 {
		t.Errorf("Expected 5 datasets, got %d.", len(data))
	}
}

func TestCreateDataset(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
	}

	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	columns := []Column{Column{ColumnType: "STRING", Name: "Test Col String"}, Column{ColumnType: "STRING", Name: "Test Col String 2"}}
	ds := DatasetDetails{Name: DSName, Description: "TestDomoGopherDatasetCreate", Rows: 0, Schema: Schema{Columns: columns}}

	dataset, err := client.CreateDataset(ds)
	if err != nil {
		t.Errorf("Unexpected Error Creating Dataset: %s", err)
	}
	if len(dataset.ID) == 0 {
		t.Errorf("Expected to have a dataset id returned with more than 0 char. Got dataset id: %s", dataset.ID)
	}
	fmt.Printf("Dataset ID: %s \n", dataset.ID)
	if dataset.Name != DSName {
		t.Errorf("Expected created dataset to have the name: %s but got: %s", DSName, dataset.Name)
	}
}
func TestClient_GetDatasets(t *testing.T) {
	type args struct {
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		want    []*DatasetDetails
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetDatasets(tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetDatasets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetDatasets() = %v, want %v", got, tt.want)
			}
		})
	}
}
