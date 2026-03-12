package main

import (
	"sync"

	"gofr.dev/pkg/gofr"
)

var (
	n  = 0
	mu sync.RWMutex
)

const duration = 3

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

	// runs every second
	app.AddCronJob("* * * * * *", "counter", count)

	app.Run()
}

func count(c *gofr.Context) {
	mu.Lock()
	defer mu.Unlock()

	n++

	c.Log("Count:", n)
}
