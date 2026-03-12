package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
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

	//HTTP service with default health check endpoint
	a.AddHTTPService("anotherService", "http://localhost:9000")

	// Add all the routes
	a.GET("/hello", HelloHandler)
	a.GET("/error", ErrorHandler)
	a.GET("/redis", RedisHandler)
	a.GET("/trace", TraceHandler)
	a.GET("/mysql", MysqlHandler)

	// Run the application
	a.Run()
}

// HelloHandler hello-handler
// @Summary hello-handler
// @Description hello-handler
// @Tags Http-Server
// @Param name query string true "name"
// @Success 200 {object} string
// @Router /hello [GET]
func HelloHandler(c *gofr.Context) (any, error) {
	name := c.Param("name")
	if name == "" {
		c.Log("Name came empty")
		name = "World"
	}

	return fmt.Sprintf("Hello %s!", name), nil
}

// ErrorHandler error-handler
// @Summary error-handler
// @Description error-handler
// @Tags Http-Server
// @Failure 500 {object} error
// @Router /error [GET]
func ErrorHandler(c *gofr.Context) (any, error) {
	return nil, errors.New("some error occurred")
}

// RedisHandler redis-handler
// @Summary redis-handler
// @Description redis-handler
// @Tags Http-Server
// @Success 200 {object} string
// @Failure 500 {object} datasource.ErrorDB
// @Router /redis [GET]
func RedisHandler(c *gofr.Context) (any, error) {
	val, err := c.Redis.Get(c, "test").Result()
	if err != nil && err != redis.Nil { // If key is not found, we are not considering this an error and returning "".
		return nil, datasource.ErrorDB{Err: err, Message: "error from redis db"}
	}

	return val, nil
}

// TraceHandler trace-handler
// @Summary trace-handler
// @Description trace-handler
// @Tags Http-Server
// @Success 200 {object} object
// @Failure 500 {object} error
// @Router /trace [GET]
func TraceHandler(c *gofr.Context) (any, error) {
	defer c.Trace("traceHandler").End()

	span2 := c.Trace("some-sample-work")
	// Waiting for 1ms to simulate workload
	<-time.After(time.Millisecond * 1) //nolint:wsl
	defer span2.End()

	// Ping redis 5 times concurrently and wait.
	count := 5
	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			c.Redis.Ping(c)
			wg.Done()
		}()
	}
	wg.Wait()

	// Call to Another service
	resp, err := c.GetHTTPService("anotherService").Get(c, "redis", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data = struct {
		Data any `json:"data"`
	}{}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data.Data, nil
}

// MysqlHandler mysql-handler
// @Summary mysql-handler
// @Description mysql-handler
// @Tags Http-Server
// @Success 200 {object} integer
// @Failure 500 {object} datasource.ErrorDB
// @Router /mysql [GET]
func MysqlHandler(c *gofr.Context) (any, error) {
	var value int
	err := c.SQL.QueryRowContext(c, "select 2+2").Scan(&value)
	if err != nil {
		return nil, datasource.ErrorDB{Err: err, Message: "error from sql db"}
	}

	return value, nil
}
