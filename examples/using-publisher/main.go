package main

import (
	"encoding/json"

	"gofr.dev/examples/using-publisher/migrations"
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

	app.Migrate(migrations.All())

	app.POST("/publish-order", order)
	app.POST("/publish-product", product)

	app.Run()
}

// order order
// @Summary order
// @Description order
// @Tags Using-Publisher
// @Accept json
// @Param orderStatus body orderStatus false "orderStatus data"
// @Success 200 {object} string
// @Failure 500 {object} error
// @Router /publish-order [POST]
func order(ctx *gofr.Context) (any, error) {
	type orderStatus struct {
		OrderId string `json:"orderId"`
		Status  string `json:"status"`
	}

	var data orderStatus

	err := ctx.Bind(&data)
	if err != nil {
		return nil, err
	}

	msg, _ := json.Marshal(data)

	err = ctx.GetPublisher().Publish(ctx, "order-logs", msg)
	if err != nil {
		return nil, err
	}

	return "Published", nil
}

// product product
// @Summary product
// @Description product
// @Tags Using-Publisher
// @Accept json
// @Param productInfo body productInfo false "productInfo data"
// @Success 200 {object} string
// @Failure 500 {object} error
// @Router /publish-product [POST]
func product(ctx *gofr.Context) (any, error) {
	type productInfo struct {
		ProductId string `json:"productId"`
		Price     string `json:"price"`
	}

	var data productInfo

	err := ctx.Bind(&data)
	if err != nil {
		return nil, err
	}

	msg, _ := json.Marshal(data)

	err = ctx.GetPublisher().Publish(ctx, "products", msg)
	if err != nil {
		return nil, err
	}

	return "Published", nil
}
