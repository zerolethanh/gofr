// versions:
// 	gofr-cli v0.6.0
// 	gofr.dev v1.37.0
// 	source: hello.proto

package server

import (
	"fmt"

	"gofr.dev/pkg/gofr"
)

// Register the gRPC service in your app using the following code in your main.go:
//
// server.RegisterHelloServerWithGofr(app, &server.NewHelloGoFrServer())
//
// HelloGoFrServer defines the gRPC server implementation.
// Customize the struct with required dependencies and fields as needed.

type HelloGoFrServer struct {
	health *healthServer
}

// SayHello say-hello
// @Summary say-hello
// @Description say-hello
// @Tags Server
// @Accept json
// @Param HelloRequest body HelloRequest false "HelloRequest data"
// @Success 200 {object} server.HelloResponse
// @Failure 500 {object} error
// @Router /say-hello [GET]
func (s *HelloGoFrServer) SayHello(ctx *gofr.Context) (any, error) {
	request := HelloRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return nil, err
	}

	name := request.Name
	if name == "" {
		name = "World"
	}

	//Performing HealthCheck
	//res, err := s.health.Check(ctx, &grpc_health_v1.HealthCheckRequest{
	//	Service: "Hello",
	//})
	//ctx.Log(res.String())

	// Setting the serving status
	//s.health.SetServingStatus(ctx, "Hello", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	return &HelloResponse{
		Message: fmt.Sprintf("Hello %s!", name),
	}, nil
}
