package client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type ApiClient struct {
	BaseUrl string
}

func (client *ApiClient) Init(baseUrl string) {
	client.BaseUrl = baseUrl
}

func (client *ApiClient) GetItem(id string) (map[string]interface{}, error) {
	url := client.BaseUrl + "?id=" + id

	log.Printf("Getting Item %s", url)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(response.Body).Decode(&result)
	return result, nil
}

func (client *ApiClient) PostItem(show map[string]string) error {
	url := client.BaseUrl
	log.Printf("Creating Item %s (%s)", url, show)

	requestBody, err := json.Marshal(show)
	if err != nil {
		return err
	}

	result, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	log.Printf("Recieved Http Response %s", result)
	return nil
}

func (client *ApiClient) DeleteItem(id string) error {

	url := client.BaseUrl + "?id=" + id
	log.Printf("Deleting Item %s", url)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	log.Printf("Recieved Http Response %s", response)
	return err
}
