package main

import (
	"github.com/soedev/soego"
	"github.com/soedev/soego/core/econf"
	"github.com/soedev/soego/core/elog"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml --watch=false
func main() {
	err := soego.New().Invoker(func() error {
		p := People{}
		err := econf.UnmarshalKey("people", &p)
		if err != nil {
			panic(err.Error())
		}
		elog.Info("people info", elog.String("name", p.Name), elog.String("type", "structByFile"))
		return nil
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

// People ...
type People struct {
	Name string
}
