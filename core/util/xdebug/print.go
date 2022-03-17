package xdebug

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/soedev/soego/core/util/xcolor"
)

// MakeReqResInfo ...
// Deprecated: MakeReqResInfo will be removed in v0.10
func MakeReqResInfo(compName string, addr string, cost time.Duration, req interface{}, reply interface{}) string {
	return fmt.Sprintf("%s %s %s %s => %s\n", xcolor.Green(compName), xcolor.Green(addr), xcolor.Yellow(fmt.Sprintf("[%vms]", float64(cost.Microseconds())/1000)), xcolor.Blue(fmt.Sprintf("%v", req)), xcolor.Blue(fmt.Sprintf("%v", reply)))
}

// MakeReqResError ...
// Deprecated: MakeReqResError will be removed in v0.10
func MakeReqResError(compName string, addr string, cost time.Duration, req string, err string) string {
	return fmt.Sprintf("%s %s %s %s => %s\n", xcolor.Red(compName), xcolor.Red(addr), xcolor.Yellow(fmt.Sprintf("[%vms]", float64(cost.Microseconds())/1000)), xcolor.Blue(fmt.Sprintf("%v", req)), xcolor.Red(err))
}

// MakeReqResInfoV2 以info级别打印行号、配置名、目标地址、耗时、请求数据、响应数据
func MakeReqResInfoV2(callerSkip int, compName string, addr string, cost time.Duration, req interface{}, reply interface{}) string {
	_, file, line, _ := runtime.Caller(callerSkip)
	return fmt.Sprintf("%s %s %s %s %s => %s \n", xcolor.Green(file+":"+strconv.Itoa(line)), xcolor.Green(compName), xcolor.Green(addr), xcolor.Yellow(fmt.Sprintf("[%vms]", float64(cost.Microseconds())/1000)), xcolor.Blue(fmt.Sprintf("%v", req)), xcolor.Blue(fmt.Sprintf("%v", reply)))
}

// MakeReqResErrorV2 以error级别打印行号、配置名、目标地址、耗时、请求数据、响应数据
func MakeReqResErrorV2(callerSkip int, compName string, addr string, cost time.Duration, req string, err string) string {
	_, file, line, _ := runtime.Caller(callerSkip)
	return fmt.Sprintf("%s %s %s %s %s => %s \n", xcolor.Green(file+":"+strconv.Itoa(line)), xcolor.Red(compName), xcolor.Red(addr), xcolor.Yellow(fmt.Sprintf("[%vms]", float64(cost.Microseconds())/1000)), xcolor.Blue(fmt.Sprintf("%v", req)), xcolor.Red(err))
}
