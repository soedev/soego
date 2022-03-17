package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/soedev/soego"
	"github.com/soedev/soego/core/elog"
	"github.com/soedev/soego/examples/helloworld"
	"github.com/soedev/soego/server/egin"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := soego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()
		server.GET("/hello", func(ctx *gin.Context) {
			ctx.JSON(200, "Hello client: "+ctx.GetHeader("app"))
			return
		})
		mock := &GreeterMock{}
		server.GET("/grpcproxyok", egin.GRPCProxy(mock.SayHelloOK))
		server.GET("/grpcproxyerr", egin.GRPCProxy(mock.SayHelloErr))
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

type GreeterMock struct{}

func (mock GreeterMock) SayHelloOK(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{
		Message: "hello",
	}, nil
}

func (mock GreeterMock) SayHelloErr(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{
		Message: "hello",
	}, fmt.Errorf("say hello err")
}
