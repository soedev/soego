package main

import (
	"errors"

	"github.com/soedev/soego"
	"github.com/soedev/soego/core/elog"
)

func main() {
	err := soego.New().Invoker(func() error {
		elog.Info("logger info", elog.String("gopher", "ego"), elog.String("type", "command"))
		return errors.New("i am panic")
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
