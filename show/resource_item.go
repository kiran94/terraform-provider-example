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
				Type:        schema.TypeInt,
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
				Type:        schema.TypeInt,
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

	show := map[string]interface{}{
		"id":     d.Get(unique_id_key).(int),
		"name":   d.Get(name_key).(string),
		"rating": d.Get(rating_key).(int),
	}

	log.Printf("Creating Resource %s", show["id"])
	err := apiClient.PostItem(show)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(d.Get(unique_id_key).(int)))
	d.Set(name_key, show["name"])
	d.Set(rating_key, show["rating"])

	log.Printf("Created Item %d", d.Get(unique_id_key).(int))
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.ApiClient)

	itemId := d.Get(unique_id_key).(int)
	log.Printf("Reading Item %d", itemId)
	item, err := apiClient.GetItem(itemId)

	if err != nil {
		return err
	}

	if len(item) == 0 {
		log.Println("No items were found.")
		return nil
	}

	d.SetId(strconv.Itoa(int(item["id"].(float64))))
	d.Set(rating_key, strconv.Itoa(int(item["rating"].(float64))))
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

	intId, err2 := strconv.Atoi(itemId)
	if err2 != nil {
		return err2
	}

	err := apiClient.DeleteItem(intId)
	if err != nil {
		return err
	}

	return nil
}

func resourceExistsItem(d *schema.ResourceData, m interface{}) (bool, error) {
	itemId := d.Id()
	log.Printf("Checking if item exists: %s", itemId)

	apiClient := m.(*client.ApiClient)

	intItemId, err := strconv.Atoi(itemId)
	if err != nil {
		return false, err
	}

	item, err := apiClient.GetItem(intItemId)
	if err != nil {
		return false, err
	}

	return len(item) != 0, nil
}
