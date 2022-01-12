package domo

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"
)

const (
	testDSName = "Test DomoGopher"
)

func Test_checkForSchemaChangeByColumnNameMatching(t *testing.T) {
	domo := Schema{Columns: []Column{{
		ColumnType: "DECIMAL",
		Name:       "Blah",
	}, {
		ColumnType: "DATE",
		Name:       "firstBlahDay",
	}, {
		ColumnType: "DATETIME",
		Name:       "firstBlahTime",
	}, {
		ColumnType: "STRING",
		Name:       "Foo",
	}, {
		ColumnType: "LONG",
		Name:       "bar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "baz",
	}, {
		ColumnType: "LONG",
		Name:       "BazBar",
	}, {
		ColumnType: "LONG",
		Name:       "obar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "obaz",
	}}}

	sample := GenerateDataSetSchema(reflect.TypeOf(DomoEmbeddedSample{}))

	diffs := checkForSchemaChangeByColumnNameMatching(sample, domo)

	if diffs.DiffsCount() > 0 {
		t.Errorf("Unexpected Difference In Schema: %v", diffs)
	}
}

func Test_checkForSchemaChangeByColumnNameMatching_addColumn(t *testing.T) {
	domo := Schema{Columns: []Column{{
		ColumnType: "DECIMAL",
		Name:       "Blah",
	}, {
		ColumnType: "DATE",
		Name:       "firstBlahDay",
	}, {
		ColumnType: "DATETIME",
		Name:       "firstBlahTime",
	}, {
		ColumnType: "STRING",
		Name:       "Foo",
	}, {
		ColumnType: "LONG",
		Name:       "bar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "baz",
	}, {
		ColumnType: "LONG",
		Name:       "BazBar",
	}, {
		ColumnType: "LONG",
		Name:       "obar",
	}}}

	sample := GenerateDataSetSchema(reflect.TypeOf(DomoEmbeddedSample{}))

	diffs := checkForSchemaChangeByColumnNameMatching(sample, domo)

	if diffs.DiffsCount() != 1 {
		t.Errorf("Expected only 1 Difference In Schema, found %d: %v", diffs.DiffsCount(), diffs)
	}
	colToAdd := diffs.ColumnsToAddToDomo[0]
	if colToAdd.ComparedColumnName != "obaz" {
		t.Errorf("Expected to have column named %s to add, found %s", "obaz", colToAdd.ComparedColumnName)
	}
	if colToAdd.ComparedColumnType != ColumnTypeDouble {
		t.Errorf("Expected to have column with type %s to add, found type %s", ColumnTypeDouble, colToAdd.ComparedColumnType)
	}
}

func Test_checkForSchemaChangeByColumnNameMatching_addColumn_andChangeDataType(t *testing.T) {
	domo := Schema{Columns: []Column{{
		ColumnType: "DECIMAL",
		Name:       "Blah",
	}, {
		ColumnType: "DATE",
		Name:       "firstBlahDay",
	}, {
		ColumnType: "DATETIME",
		Name:       "firstBlahTime",
	}, {
		ColumnType: "STRING",
		Name:       "Foo",
	}, {
		ColumnType: "LONG",
		Name:       "bar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "baz",
	}, {
		ColumnType: "STRING",
		Name:       "BazBar",
	}, {
		ColumnType: "LONG",
		Name:       "obar",
	}}}

	sample := GenerateDataSetSchema(reflect.TypeOf(DomoEmbeddedSample{}))

	diffs := checkForSchemaChangeByColumnNameMatching(sample, domo)

	if diffs.DiffsCount() != 2 {
		t.Errorf("Expected only 1 Difference In Schema, found %d: %v", diffs.DiffsCount(), diffs)
	}
	colToAdd := diffs.ColumnsToAddToDomo[0]
	if colToAdd.ComparedColumnName != "obaz" {
		t.Errorf("Expected to have column named %s to add, found %s", "obaz", colToAdd.ComparedColumnName)
	}
	if colToAdd.ComparedColumnType != ColumnTypeDouble {
		t.Errorf("Expected to have column with type %s to add, found type %s", ColumnTypeDouble, colToAdd.ComparedColumnType)
	}
	colToChangeType := diffs.ColumnTypeMismatch[0]
	if colToChangeType.DomoColumnName != "BazBar" {
		t.Errorf("Expected to have a column data type change for column BazBar. Found column data type change for %s", colToChangeType.DomoColumnName)
	}
	if colToChangeType.ComparedColumnType == colToChangeType.DomoColumnType {
		t.Errorf("Expected to have a mismatch between ComparedColumnType and DomoColumnType but both where %s", colToChangeType.ComparedColumnType)
	}
}

func Test_checkForSchemaChangeByColumnNameMatching_deleteColumn(t *testing.T) {
	domo := Schema{Columns: []Column{{
		ColumnType: "DECIMAL",
		Name:       "Blah",
	}, {
		ColumnType: "DATE",
		Name:       "firstBlahDay",
	}, {
		ColumnType: "DATETIME",
		Name:       "firstBlahTime",
	}, {
		ColumnType: "STRING",
		Name:       "Foo",
	}, {
		ColumnType: "LONG",
		Name:       "bar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "baz",
	}, {
		ColumnType: "LONG",
		Name:       "BazBar",
	}, {
		ColumnType: "LONG",
		Name:       "obar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "obaz",
	}, {
		ColumnType: "STRING",
		Name:       "ColumnInDomoToDelete",
	}}}

	sample := GenerateDataSetSchema(reflect.TypeOf(DomoEmbeddedSample{}))

	diffs := checkForSchemaChangeByColumnNameMatching(sample, domo)

	if diffs.DiffsCount() != 1 {
		t.Errorf("Expected only 1 Difference In Schema, found %d: %v", diffs.DiffsCount(), diffs)
	}
	colToDelete := diffs.ColumnsToDeleteFromDomo[0]
	if colToDelete.DomoColumnName != "ColumnInDomoToDelete" {
		t.Errorf("Expected to have column named %s to add, found %s", "ColumnInDomoToDelete", colToDelete.DomoColumnName)
	}
	if colToDelete.DomoColumnType != ColumnTypeString {
		t.Errorf("Expected to have column with type %s to add, found type %s", ColumnTypeString, colToDelete.DomoColumnType)
	}
}

func Test_checkForSchemaChangeByColumnIndexComparision(t *testing.T) {
	domo := Schema{Columns: []Column{{
		ColumnType: "DECIMAL",
		Name:       "Blah",
	}, {
		ColumnType: "DATE",
		Name:       "firstBlahDay",
	}, {
		ColumnType: "DATETIME",
		Name:       "firstBlahTime",
	}, {
		ColumnType: "STRING",
		Name:       "Foo",
	}, {
		ColumnType: "LONG",
		Name:       "bar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "baz",
	}, {
		ColumnType: "LONG",
		Name:       "BazBar",
	}, {
		ColumnType: "LONG",
		Name:       "obar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "obaz",
	}}}

	sample := GenerateDataSetSchema(reflect.TypeOf(DomoEmbeddedSample{}))

	diffs := checkForSchemaChangeByColumnIndexComparision(sample, domo)

	if diffs.DiffsCount() > 0 {
		t.Errorf("Unexpected Difference In Schema: %v", diffs)
	}

}

func Test_checkForSchemaChangeByColumnIndexComparision_Name_andType_Diff(t *testing.T) {
	domo := Schema{Columns: []Column{{
		ColumnType: "DECIMAL",
		Name:       "Blah",
	}, {
		ColumnType: "DATE",
		Name:       "firstBlahDay",
	}, {
		ColumnType: "DATETIME",
		Name:       "firstBlahTime",
	}, {
		ColumnType: "STRING",
		Name:       "Foo",
	}, {
		ColumnType: "LONG",
		Name:       "bar",
	}, {
		ColumnType: "DOUBLE",
		Name:       "baz",
	}, {
		ColumnType: "LONG",
		Name:       "BazBar",
	}, {
		ColumnType: "LONG",
		Name:       "obar",
	}, {
		ColumnType: "DECIMAL",
		Name:       "OBizzle",
	}}}

	sample := GenerateDataSetSchema(reflect.TypeOf(DomoEmbeddedSample{}))

	diffs := checkForSchemaChangeByColumnIndexComparision(sample, domo)

	if diffs.DiffsCount() != 2 {
		t.Errorf("Unexpected number of Differences In Schema. Expected 2, found %d:\n%v", diffs.DiffsCount(), diffs)
	}
	nameDiff := diffs.NameMismatch[0]
	if nameDiff.DomoColumnName == nameDiff.ComparedColumnName {
		t.Errorf("Expected column names to be different. Found %s and %s", nameDiff.DomoColumnType, nameDiff.ComparedColumnName)
	}
	typeDiff := diffs.ColumnTypeMismatch[0]
	if typeDiff.DomoColumnType == typeDiff.ComparedColumnType && typeDiff.ComparedColumnType != ColumnTypeDouble {
		t.Errorf("Expected to have a change to %s. Found %s", ColumnTypeDouble, typeDiff.ComparedColumnType)
	}
}

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

	sqlQuery := fmt.Sprintf("SELECT %s FROM table WHERE `GL Account No.`=%s", "`GL Account No.`, `Category`", "1000")
	data, _, err := client.Datasets.QueryData(ctx, "5a60b561-7030-42fa-a8f6-d04ddbe89864", sqlQuery)
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
		ds := Dataset{Name: testDSName, Description: "TestDomoGopherDatasetCreate", Rows: 0, Schema: Schema{Columns: columns}}
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

func TestDatasetsService_DownloadDatasetCSV(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	// if testing.Short() {
	// 	t.Skip()
	// }
	// Don't run these integration tests unless the "domoGopher" flag is passed. i.e. `go test -domo`
	// flag.Parse()
	// if *domogopher {
	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	auth := NewAuthenticator(ScopeData)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	ctx := context.Background()

	// ds_id := "8ef52283-a1ee-4890-a17d-241b149f2d9f"
	ds_id2 := "b6d07391-0887-4247-9b1e-43089559816f"
	dataset, _, err := client.Datasets.DownloadDatasetCSV(ctx, ds_id2, false)
	if err != nil {
		t.Errorf("Unexpected Error Creating Dataset: %s", err)
	}
	if len(dataset) == 0 {
		t.Errorf("Expected to have a dataset CSV returned with more than 0 length. Got dataset id: %s", dataset)
	}
	fmt.Printf("Dataset:\n%s\n", dataset)
	// }
}
