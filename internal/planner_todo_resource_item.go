package internal

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kiran94/terraform-provider-example/pkg"
)

// Creates a new Terraform Resource object for todo items.
func resourceTodoNote() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the todo note",
				ForceNew:    false,
			},
			"message": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The body of the todo note",
				ForceNew:    false,
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "The priority of the todo note",
				ForceNew:    false,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 10 {
						errs = append(errs, fmt.Errorf("priority was %d but should be between 1-10", v))
					}

					return
				},
			},
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceUpdateItem,
		Delete: resourceDeleteItem,
		Exists: resourceExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Creates a new Todo Resource
func resourceCreateItem(rd *schema.ResourceData, i interface{}) error {
	apiClient := i.(*pkg.TodoApiClient)
	todo, err := createTodoItemFromResource(rd, false)
	if err != nil {
		return err
	}

	if err := apiClient.Create(todo); err != nil {
		return err
	}

	id := strconv.Itoa(todo.Id)
	rd.SetId(id)
	rd.Set("title", todo.Title)
	rd.Set("message", todo.Message)
	rd.Set("priority", todo.Priority)
	return nil
}

// Reads a Todo Resource
func resourceReadItem(rd *schema.ResourceData, i interface{}) error {
	apiClient := i.(*pkg.TodoApiClient)
	idStr := rd.Id()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	todo, err := apiClient.Get(id)
	if err != nil {
		return err
	}

	rd.SetId(idStr)
	rd.Set("title", todo.Title)
	rd.Set("message", todo.Message)
	rd.Set("priority", todo.Priority)
	return nil
}

// Updates the Todo Resource
func resourceUpdateItem(rd *schema.ResourceData, i interface{}) error {
	apiClient := i.(*pkg.TodoApiClient)

	todo, err := createTodoItemFromResource(rd, true)
	if err != nil {
		return err
	}

	if err := apiClient.Update(todo); err != nil {
		return err
	}

	return resourceCreateItem(rd, i)
}

// Deletes the Todo Resource
func resourceDeleteItem(rd *schema.ResourceData, i interface{}) error {
	apiClient := i.(*pkg.TodoApiClient)
	id := rd.Id()
	idParsed, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return apiClient.Delete(idParsed)
}

// Checks if the Todo Resource exists
func resourceExists(rd *schema.ResourceData, i interface{}) (bool, error) {
	apiClient := i.(*pkg.TodoApiClient)
	id := rd.Id()
	idParsed, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}

	todo, err := apiClient.Get(idParsed)
	todoExists := todo != nil
	return todoExists, err
}

// Creates a todo structure from the incoming resource data
func createTodoItemFromResource(rd *schema.ResourceData, extractId bool) (*pkg.Todo, error) {

	var id int
	if extractId {
		idStr := rd.Id()
		i, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}

		id = i
	}

	priority, ok := rd.Get("priority").(int)
	if !ok {
		return nil, fmt.Errorf("priority %s cannot be converted into an int", rd.Get("priority"))
	}

	todo := pkg.Todo{
		Id:       id,
		Title:    rd.Get("title").(string),
		Message:  rd.Get("message").(string),
		Priority: int8(priority),
	}

	return &todo, nil
}
