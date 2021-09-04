package show

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kiran94/terraform-provider-example/client"
)

const (
	unique_id_key = "unique_id"
	name_key      = "name"
	rating_key    = "rating"
)

func resourceTvShow() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			unique_id_key: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier for the show",
				ForceNew:    true,
				// ValidateFunc: validateName,
			},
			name_key: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the show",
				ForceNew:    false,
				// ValidateFunc: validateName,
			},
			rating_key: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rating of the show",
				ForceNew:    false,
				// ValidateFunc: validateName,
			},
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceUpdateItem,
		Delete: resourceDeleteItem,
		Exists: resourceExistsItem,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.ApiClient)

	show := map[string]string{
		"id":     d.Get(unique_id_key).(string),
		"name":   d.Get(name_key).(string),
		"rating": d.Get(rating_key).(string),
	}

	log.Printf("Creating Resource %s", show["id"])
	err := apiClient.PostItem(show)
	if err != nil {
		return err
	}

	d.SetId(d.Get(unique_id_key).(string))
	d.Set(name_key, show["name"])
	d.Set(rating_key, show["rating"])

	log.Printf("Created Item %s", d.Get(unique_id_key).(string))
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.ApiClient)

	itemId := d.Get(unique_id_key).(string)
	log.Printf("Reading Item %s", itemId)
	item, err := apiClient.GetItem(itemId)

	if err != nil {
		return err
	}

	if len(item) == 0 {
		log.Println("No items were found.")
		return nil
	}

	d.SetId(strconv.FormatFloat(item["id"].(float64), 'f', 0, 64))
	d.Set(rating_key, strconv.FormatFloat(item["rating"].(float64), 'f', 0, 64))
	d.Set(name_key, item["name"].(string))
	return nil
}

func resourceUpdateItem(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating Resource")
	return resourceCreateItem(d, m)
}

func resourceDeleteItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.ApiClient)

	itemId := d.Id()
	log.Printf("Deleting Item %s", itemId)

	err := apiClient.DeleteItem(itemId)
	if err != nil {
		return err
	}

	return nil
}

func resourceExistsItem(d *schema.ResourceData, m interface{}) (bool, error) {
	itemId := d.Id()
	log.Printf("Checking if item exists: %s", itemId)

	apiClient := m.(*client.ApiClient)

	item, err := apiClient.GetItem(itemId)
	if err != nil {
		return false, err
	}

	return len(item) != 0, nil
}
