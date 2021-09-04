package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/kiran94/terraform-provider-example/example"
	// "log"
	// "github.com/kiran94/terraform-provider-example/client"
)

func main() {
    // testClient()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: example.Provider,
	})
}

func testClient() {
	apiClient := new(client.ApiClient)
	apiClient.Init("https://localhost:5001")

	show := map[string]string{
		"id":     "12",
		"name":   "Toy Story 12",
		"rating": "100",
	}

	Add Item
	err := apiClient.PostItem(show)
	if err != nil {
		panic(err)
	}

	Get Item
	result, err := apiClient.GetItem(show["id"])
	if err != nil {
		panic(err)
	} else {
		log.Println(result)
	}

	if len(result) == 0 {
		log.Println("No Result found")
	}

	Delete Item
	err2 := apiClient.DeleteItem(show["id"])
	if err2 != nil {
		panic(err2)
	}
}
