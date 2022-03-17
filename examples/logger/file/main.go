package main

import (
	"github.com/soedev/soego"
	"github.com/soedev/soego/core/elog"
)

//  export EGO_DEBUG=false && go run main.go
func main() {
	err := soego.New().Invoker(func() error {
		elog.Info("logger info", elog.String("gopher", "ego"), elog.String("type", "command"))
		return nil
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
