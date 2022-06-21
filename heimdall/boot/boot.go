package boot

import (
	"context"
	"flag"

	"github.com/vfluxus/heimdall/core"
	"github.com/vfluxus/heimdall/infra"
	"github.com/vfluxus/heimdall/webserver"
	"github.com/vfluxus/workflow-utils/booting"
)

var (
	fileConfig = flag.String("config", "dev.yaml", "Config file")
)

// Boot straping server
func Boot() {
	flag.Parse()
	core.ReadConfig(*fileConfig)

	// init the mailservice
	mailConfig := core.GetConfig().MailService
	mailSrv, err := infra.NewMailService(mailConfig.Address, mailConfig.AccountID, mailConfig.NotifyTemplate)
	if err != nil {
		panic(err)
	}
	infra.ReplaceGlobalMailSrv(mailSrv)

	booting.BootstrapDaemons(context.Background(), webserver.WebServer)
}
