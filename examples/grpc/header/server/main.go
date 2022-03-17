package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/soedev/soego"
	"github.com/soedev/soego/core/elog"
	"github.com/soedev/soego/examples/helloworld"
	"github.com/soedev/soego/server"
	"github.com/soedev/soego/server/egrpc"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := soego.New().Serve(func() server.Server {
		server := egrpc.Load("server.grpc").Build()
		helloworld.RegisterGreeterServer(server.Server, &Greeter{server: server})
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

// Greeter ...
type Greeter struct {
	server *egrpc.Component
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g Greeter) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	if request.Name == "error" {
		return nil, status.Error(codes.Unavailable, "error")
	}
	header := metadata.Pairs("x-header-key", "val")
	err := grpc.SendHeader(context, header)
	if err != nil {
		return nil, fmt.Errorf("set header fail, %w", err)
	}
	return &helloworld.HelloResponse{
		Message: "Hello EGO, I'm " + g.server.Address(),
	}, nil
}
