package domo

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"gitlab.com/buildintelligence/domo-gopher/domo"
	"os"
	"testing"

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

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := domo.NewAuthenticator(domo.ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()
		ctx := context.Background()

		off0URL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 3, 0)
		off2URL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 3, 2)

		rq1, err := client.NewRequest("GET", off0URL, nil)
		off0Buf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq1, off0Buf)
		off0Resp := off0Buf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq2, err := client.NewRequest("GET", off2URL, nil)
		off2Buf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq2, off2Buf)
		off2Resp := off2Buf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if off0Resp == off2Resp {
			t.Error("Expected offset 0 to return a different response than offset 2")
		}

		listOff0, _, err := client.Streams.List(ctx, 3, 0)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		listOff2, _, err := client.Streams.List(ctx, 3, 2)
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

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := domo.NewAuthenticator(domo.ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()
		ctx := context.Background()

		lim1URL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 1, 0)
		lim5URL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 5, 0)

		rq1, err := client.NewRequest("GET", lim1URL, nil)
		lim1Buf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq1, lim1Buf)
		lim1Resp := lim1Buf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq2, err := client.NewRequest("GET", lim5URL, nil)
		lim5Buf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq2, lim5Buf)
		lim5Resp := lim5Buf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		if lim1Resp == lim5Resp {
			t.Error("Expected limit 1 to return a different response than limit 5")
		}

		listLim1, _, err := client.Streams.List(ctx, 1, 0)
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		listLim5, _, err := client.Streams.List(ctx, 5, 0)
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

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := domo.NewAuthenticator(domo.ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()
		ctx := context.Background()

		sortAscURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&sort=%s", 5, 0, "name")
		sortDescURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&sort=%s", 5, 0, "-name")

		rq1, err := client.NewRequest("GET", sortAscURL, nil)
		ascBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq1, ascBuf)
		sortAscResp := ascBuf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq2, err := client.NewRequest("GET", sortDescURL, nil)
		descBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq2, descBuf)
		sortDescResp := descBuf.String()
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

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := domo.NewAuthenticator(domo.ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()
		ctx := context.Background()
		noFieldsParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d", 5, 0)
		datasetFieldParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&fields=%s", 5, 0, "dataSet")
		updateMethodFieldParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&fields=%s", 5, 0, "updateMethod")
		createdAtFieldParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&fields=%s", 5, 0, "createdAt")
		modifiedAtFieldParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?limit=%d&offset=%d&fields=%s", 5, 0, "modifiedAt")

		rq, err := client.NewRequest("GET", noFieldsParamURL, nil)
		nfpBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq, nfpBuf)
		noFieldsParamResp := nfpBuf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq1, err := client.NewRequest("GET", datasetFieldParamURL, nil)
		dfpBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq1, dfpBuf)
		datasetFieldParamResp := dfpBuf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq2, err := client.NewRequest("GET", updateMethodFieldParamURL, nil)
		umfpBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq2, umfpBuf)
		updateMethodFieldParamResp := umfpBuf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq3, err := client.NewRequest("GET", createdAtFieldParamURL, nil)
		cafpBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq3, cafpBuf)
		createdAtFieldParamResp := cafpBuf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq4, err := client.NewRequest("GET", modifiedAtFieldParamURL, nil)
		mafpBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq4, mafpBuf)
		modifiedAtFieldParamResp := mafpBuf.String()
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

		clientID := os.Getenv("DOMO_CLIENT_ID")
		clientSecret := os.Getenv("DOMO_SECRET")
		auth := domo.NewAuthenticator(domo.ScopeData)
		auth.SetAuthInfo(clientID, clientSecret)
		client := auth.NewClient()
		ctx := context.Background()

		noFieldsParamURL := fmt.Sprintf("https://api.domo.com/v1/streams?offset=%d&limit=%d", 0, 5)

		ownerIDQueryURL := fmt.Sprintf("https://api.domo.com/v1/streams?q=dataSource.owner.id:%d&offset=%d&limit=%d", 1704739518, 0, 5)

		dataSourceNameQueryURL := fmt.Sprintf("https://api.domo.com/v1/streams?q=dataSource.name:%s&offset=%d&limit=%d", "Rusty", 0, 5)

		dataSourceIDQueryURL := fmt.Sprintf("https://api.domo.com/v1/streams?q=dataSource.id:%s&offset=%d&limit=%d", "59682470-ff7b-43c7-9024-a9deea824eb6", 0, 5)


		rq1, err := client.NewRequest("GET", noFieldsParamURL, nil)
		nfpBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq1, nfpBuf)
		noFieldsParamResp := nfpBuf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq2, err := client.NewRequest("GET", ownerIDQueryURL, nil)
		oiqBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq2, oiqBuf)
		ownerIDQueryResp := oiqBuf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq3, err := client.NewRequest("GET", dataSourceNameQueryURL, nil)
		dsnqBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq3, dsnqBuf)
		dataSourceNameQueryResp := dsnqBuf.String()
		if err != nil {
			t.Errorf("Unexpected Error Retrieving Streams List: %s", err)
		}

		rq4, err := client.NewRequest("GET", dataSourceIDQueryURL, nil)
		dsiqBuf := new(bytes.Buffer)
		_, err = client.Do(ctx, rq4, dsiqBuf)
		dataSourceIDQueryResp := dsiqBuf.String()
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
