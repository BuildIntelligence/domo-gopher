# Getting Started

## Example using packaged OAuth2 Client Helper
``` golang
	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_SECRET")

	auth := domo.NewAuthenticator(domo.ScopeData)
	auth.SetAuthInfo(clientID, clientSecret) // It defaults to using the same ENV Vars so this is potentially Optional 
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
```

## Example using your own http client with OAuth2 configured on it
``` golang
	client := domo.NewClient(YourClient)
	ctx := context.Background()
	data, _, err := client.Datasets.List(ctx, 5, 0)
	if err != nil {
		fmt.Println("error domo dataset")
	}
	for _, ds := range data {
		out := fmt.Sprintf("DomoDomo Dataset name: %s, ID: %s", ds.Name, ds.ID)
		fmt.Println(out)
	}
```

# TODO:
- [x] improve auth scope configuration to include scope in the url auth params based on input flags
- [x] Dataset API wrapper methods
- [x] Stream API wrapper methods
- [x] User API wrapper methods
- [x] Group API wrapper methods
- [x] Page API wrapper methods
- [x] Go Modules for dependency management
- [ ] Dataset/Stream upload methods that take an array/slice of structs. i.e. it handles the serialization to CSV as well as schema generation/updating.
- [ ] Projects & Tasks API
- [ ] Account API