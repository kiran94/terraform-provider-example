package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

type HttpAccess interface {
	Do(req *http.Request) (*http.Response, error)
}

type TodoApiClient struct {
	HttpClient HttpAccess
	BaseUrl    string
}

type TodoApi interface {
	Create(todo *Todo) error
	Update(todo *Todo) error
	Get(id int) (*Todo, error)
	Delete(id int) error
}

const (
	contentTypeJson = "application/json"
)

// Creates a new Todo Item
func (t TodoApiClient) Create(todo *Todo) error {
	if todo == nil {
		return errors.New("todo was nil")
	}

	logrus.WithField("id", todo.Id).Debug("Creating Todo Item")

	url, err := url.Parse(t.BaseUrl + "/todo")
	if err != nil {
		return err
	}

	r, err := t.sendRequest(url, http.MethodPost, todo)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status code was not accepted %d", r.StatusCode)
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, &todo); err != nil {
		return err
	}

	return nil

}

// Updates the given todo
func (t TodoApiClient) Update(todo *Todo) error {
	logrus.WithField("id", todo.Id).Debug("Updating Todo Item")

	url, err := url.Parse(t.BaseUrl + "/todo/" + strconv.Itoa(todo.Id))
	if err != nil {
		return err
	}

	response, err := t.sendRequest(url, http.MethodPatch, todo)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return errors.New("status code was not accepted")
	}

	return nil
}

// Gets the Todo with the given identifier
func (t TodoApiClient) Get(id int) (*Todo, error) {
	logrus.WithField("id", id).Debug("Getting Todo Item")

	url, err := url.Parse(t.BaseUrl + "/todo/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	response, err := t.sendRequest(url, http.MethodGet, nil)
	if err != nil {
		fmt.Println("error sending request")
		return nil, err
	}
	defer response.Body.Close()

	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading body")
		return nil, err
	}

	var recieveTodo Todo
	if err := json.Unmarshal(rawBody, &recieveTodo); err != nil {
		return nil, err
	}

	return &recieveTodo, nil
}

// Deletes the given Todo.
func (t TodoApiClient) Delete(id int) error {
	logrus.WithField("id", id).Debug("Deleting Todo Item")

	url, err := url.Parse(t.BaseUrl + "/todo/" + strconv.Itoa(id))
	if err != nil {
		return err
	}

	response, err := t.sendRequest(url, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// Executes a new request to the underlying data store.
func (t TodoApiClient) sendRequest(url *url.URL, method string, todo *Todo) (*http.Response, error) {
	logrus.WithField("method", method).Info("Sending Todo Reqeuest")

	headers := make(http.Header, 1)
	headers["content-type"] = []string{contentTypeJson}

	request := &http.Request{
		Method: method,
		URL:    url,
		Header: headers,
	}

	if todo != nil {
		rawBytes, err := json.Marshal(todo)
		if err != nil {
			return nil, err
		}

		br := bytes.NewReader(rawBytes)
		bc := io.NopCloser(br)
		defer bc.Close()

		request.Body = bc
	}

	logrus.WithFields(logrus.Fields{
		"method": request.Method,
		"url":    request.URL,
		"header": request.Header,
	}).Info("Fetching Todo from Api")

	r, err := t.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return r, nil
}
