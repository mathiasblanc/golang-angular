package todo

import (
	"errors"
	"sync"

	"github.com/rs/xid"
)

var (
	list []Todo
	mtx  sync.RWMutex
	once sync.Once
)

func init() {
	once.Do(initializeList)
}

func initializeList() {
	list = []Todo{}
}

//Todo data structure for a task with a description of what to do
type Todo struct {
	Id       string `json:"id"`
	Message  string `json:"message"`
	Complete bool   `json:"complete"`
}

//Get retrieves all elements from the todo list
func Get() []Todo {
	return list
}

//Add adds a new todo based on a message
func Add(message string) string {
	todo := newTodo(message)
	mtx.Lock()
	list = append(list, todo)
	mtx.Unlock()

	return todo.Id
}

//Delete removes a todo from the todo list
func Delete(id string) error {
	location, err := findTodoLocation(id)

	if err != nil {
		return err
	}

	removeTodoByLocation(location)
	return nil
}

//Complete sets the complete flag to true, marking a todo as completed
func Complete(id string) error {
	location, err := findTodoLocation(id)

	if err != nil {
		return err
	}

	setTodoCompleteByLocation(location)
	return nil
}

func newTodo(message string) Todo {
	return Todo{
		Id:       xid.New().String(),
		Message:  message,
		Complete: false,
	}
}

func findTodoLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()

	for i, t := range list {
		if isMatchingId(t.Id, id) {
			return i, nil
		}
	}

	return 0, errors.New("Could not find todo based on id")
}

func removeTodoByLocation(location int) {
	mtx.Lock()
	list = append(list[:location], list[location+1:]...)
	mtx.Unlock()
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
}

func isMatchingId(a, b string) bool {
	return a == b
}
