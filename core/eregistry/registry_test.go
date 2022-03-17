package eregistry

import (
	"testing"
	"time"

	"github.com/soedev/soego/core/constant"
	"github.com/soedev/soego/core/util/xtime"
	"github.com/soedev/soego/server"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceKey(t *testing.T) {
	svc := server.ServiceInfo{
		Name:    "myserver",
		Scheme:  "http",
		Address: "localhost",
		Weight:  100,
		Enable:  true,
		Healthy: true,
		Kind:    constant.ServiceProvider,
		Metadata: map[string]string{
			"appHost":    "",
			"appMode":    "",
			"appVersion": "",
			"buildTime":  "",
			"egoVersion": "v0.7.0",
			"key":        "val",
			"startTime":  xtime.TS.Format(time.Now()),
		},
	}
	assert.Equal(t, "/ego/myserver/providers/http://localhost", GetServiceKey("ego", &svc))
}

func TestGetServiceValue(t *testing.T) {
	svc := server.ServiceInfo{
		Name:    "myserver",
		Scheme:  "http",
		Address: "localhost",
		Weight:  100,
		Enable:  true,
		Healthy: true,
		Kind:    constant.ServiceProvider,
		Metadata: map[string]string{
			"appHost":    "",
			"appMode":    "",
			"appVersion": "",
			"buildTime":  "",
			"egoVersion": "v0.7.0",
			"key":        "val",
			"startTime":  xtime.TS.Format(time.Now()),
		},
	}
	assert.Contains(t, GetServiceValue(&svc), "v0.7.0")
	assert.Contains(t, GetServiceValue(&svc), "localhost")
	assert.Contains(t, GetServiceValue(&svc), "myserver")
	assert.Contains(t, GetServiceValue(&svc), "http")
}
