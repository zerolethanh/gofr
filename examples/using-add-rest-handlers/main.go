package main

import (
	"gofr.dev/examples/using-add-rest-handlers/migrations"
	"gofr.dev/pkg/gofr"
)

type user struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	IsEmployed bool   `json:"isEmployed"`
}

// GetAll : User can overwrite the specific handlers by implementing them like this
// @Summary get-all
// @Description get-all
// @Tags Using-Add-Rest-Handlers
// @Param id path string true "id"
// @Success 200 {object} string
// @Router /{id}/get-all [GET]
func (u *user) GetAll(c *gofr.Context) (any, error) {
	return "user GetAll called", nil
}

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
	// Create a new application
	a := gofr.New()

	// Add migrations to run
	a.Migrate(migrations.All())

	// AddRESTHandlers creates CRUD handles for the given entity
	err := a.AddRESTHandlers(&user{})
	if err != nil {
		return
	}

	// Run the application
	a.Run()
}
