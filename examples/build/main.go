package main

import (
	"github.com/gin-gonic/gin"
	"github.com/soedev/soego"
	"github.com/soedev/soego/core/elog"
	"github.com/soedev/soego/server/egin"
	"github.com/soedev/soego/server/egovernor"
)

func main() {
	if err := soego.New().
		Serve(
			egovernor.Load("server.governor").Build(),
			serverHTTP(),
		).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

func serverHTTP() *egin.Component {
	server := egin.Load("server.http").Build()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, "Hello")
		return
	})
	return server
}
