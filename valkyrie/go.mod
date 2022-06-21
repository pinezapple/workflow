module github.com/vfluxus/valkyrie

go 1.14

replace github.com/vfluxus/workflow-utils => /home/pine/go/src/github.com/vfluxus/workflow-utils

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/aws/aws-sdk-go v1.35.14
	github.com/google/uuid v1.3.0
	github.com/labstack/echo/v4 v4.3.0
	github.com/spf13/viper v1.7.1
	github.com/swaggo/echo-swagger v1.1.0
	github.com/swaggo/swag v1.7.0
	github.com/uc-cdis/go-authutils v0.0.0-20201026165355-17b5f353bf4f
	github.com/vfluxus/workflow-utils v0.0.0-20201015040921-5f23d1a17a73
	go.temporal.io/sdk v1.15.0
	go.uber.org/zap v1.16.0
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gorm.io/driver/postgres v1.0.5
	gorm.io/gorm v1.20.7
)
