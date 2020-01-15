package domo

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
)

const (
	testDSName = "Test DomoGopher"
)

func TestDatasetsService_List(t *testing.T) {
	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	ctx := context.Background()

	data, _, err := client.Datasets.List(ctx, 5, 0)
	if err != nil {
		t.Errorf("Unexpected Error Retrieving Datasets: %s", err)
	}
	if len(data) != 5 {
		t.Errorf("Expected 5 datasets, got %d.", len(data))
	}
}

func TestDatasetsService_QueryData(t *testing.T) {
	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	ctx := context.Background()

	sqlQuery := fmt.Sprintf("SELECT %s FROM table WHERE `WMS Community`='%s'", "`Lot No. for Community Maps`, `Sales Status Series for Community Map`, `Sales Status`", "ABM")
	data, _, err := client.Datasets.QueryData(ctx, "447a2858-9c1c-42a9-b90b-a5340268d90e", sqlQuery)
	if err != nil {
		t.Errorf("Unexpected Error Retrieving Data: %s", err)
	}
	fmt.Println(data)
	if len(data) == 0 {
		t.Errorf("Expected data, got %s.", data)
	}
}

func TestDatasetsService_Create(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domoGopher" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domogopher {
		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()
		ctx := context.Background()

		columns := []Column{Column{ColumnType: "STRING", Name: "Test Col String"}, Column{ColumnType: "STRING", Name: "Test Col String 2"}}
		ds := DatasetDetails{Name: testDSName, Description: "TestDomoGopherDatasetCreate", Rows: 0, Schema: Schema{Columns: columns}}

		dataset, _, err := client.Datasets.Create(ctx, ds)
		if err != nil {
			t.Errorf("Unexpected Error Creating Dataset: %s", err)
		}
		if len(dataset.ID) == 0 {
			t.Errorf("Expected to have a dataset id returned with more than 0 char. Got dataset id: %s", dataset.ID)
		}
		fmt.Printf("Dataset ID: %s \n", dataset.ID)
		if dataset.Name != testDSName {
			t.Errorf("Expected created dataset to have the name: %s but got: %s", testDSName, dataset.Name)
		}
	}
}
