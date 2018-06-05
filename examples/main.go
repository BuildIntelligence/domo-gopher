// Example usage of domo-gopher
package main

import (
	domoGopher "../domo-gopher"
	"github.com/bitly/go-simplejson"
	"fmt"
)
func main() {

	// Create a new domo obj
	domo := domoGopher.New("clientID","clientSecret")

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
}

