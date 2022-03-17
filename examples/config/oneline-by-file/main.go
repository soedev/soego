package main

import (
	"github.com/soedev/soego"
	"github.com/soedev/soego/core/econf"
	"github.com/soedev/soego/core/elog"
)

// export EGO_DEBUG=true && go run main.go  --config=config.toml --watch=false
func main() {
	if err := soego.New().Invoker(func() error {
		peopleName := econf.GetString("people.name")
		elog.Info("people info", elog.String("name", peopleName), elog.String("type", "onelineByFile"))
		return nil
	}).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
