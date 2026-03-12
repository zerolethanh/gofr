package main

import (
	"errors"
	"fmt"

	"gofr.dev/examples/using-migrations/migrations"
	"gofr.dev/pkg/gofr"
)

const (
	queryGetEmployee    = "SELECT id,name,gender,contact_number,dob from employee where name = ?"
	queryInsertEmployee = "INSERT INTO employee (id, name, gender, contact_number,dob) values (?, ?, ?, ?, ?)"
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
	// Create a new application
	a := gofr.New()

	// Add migrations to run
	a.Migrate(migrations.All())

	// Add all the routes
	a.GET("/employee", GetHandler)
	a.POST("/employee", PostHandler)

	// Run the application
	a.Run()
}

type Employee struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Phone  int    `json:"contact_number"`
	DOB    string `json:"dob"`
}

// GetHandler handles GET requests for retrieving employee information
// @Summary get-handler
// @Description get-handler
// @Tags Using-Migrations
// @Param name query string true "name"
// @Success 200 {object} using-migrations.Employee
// @Failure 500 {object} error
// @Router /employee [GET]
func GetHandler(c *gofr.Context) (any, error) {
	name := c.Param("name")
	if name == "" {
		return nil, errors.New("name can't be empty")
	}

	var emp Employee

	err := c.SQL.QueryRowContext(c, queryGetEmployee, name).
		Scan(&emp.ID, &emp.Name, &emp.Gender, &emp.Phone, &emp.DOB)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("DB Error: %v", err))
	}

	return emp, nil
}

// PostHandler handles POST requests for creating new employees
// @Summary post-handler
// @Description post-handler
// @Tags Using-Migrations
// @Accept json
// @Param Employee body Employee false "Employee data"
// @Success 200 {object} string
// @Failure 500 {object} error
// @Router /employee [POST]
func PostHandler(c *gofr.Context) (any, error) {
	var emp Employee
	if err := c.Bind(&emp); err != nil {
		c.Logger.Errorf("error in binding: %v", err)
		return nil, errors.New("invalid body")
	}

	// Execute the INSERT query
	_, err := c.SQL.ExecContext(c, queryInsertEmployee, emp.ID, emp.Name, emp.Gender, emp.Phone, emp.DOB)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("DB Error: %v", err))
	}

	return fmt.Sprintf("successfully posted entity: %v", emp.Name), nil
}
