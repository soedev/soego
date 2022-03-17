package main

import (
	"context"
	"fmt"

	"github.com/soedev/soego"
	"github.com/soedev/soego/client/ehttp"
	"github.com/soedev/soego/core/elog"
	"github.com/soedev/soego/core/etrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	if err := soego.New().Invoker(
		invokerHTTP,
		callHTTP,
	).Run(); err != nil {
		elog.Error("startup", elog.FieldErr(err))
	}
}

var httpComp *ehttp.Component

func invokerHTTP() error {
	httpComp = ehttp.Load("http.test").Build()
	return nil
}

func callHTTP() error {
	tracer := etrace.NewTracer(trace.SpanKindClient)

	req := httpComp.R()

	ctx, span := tracer.Start(context.Background(), "callHTTP()", propagation.HeaderCarrier(req.Header))
	defer span.End()

	// Inject traceId Into Header
	//c1 := etrace.HeaderInjector(ctx, req.Header)
	fmt.Println(span.SpanContext().TraceID())
	info, err := req.SetContext(ctx).Get("http://127.0.0.1:9007/hello")
	if err != nil {
		return err
	}
	fmt.Println(info)
	return nil
}
