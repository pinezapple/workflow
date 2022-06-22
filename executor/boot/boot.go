package boot

import (
	"context"

	"workflow/executor/business"
	"workflow/executor/core"
	"workflow/executor/webserver"
	"workflow/workflow-utils/booting"

	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
)

func Boot() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	var ctx = context.Background()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	mainConfig := &core.MainConfig{
		ServiceName:          viper.GetString("SERVICE_NAME"),
		Environment:          viper.GetString("ENVIRONMENT"),
		LogLevel:             viper.GetString("LOG_LEVEL"),
		FailOverTime:         viper.GetInt("FAIL_OVER_TIME"),
		MaximumConcurrentJob: viper.GetInt("MAX_CONCURRENT_JOB"),
		HttpServerConfig: &core.HTTPServerConf{
			Port: viper.GetInt("PORT"),
		},
		K8SConfig: &core.K8SConf{
			K8SConfigFile:          viper.GetString("K8S_CONFIG_FILE"),
			K8SNameSpace:           viper.GetString("K8S_NAMESPACE"),
			OutputDirPrefix:        viper.GetString("OUTPUT_DIR_PREFIX"),
			DeleteJob:              viper.GetBool("K8S_DELETE_JOB"),
			JobTTLAfterFinished:    viper.GetInt("K8S_JOB_TTL_AFTER_FINISHED"),
			JobDeleteIntervalCheck: viper.GetInt("K8S_JOB_DELETE_INTERVAL_CHECK"),
			NodeLabelKey:           viper.GetString("K8S_NODE_LABEL_KEY"),
			NodeLabelValue:         viper.GetString("K8S_NODE_LABEL_VALUE"),
		},
		MinioEndpoint:  viper.GetString("MINIO_ENDPOINT"),
		FUSEMountpoint: viper.GetString("FUSE_MOUNT_POINT"),
	}

	core.SetMainConfig(mainConfig)
	core.InitCore(ctx)
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic("unable to create Temporal client")
	}
	business.SetExecutorTemporal(c)

	booting.BootstrapDaemons(ctx,
		webserver.WebServer,
		business.RunK8SDaemon,
		business.RunTemporalDaemon,
	)
}
