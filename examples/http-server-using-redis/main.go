package main

import (
	"time"

	"gofr.dev/pkg/gofr"
)

const redisExpiryTime = 5

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
	app := gofr.New()

	// Add routes for Redis operations
	app.GET("/redis/{key}", RedisGetHandler)
	app.POST("/redis", RedisSetHandler)
	app.GET("/redis-pipeline", RedisPipelineHandler)

	// Register an OnStart hook to warm up a cache.
	// This runs before route registration as intended.
	app.OnStart(func(ctx *gofr.Context) error {
		ctx.Logger.Info("Warming up the cache...")

		// Example: Fetch some data and store it in Redis.
		// In a real app, this might come from a database or another service.
		cacheKey := "initial-data"
		cacheValue := "This is some data cached at startup."

		err := ctx.Redis.Set(ctx, cacheKey, cacheValue, 0).Err()
		if err != nil {
			ctx.Logger.Errorf("Failed to warm up cache: %v", err)
			return err // Return the error to halt startup if caching fails.
		}

		ctx.Logger.Info("Cache warmed up successfully!")

		return nil
	})

	// Run the application
	app.Run()
}

// RedisSetHandler sets a key-value pair in Redis using the Set Command.
// @Summary redis-set-handler
// @Description redis-set-handler
// @Tags Http-Server-Using-Redis
// @Param input body map[string]string false "body map[string]string data"
// @Success 200 {object} string
// @Failure 500 {object} error
// @Router /redis [POST]
func RedisSetHandler(c *gofr.Context) (any, error) {
	input := make(map[string]string)

	if err := c.Request.Bind(&input); err != nil {
		return nil, err
	}

	for key, value := range input {
		err := c.Redis.Set(c, key, value, redisExpiryTime*time.Minute).Err()
		if err != nil {
			return nil, err
		}
	}

	return "Successful", nil
}

// RedisGetHandler gets the value from Redis.
// @Summary redis-get-handler
// @Description redis-get-handler
// @Tags Http-Server-Using-Redis
// @Param key path string true "key"
// @Success 200 {object} map[string]string
// @Failure 500 {object} error
// @Router /redis/{key} [GET]
func RedisGetHandler(c *gofr.Context) (any, error) {
	key := c.PathParam("key")

	value, err := c.Redis.Get(c, key).Result()
	if err != nil {
		return nil, err
	}

	resp := make(map[string]string)
	resp[key] = value

	return resp, nil
}

// RedisPipelineHandler demonstrates using multiple Redis commands efficiently within a pipeline.
// @Summary redis-pipeline-handler
// @Description redis-pipeline-handler
// @Tags Http-Server-Using-Redis
// @Success 200 {object} []v9.Cmder
// @Failure 500 {object} error
// @Router /redis-pipeline [GET]
func RedisPipelineHandler(c *gofr.Context) (any, error) {
	pipe := c.Redis.Pipeline()

	// Add multiple commands to the pipeline
	pipe.Set(c, "testKey1", "testValue1", redisExpiryTime*time.Minute)
	pipe.Get(c, "testKey1")

	// Execute the pipeline and get results
	cmds, err := pipe.Exec(c)
	if err != nil {
		return nil, err
	}

	// Process or return the results of each command in the pipeline (implementation omitted for brevity)
	return cmds, nil
}
