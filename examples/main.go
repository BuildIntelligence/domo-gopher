// Example usage of domo-gopher
package main

import (
	"fmt"
	"log"
	"os"

	domodomo ".."
	domoGopher "../domo-gopher"
	"github.com/bitly/go-simplejson"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
	}

	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")
	// Start of original Domo-Gopher
	// Create a new domo obj
	domo := domoGopher.New(clientID, clientSecret)

	// Authorize against Domo
	authorized, _ := domo.Authorize()
	if authorized {

		// get a dataset meta data, 'Api Test Web Form' dataset_id: 77faea51-68ab-4dd3-ae1a-8992bc1b58a8
		response, _ := domo.Get("datasets/%s", nil, "77faea51-68ab-4dd3-ae1a-8992bc1b58a8")

		// Parse resonse to JSON object and get dataset name
		json, _ := simplejson.NewJson(response)
		jsonData, exists := json.CheckGet("name")

		if exists {
			datasetName, _ := jsonData.String()
			fmt.Println("Dataset name is ", datasetName)
		} else {
			fmt.Println("Didn't work :(")
		}

	}
	// End of original domo-gopher

	// New Domo Package V2
	auth := domodomo.NewAuthenticator(domodomo.ScopeData)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	data, err := client.GetDatasets(5, 0)
	if err != nil {
		fmt.Println("error domo dataset")
	}
	for _, ds := range data {
		out := fmt.Sprintf("DomoDomo Dataset name: %s, ID: %s", ds.Name, ds.ID)
		fmt.Println(out)
	}
}
