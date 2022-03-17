package ehttp

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/soedev/soego/core/eapp"
	"github.com/soedev/soego/core/elog"
	"github.com/soedev/soego/core/emetric"
	"github.com/soedev/soego/core/etrace"
	"github.com/soedev/soego/core/util/xdebug"
)

type interceptor func(name string, cfg *Config, logger *elog.Component) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook)

func logAccess(name string, config *Config, logger *elog.Component, req *resty.Request, res *resty.Response, err error) {
	rr := req.RawRequest
	var url, host string
	// 修复err 不是 *resty.ResponseError错误的时候，可能为nil
	if rr != nil {
		url = rr.URL.RequestURI()
		host = rr.URL.Host
	}

	fullMethod := req.Method + "." + url // GET./hello
	var cost = time.Since(beg(req.Context()))
	var respBody string
	if res != nil {
		respBody = string(res.Body())
	}
	if eapp.IsDevelopmentMode() {
		if err != nil {
			log.Println("http.response", xdebug.MakeReqResErrorV2(6, name, config.Addr, cost, fullMethod, err.Error()))
		} else {
			log.Println("http.response", xdebug.MakeReqResInfoV2(6, name, config.Addr, cost, fullMethod, respBody))
		}
	}

	var fields = make([]elog.Field, 0, 15)
	fields = append(fields,
		elog.FieldMethod(fullMethod),
		elog.FieldName(name),
		elog.FieldCost(cost),
		elog.FieldAddr(host),
	)

	// 开启了链路，那么就记录链路id
	if config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
		fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(req.Context())))
	}

	if config.EnableAccessInterceptorRes {
		fields = append(fields, elog.FieldValueAny(respBody))
	}

	if config.SlowLogThreshold > time.Duration(0) && cost > config.SlowLogThreshold {
		logger.Warn("slow", fields...)
	}

	if err != nil {
		fields = append(fields, elog.FieldEvent("error"), elog.FieldErr(err))
		if res == nil {
			// 无 res 的是连接超时等系统级错误
			logger.Error("access", fields...)
			return
		}
		logger.Warn("access", fields...)
		return
	}

	if config.EnableAccessInterceptor {
		fields = append(fields, elog.FieldEvent("normal"))
		logger.Info("access", fields...)
	}
}

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
// https://blog.golang.org/context#TOC_3.2.
// https://golang.org/pkg/context/#WithValue ，这边文章说明了用struct，可以避免分配
type begKey struct{}

func beg(ctx context.Context) time.Time {
	begTime, _ := ctx.Value(begKey{}).(time.Time)
	return begTime
}

func fixedInterceptor(name string, config *Config, logger *elog.Component) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	return func(cli *resty.Client, req *resty.Request) error {
		req.SetContext(context.WithValue(req.Context(), begKey{}, time.Now()))
		return nil
	}, nil, nil
}

func logInterceptor(name string, config *Config, logger *elog.Component) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	afterFn := func(cli *resty.Client, response *resty.Response) error {
		logAccess(name, config, logger, response.Request, response, nil)
		return nil
	}
	errorFn := func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			logAccess(name, config, logger, req, v.Response, v.Err)
		} else {
			logAccess(name, config, logger, req, nil, err)
		}
	}
	return nil, afterFn, errorFn
}

func metricInterceptor(name string, config *Config, logger *elog.Component) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	addr := strings.TrimRight(config.Addr, "/")
	afterFn := func(cli *resty.Client, res *resty.Response) error {
		method := res.Request.Method + "." + res.Request.RawRequest.URL.Path
		emetric.ClientHandleCounter.Inc(emetric.TypeHTTP, name, method, addr, http.StatusText(res.StatusCode()))
		emetric.ClientHandleHistogram.Observe(res.Time().Seconds(), emetric.TypeHTTP, name, method, addr)
		return nil
	}
	errorFn := func(req *resty.Request, err error) {
		method := req.Method + "." + req.RawRequest.URL.Path
		if v, ok := err.(*resty.ResponseError); ok {
			emetric.ClientHandleCounter.Inc(emetric.TypeHTTP, name, method, addr, http.StatusText(v.Response.StatusCode()))
		} else {
			emetric.ClientHandleCounter.Inc(emetric.TypeHTTP, name, method, addr, "biz error")
		}

		emetric.ClientHandleHistogram.Observe(time.Since(beg(req.Context())).Seconds(), emetric.TypeHTTP, name, method, addr)
	}
	return nil, afterFn, errorFn
}

func traceInterceptor(name string, config *Config, logger *elog.Component) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	tracer := etrace.NewTracer(trace.SpanKindClient)
	beforeFn := func(cli *resty.Client, req *resty.Request) error {
		ctx, span := tracer.Start(req.Context(), req.Method, nil)

		span.SetAttributes(
			etrace.String("peer.service", name),
			etrace.String("http.method", req.Method),
			etrace.String("http.url", req.URL),
		)

		req.SetContext(ctx)
		return nil
	}

	afterFn := func(cli *resty.Client, res *resty.Response) error {
		span := trace.SpanFromContext(res.Request.Context())
		span.SetAttributes(
			etrace.String("http.status_code", cast.ToString(res.StatusCode())),
		)

		span.End()
		return nil
	}
	errorFn := func(req *resty.Request, err error) {
		span := trace.SpanFromContext(req.Context())

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}
	return beforeFn, afterFn, errorFn
}
