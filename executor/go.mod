module github.com/vfluxus/workflow/executor

go 1.14

replace (
	github.com/vfluxus/heimdall => /home/pine/go/src/github.com/vfluxus/heimdall/app
	github.com/vfluxus/workflow-utils => /home/pine/go/src/github.com/vfluxus/workflow-utils
)

require (
	github.com/Shopify/sarama v1.28.0
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/spf13/viper v1.7.1
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/vfluxus/dvergr v0.3.0
	github.com/vfluxus/heimdall v0.0.0-00010101000000-000000000000
	github.com/vfluxus/workflow-utils v0.0.0-20201130034434-cbda7c2a1b87
	go.temporal.io/sdk v1.15.0
	go.uber.org/zap v1.16.0
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/utils v0.0.0-20200821003339-5e75c0163111 // indirect
)
