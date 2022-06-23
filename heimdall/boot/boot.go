package boot

import (
	"context"
	"flag"

	"workflow/heimdall/core"
	"workflow/heimdall/webserver"
	"workflow/workflow-utils/booting"
)

var (
	fileConfig = flag.String("config", "dev.yaml", "Config file")
)

// Boot straping server
func Boot() {
	flag.Parse()
	core.ReadConfig(*fileConfig)

	booting.BootstrapDaemons(context.Background(), webserver.WebServer)
}
