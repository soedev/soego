package main

import (
	"errors"
	"fmt"

	"github.com/soedev/soego"
	"github.com/soedev/soego/core/elog"
	"github.com/soedev/soego/core/etrace"
	"github.com/soedev/soego/task/ejob"
	"go.uber.org/zap"
)

// export EGO_DEBUG=true && go run main.go --job=jobrunner  --config=config.toml
func main() {
	if err := soego.New().Job(
		ejob.Job("job1", job1),
		ejob.Job("job2", job2),
	).Run(); err != nil {
		elog.Error("start up", zap.Error(err))
	}
}

func job2(ctx ejob.Context) error {
	fmt.Println("i am error job runner, traceId: ", etrace.ExtractTraceID(ctx.Ctx))
	return errors.New("i am error")
}

func job1(ctx ejob.Context) error {
	fmt.Println("i am job runner, traceId: ", etrace.ExtractTraceID(ctx.Ctx))
	return nil
}
