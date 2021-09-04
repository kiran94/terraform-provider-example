package show

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kiran94/terraform-provider-example/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HOST", "https://localhost:5001"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"example_tv_show": resourceTvShow(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	apiClient := new(client.ApiClient)
	apiClient.Init(host)
	return apiClient, nil
}
