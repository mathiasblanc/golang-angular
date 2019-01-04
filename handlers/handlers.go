package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mathiasblanc/golang-angular/todo"
)

//GetTodoListHandler returns all todo list items
func GetTodoListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, todo.Get())
}

//AddTodoHandler adds a todo item to the todo list
func AddTodoHandler(c *gin.Context) {
	item, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)

	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	c.JSON(statusCode, gin.H{"id": todo.Add(item.Message)})
}

//DeleteTodoHandler deletes a specified todo based on user http input
func DeleteTodoHandler(c *gin.Context) {
	id := c.Param("id")

	if err := todo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

//CompleteTodoHandler sets a todo as completed
func CompleteTodoHandler(c *gin.Context) {
	item, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)

	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	if todo.Complete(item.Id) != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToTodo(httpBody io.ReadCloser) (todo.Todo, int, error) {
	body, err := ioutil.ReadAll(httpBody)

	if err != nil {
		return todo.Todo{}, http.StatusInternalServerError, err
	}

	defer httpBody.Close()
	return convertJSONBodyToTodo(body)
}

func convertJSONBodyToTodo(jsonBody []byte) (todo.Todo, int, error) {
	var item todo.Todo
	err := json.Unmarshal(jsonBody, &item)

	if err != nil {
		return todo.Todo{}, http.StatusBadRequest, err
	}

	return item, http.StatusOK, nil
}
