package internal

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kiran94/terraform-provider-example/pkg"
	"github.com/sirupsen/logrus"
)

// Creates a new Terrform Provider Schema
// Resources in the map follow a naming convention
// in this case having a prefix 'planner_'
// This is important on the consumer side where terraform
// will use that prefix to identify which provider to use
// when the user declares resources
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HOST", "http://localhost:8080"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"planner_todo_note": resourceTodoNote(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// Configures the providerConfigure
// The parameters provided in the above schema are accessible
// and the returned client is usable in the rest of the resources
func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	host := d.Get("host").(string)
	logrus.WithField("host", host).Debug("using host")

	apiClient := pkg.TodoApiClient{
		BaseUrl:    host,
		HttpClient: http.DefaultClient,
	}
	return &apiClient, nil
}
