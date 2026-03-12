package main

import (
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
)

// @title Swa App
// @version 1.0
// @description This is a sample server.
// @termOfService https://swagger.io/terms/
// @schemes http https
// @contact.name API Support
// @contact.url https://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html
// @basePath /
func main() {
	app := gofr.New()
	app.GET("/list", listHandler)
	app.AddStaticFiles("/", "./static")
	app.Run()
}

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

// listHandler list-handler
// @Summary list-handler
// @Description list-handler
// @Tags Using-Html-Template
// @Success 200 {object} response.Template
// @Router /list [GET]
func listHandler(*gofr.Context) (any, error) {
	// Get data from somewhere
	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Expand on Gofr documentation ", Done: false},
			{Title: "Add more examples", Done: true},
			{Title: "Write some articles", Done: false},
		},
	}

	return response.Template{Data: data, Name: "todo.html"}, nil
}
