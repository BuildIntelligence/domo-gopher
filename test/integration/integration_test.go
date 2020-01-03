package domo

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestListStreamsOffsetParamIsNotIgnoredByDomo(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domod {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		off0URL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 3, 0)
		off2URL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 3, 2)

		off0Resp, err := client.getRespBody(off0URL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		off2Resp, err := client.getRespBody(off2URL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if off0Resp == off2Resp {
			t.Error("Expected offset 0 to return a different response than offset 2")
		}

		listOff0, err := client.ListStreams(3, 0)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		listOff2, err := client.ListStreams(3, 2)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if listOff0[0].ID == listOff2[0].ID {
			t.Error("Expected first entry of Offset 0 list to have a different ID than first entry of Offset 2 list.")
		}
		if listOff0[2].ID != listOff2[0].ID {
			t.Error("Expected third entry of Offset 0 list to be the same as first entry of Offset 2 list.")
		}
	}
}

func TestListStreamsLimitParamIsNotIgnoredByDomo(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domod {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		lim1URL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 1, 0)
		lim5URL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 5, 0)

		lim1Resp, err := client.getRespBody(lim1URL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		lim5Resp, err := client.getRespBody(lim5URL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if lim1Resp == lim5Resp {
			t.Error("Expected limit 1 to return a different response than limit 5")
		}

		listLim1, err := client.ListStreams(1, 0)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		listLim5, err := client.ListStreams(5, 0)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if len(listLim1) == len(listLim5) {
			t.Error("Expected Limit 1 and Limit 5 to return different List Lengths")
		}

		if len(listLim1) != 1 {
			t.Errorf("Expected to return 1 stream, returned %d streams", len(listLim1))
		}
		if len(listLim5) != 5 {
			t.Errorf("Expected to return 5 stream, returned %d streams", len(listLim5))
		}
	}
}

func TestListStreamsSortParamIsNotIgnoredByDomo(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domod {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		sortAscURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&sort=%s", 5, 0, "name")
		sortDescURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&sort=%s", 5, 0, "-name")

		sortAscResp, err := client.getRespBody(sortAscURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		sortDescResp, err := client.getRespBody(sortDescURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if sortAscResp == sortDescResp {
			t.Error("Expected Sort by Name Asc and Sort by Name Desc to return different responses")
		}
	}
}

func TestListStreamsFieldsParamIsNotIgnoredByDomo(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domod {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		noFieldsParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 5, 0)
		datasetFieldParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&fields=%s", 5, 0, "dataSet")
		updateMethodFieldParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&fields=%s", 5, 0, "updateMethod")
		createdAtFieldParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&fields=%s", 5, 0, "createdAt")
		modifiedAtFieldParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&fields=%s", 5, 0, "modifiedAt")

		noFieldsParamResp, err := client.getRespBody(noFieldsParamURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		datasetFieldParamResp, err := client.getRespBody(datasetFieldParamURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		updateMethodFieldParamResp, err := client.getRespBody(updateMethodFieldParamURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		createdAtFieldParamResp, err := client.getRespBody(createdAtFieldParamURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		modifiedAtFieldParamResp, err := client.getRespBody(modifiedAtFieldParamURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if datasetFieldParamResp == updateMethodFieldParamResp {
			t.Error("Expected dataset field param and updateMethod field param to return different responses")
		}
		if datasetFieldParamResp == createdAtFieldParamResp {
			t.Error("Expected dataset field param and createdAt field param to return different responses")
		}
		if datasetFieldParamResp == modifiedAtFieldParamResp {
			t.Error("Expected dataset field param and modifiedAt field param to return different responses")
		}
		if updateMethodFieldParamResp == datasetFieldParamResp {
			t.Error("Expected updateMethod field param and createdAt field param to return different responses")
		}
		if updateMethodFieldParamResp == modifiedAtFieldParamResp {
			t.Error("Expected updateMethod field param and modifiedAt field param to return different responses")
		}
		if createdAtFieldParamResp == modifiedAtFieldParamResp {
			t.Error("Expected createdAt field param and modifiedAt field param to return different responses")
		}
		if noFieldsParamResp == datasetFieldParamResp {
			t.Error("Expected no fields param and dataset field param to return different responses")
		}
		if noFieldsParamResp == updateMethodFieldParamResp {
			t.Error("Expected no fields param and updateMethod field param to return different responses")
		}
		if noFieldsParamResp == createdAtFieldParamResp {
			t.Error("Expected no fields param and createdAt field param to return different responses")
		}
		if noFieldsParamResp == modifiedAtFieldParamResp {
			t.Error("Expected no fields param and modifiedAt field param to return different responses")
		}
	}
}

// https://api.domo.com/v1/streams?q=dataSource.owner.id:1704739518&offset=0&limit=5
func TestListStreamsQParamIsNotIgnoredByDomo(t *testing.T) {
	// if the -short flag is passed this will be skipped. Since this requires a flag to be
	// passed everytime I added an opt in flag to actually run these.
	// This -short flag will help out with the UI in some IDEs though so I have both despite the redundancy.
	if testing.Short() {
		t.Skip()
	}
	// Don't run these integration tests unless the "domo" flag is passed. i.e. `go test -domo`
	flag.Parse()
	if *domod {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
		}

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := NewAuthenticator(ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()

		noFieldsParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?offset=%d&limit=%d", 0, 5)

		ownerIDQueryURL := fmt.Sprintf("https://api.domo.com/v1/streams?q=dataSource.owner.id:%d&offset=%d&limit=%d", 1704739518, 0, 5)

		dataSourceNameQueryURL := fmt.Sprintf("https://api.domo.com/v1/streams?q=dataSource.name:%s&offset=%d&limit=%d", "Rusty", 0, 5)

		dataSourceIDQueryURL := fmt.Sprintf("https://api.domo.com/v1/streams?q=dataSource.id:%s&offset=%d&limit=%d", "59682470-ff7b-43c7-9024-a9deea824eb6", 0, 5)

		noFieldsParamResp, err := client.getRespBody(noFieldsParamURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		ownerIDQueryResp, err := client.getRespBody(ownerIDQueryURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		dataSourceNameQueryResp, err := client.getRespBody(dataSourceNameQueryURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		dataSourceIDQueryResp, err := client.getRespBody(dataSourceIDQueryURL)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if noFieldsParamResp == ownerIDQueryResp {
			t.Errorf("Expected no fields param and query by owner ID to return different responses \n%s\n%s", "", "") //noFieldsParamResp, ownerIDQueryResp)
		}

		if noFieldsParamResp == dataSourceNameQueryResp {
			t.Errorf("Expected no fields param and query by name to return different responses \n%s\n%s", "", "") //noFieldsParamResp, dataSourceNameQueryResp)
		}

		if noFieldsParamResp == dataSourceIDQueryResp {
			t.Errorf("Expected no fields param and query by datasource ID to return different responses \n%s\n%s", "", "") //noFieldsParamResp, dataSourceIDQueryResp)
		}
	}
}
