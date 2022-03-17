package main

import (
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"

	"github.com/soedev/soego/core/econf"
	"github.com/soedev/soego/core/elog"
)

func main() {
	var err error
	conf := `
[stderr]
level = "info"
enableAsync = false
writer = "stderr"
`
	if err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal); err != nil {
		log.Println("load conf fail", err)
		return
	}
	logger := elog.Load("stderr").Build()
	logger.Info("an stderr msg", zap.Any("lee", 17))
}
