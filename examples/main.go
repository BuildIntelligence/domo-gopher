// Example usage of domo-gopher
package main

import (
	"context"
	"fmt"
	"gitlab.com/buildintelligence/domo-gopher/domo"
	"os"
)

func main() {

	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")

	// New Domo Package V2
	auth := domo.NewAuthenticator(domo.ScopeData)
	auth.SetAuthInfo(clientID, clientSecret)
	client := auth.NewClient()
	ctx := context.Background()
	data, _, err := client.Datasets.List(ctx, 5, 0)
	if err != nil {
		fmt.Println("error domo dataset")
	}
	for _, ds := range data {
		out := fmt.Sprintf("DomoDomo Dataset name: %s, ID: %s", ds.Name, ds.ID)
		fmt.Println(out)
	}
}
