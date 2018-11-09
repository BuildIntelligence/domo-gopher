// Example usage of domo-gopher
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	domo "gitlab.com/buildintelligence/domo-gopher"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, make sure you've created one in the same directory as main.go")
	}

	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")

	// New Domo Package V2
	auth := domo.NewAuthenticator(domo.ScopeData)
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
