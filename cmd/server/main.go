package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/kiran94/terraform-provider-example/pkg"
	ginlogrus "github.com/toorop/gin-logrus"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	port  string = ":8080"
	local bool   = true
)

// A Simulated Database
var (
	todoDatabase = map[int]pkg.Todo{
		1: {Id: 1, Title: "My Todo", Message: "My Message of the Todo", Priority: 2},
		2: {Id: 2, Title: "My Todo", Message: "My Message of the Todo", Priority: 4},
		3: {Id: 3, Title: "My Todo", Message: "My Message of the Todo", Priority: 8},
		4: {Id: 4, Title: "My Todo", Message: "My Message of the Todo", Priority: 9},
		5: {Id: 5, Title: "My Todo", Message: "My Message of the Todo", Priority: 1},
	}

	databaseIncrementMutex sync.Mutex
	databaseIncrement      = 5
)

/*
Starts a new Example Todo Planner Server which is
used as the underlying API for the Terraform Provider.

This is a just a simple todo application which
provides CRUD functionality.
*/
func main() {
	r := getGinEngine()

	// Get a todo from the database
	r.GET("/todo/:id", func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok || id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
			return
		}

		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("id %s could not be converted into number", id)})
			return
		}

		todo, ok := todoDatabase[idParsed]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("id %s could not be found in db", id)})
			return
		}

		c.Header("Content-Type", gin.MIMEJSON)
		c.JSON(http.StatusOK, todo)
	})

	// Creates a new todo
	r.POST("/todo", func(c *gin.Context) {
		defer c.Request.Body.Close()

		allBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			sendBadRequest(c, http.StatusBadRequest, "could not read body", err)
			return
		}

		var newTodo pkg.Todo
		if err := json.Unmarshal(allBody, &newTodo); err != nil {
			sendBadRequest(c, http.StatusBadRequest, "could not deserialize body", err)
			return
		}

		newTodo.Id = incrementDatabaseId()
		todoDatabase[newTodo.Id] = newTodo

		logrus.WithField("id", newTodo.Id).Debug("Created new todo")
		c.JSON(http.StatusAccepted, newTodo)
	})

	// Updates an existing todo
	r.PATCH("/todo/:id", func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok || id == "" {
			sendBadRequest(c, http.StatusBadRequest, "id parameter is required", nil)
			return
		}

		idParsed, err := strconv.Atoi(id)
		if err != nil {
			sendBadRequest(c, http.StatusBadRequest, fmt.Sprintf("id %s could not be converted into number", id), err)
			return
		}

		defer c.Request.Body.Close()
		allBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			sendBadRequest(c, http.StatusBadRequest, "could not read body", err)
			return
		}

		_, ok = todoDatabase[idParsed]
		if !ok {
			logrus.WithField("id", id).Warn("todo already exists. creating new entry")
			sendBadRequest(c, http.StatusBadRequest, "todo already exists", nil)
			return
		}

		var updatedTodo pkg.Todo
		if err := json.Unmarshal(allBody, &updatedTodo); err != nil {
			sendBadRequest(c, http.StatusBadRequest, "could not deserialize body", err)
			return
		}

		todoDatabase[idParsed] = updatedTodo
		c.Status(http.StatusAccepted)
	})

	r.DELETE("/todo/:id", func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok || id == "" {
			sendBadRequest(c, http.StatusBadRequest, "id parameter is required", nil)
			return
		}

		idParsed, err := strconv.Atoi(id)
		if err != nil {
			sendBadRequest(c, http.StatusBadRequest, fmt.Sprintf("id %s could not be converted into number", id), err)
			return
		}

		_, ok = todoDatabase[idParsed]
		if !ok {
			sendBadRequest(c, http.StatusBadGateway, fmt.Sprintf("id %s was not found", id), nil)
			return
		}

		delete(todoDatabase, idParsed)
		c.Status(http.StatusAccepted)
	})

	logrus.WithField("port", port).Info("Listening...")
	r.Run(port)
}

// Sends a bad request with the given message
func sendBadRequest(c *gin.Context, status int, message string, err error) {
	logrus.WithFields(logrus.Fields{
		"status":  status,
		"message": message,
	}).WithError(err).Error("Sending bad request")

	c.JSON(status, gin.H{"error": message})
}

// Gets a new gin Engine configured according to the environment
func getGinEngine() *gin.Engine {
	var r *gin.Engine

	if local {
		r = gin.Default()
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.InfoLevel)

		r = gin.New()
		r.Use(ginlogrus.Logger(logrus.New()), gin.Recovery())
	}

	return r
}

// Atomically increment the database id across the server
// not multi-container safe
func incrementDatabaseId() int {
	databaseIncrementMutex.Lock()
	databaseIncrement += 1
	databaseIncrementMutex.Unlock()
	return databaseIncrement
}
