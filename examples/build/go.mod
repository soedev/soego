module build

go 1.15

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/soedev/soego v0.0.0-incompatible
)

replace github.com/soedev/soego => ../../../soego
