package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/kiran94/terraform-provider-example/internal"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: internal.Provider,
	})
}
