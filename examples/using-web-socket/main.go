package main

import (
	"gofr.dev/pkg/gofr"
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

	app.WebSocket("/ws", WSHandler)

	app.Run()
}

// WSHandler ws-handler
// @Summary ws-handler
// @Description ws-handler
// @Tags Using-Web-Socket
// @Param message query string false "query string data"
// @Success 200 {object} string
// @Failure 500 {object} error
// @Router /ws-handler [GET]
func WSHandler(ctx *gofr.Context) (any, error) {
	var message string

	err := ctx.Bind(&message)
	if err != nil {
		ctx.Logger.Errorf("Error binding message: %v", err)
		return nil, err
	}

	ctx.Logger.Infof("Received message: %s", message)

	err = ctx.WriteMessageToSocket("Hello! GoFr")
	if err != nil {
		return nil, err
	}

	return message, nil
}
