package main

import (
	"github.com/soedev/soego"
	"github.com/soedev/soego/core/elog"
)

//  export EGO_DEBUG=true && go run main.go
func main() {
	err := soego.New().Invoker(func() error {
		elog.Info("logger info", elog.String("gopher", "ego"), elog.String("type", "command"), elog.Any("aaa", map[string]interface{}{"aa": "bb"}))
		return nil
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
