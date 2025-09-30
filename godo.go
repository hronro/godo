package main

import (
	_ "embed"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"godo/data"
	"godo/layout"
	"godo/templates"
)

//go:embed static/style.css
var styleCss []byte

//go:embed static/htmx.min.js
var htmxJs []byte

//go:embed static/idiomorph-ext.min.js
var idiomorphExtJs []byte

func main() {
	app := fiber.New()

	app.Get("/static/style.css", func(c fiber.Ctx) error {
		c.Type("css")
		return c.Send(styleCss)
	})

	app.Get("/static/htmx.min.js", func(c fiber.Ctx) error {
		c.Type("js")
		return c.Send(htmxJs)
	})

	app.Get("/static/idiomorph-ext.min.js", func(c fiber.Ctx) error {
		c.Type("js")
		return c.Send(idiomorphExtJs)
	})

	app.Get("/", func(c fiber.Ctx) error {
		todos := data.GetTodos()
		c.Type("html")
		return layout.Render(c, c, "Todo List", templates.App(todos))
	})

	app.Get("/todos", func(c fiber.Ctx) error {
		todos := data.GetTodos()
		c.Type("html")
		return templates.Todos(todos).Render(c, c)
	})

	app.Post("/todo", func(c fiber.Ctx) error {
		newTodoText := c.FormValue("text")
		if len(newTodoText) == 0 {
			return c.Status(fiber.StatusBadRequest).SendString("empty TODO text")
		}
		newTodoText = strings.Clone(newTodoText)

		data.AddTodo(newTodoText)
		todos := data.GetTodos()
		c.Type("html")
		return templates.Todos(todos).Render(c, c)
	})

	app.Get("/todo/:todoId", func(c fiber.Ctx) error {
		todoId, paramErr := strconv.Atoi(c.Params("todoId"))
		if paramErr != nil {
			return c.Status(fiber.StatusBadRequest).SendString(paramErr.Error())
		}

		todo := data.GetTodo(todoId)
		if todo == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		_, isEditing := c.Queries()["edit"]

		c.Type("html")
		return templates.Todo(*todo, isEditing).Render(c, c)
	})

	app.Patch("/todo/:todoId", func(c fiber.Ctx) error {
		todoId, paramErr := strconv.Atoi(c.Params("todoId"))
		if paramErr != nil {
			return c.Status(fiber.StatusBadRequest).SendString(paramErr.Error())
		}

		var text *string
		if t := strings.Clone(c.FormValue("text")); t != "" {
			text = &t
		} else {
			text = nil
		}
		var done *bool
		if d := strings.Clone(c.FormValue("done")); d != "" {
			var dBool bool
			switch d {
			case "true":
				dBool = true
			case "false":
				dBool = false
			default:
				return c.Status(fiber.StatusBadRequest).SendString("invalid done value")
			}
			done = &dBool
		} else {
			done = nil
		}

		updatedTodo := data.UpdateTodo(todoId, text, done)

		if updatedTodo == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		c.Type("html")
		return templates.Todo(*updatedTodo, false).Render(c, c)
	})

	app.Delete("/todo/:todoId", func(c fiber.Ctx) error {
		todoId, paramErr := strconv.Atoi(c.Params("todoId"))
		if paramErr != nil {
			return c.Status(fiber.StatusBadRequest).SendString(paramErr.Error())
		}

		deletedTodo := data.RemoveTodo(todoId)
		if deletedTodo == nil {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			c.SendStatus(fiber.StatusOK)
			return nil
		}
	})

	data.AddTodo("First Todo")
	data.AddTodo("Learn Go")
	data.AddTodo("Learn Fiber")
	data.AddTodo("Learn Templ")
	data.AddTodo("Learn HTMX")

	log.Fatal(app.Listen(":3000", fiber.ListenConfig{
		// EnablePrefork: true,
	}))

	log.Println(p())
}
