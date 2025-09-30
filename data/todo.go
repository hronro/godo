package data

import (
	"log"
	"sync"
)

type Todo struct {
	Id   int
	Text string
	Done bool
}

type Todos struct {
	mu     sync.Mutex
	nextId int
	todos  []Todo
}

var todos = Todos{}

// Do read-only operations for the returned value,
// any modifications are not allowed.
func GetTodos() []Todo {
	return todos.todos
}

func AddTodo(text string) {
	todos.mu.Lock()
	defer todos.mu.Unlock()
	id := todos.nextId
	todos.todos = append(todos.todos, Todo{Id: id, Text: text, Done: false})
	todos.nextId++
	log.Printf("Added todo: %s (ID: %d)", text, id)
}

func GetTodo(id int) *Todo {
	for _, todo := range todos.todos {
		if todo.Id == id {
			return &todo
		}
	}
	return nil
}

func UpdateTodo(id int, text *string, done *bool) *Todo {
	todos.mu.Lock()
	defer todos.mu.Unlock()
	for index := range todos.todos {
		if todos.todos[index].Id == id {
			if text != nil {
				todos.todos[index].Text = *text
			}
			
			if done != nil {
				todos.todos[index].Done = *done
			}

			log.Printf("Updating todo (ID: %d): text=%v, done=%v", id, text, done)

			return &todos.todos[index]
		}
	}
	return nil
}

func RemoveTodo(id int) *Todo {
	todos.mu.Lock()
	defer todos.mu.Unlock()

	for index, todo := range todos.todos {
		if todo.Id == id {
			todos.todos = append(todos.todos[:index], todos.todos[index+1:]...)
			log.Printf("Deleted todo: %s (ID: %d)", todo.Text, todo.Id)
			return &todo
		}
	}

	return nil
}
