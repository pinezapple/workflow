package boot

import (
	"context"

	"workflow/workflow-utils/booting"
	"workflow/executor/business"
	"workflow/executor/core"
	"workflow/executor/webserver"

	"github.com/spf13/viper"
)

func Boot() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	var ctx = context.Background()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	core.CPULeft = viper.GetInt64("K8S_MAX_CPU")
	core.RAMLeft = viper.GetInt64("K8S_MAX_RAM")

	mainConfig := &core.MainConfig{
		ServiceName:            viper.GetString("SERVICE_NAME"),
		Environment:            viper.GetString("ENVIRONMENT"),
		LogLevel:               viper.GetString("LOG_LEVEL"),
		SelectTaskIntervalTime: viper.GetInt64("SELECT_JOB_INTERVAL_TIME"),
		APIRetryCount:          viper.GetInt("API_RETRY_COUNT"),
		FailOverTime:           viper.GetInt("FAIL_OVER_TIME"),
		MaximumConcurrentJob:   viper.GetInt("MAX_CONCURRENT_JOB"),
		HttpServerConfig: &core.HTTPServerConf{
			Port: viper.GetInt("PORT"),
		},
		KafkaConfig: &core.KafkaConf{
			ProducerBrokers: viper.GetStringSlice("PRODUCER_BROKERS"),
			ConsumerBrokers: viper.GetStringSlice("CONSUMER_BROKERS"),
			ConsumerTopics:  viper.GetStringSlice("CONSUMER_TOPICS"),
			ProducerTopics:  viper.GetStringSlice("PRODUCER_TOPICS"),
			ConsumerGroup:   viper.GetString("CONSUMER_GROUP"),
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
		//OutterBoundHTTP config
		SchedulerConfig: &core.OutterBoundHTTPOption{
			Addr: viper.GetString("SCHEDULER_ADDRESS"),
		},
		EddaConfig: &core.OutterBoundHTTPOption{
			Addr: viper.GetString("EDDA_ADDRESS"),
		},

		MinioEndpoint:  viper.GetString("MINIO_ENDPOINT"),
		FUSEMountpoint: viper.GetString("FUSE_MOUNT_POINT"),
	}

	core.SetMainConfig(mainConfig)
	core.InitCore(ctx)

	booting.BootstrapDaemons(ctx,
		webserver.WebServer,
		business.RunK8SDaemon,
		business.ReceiveMessageFromKafka,
		business.SelectTaskDaemons,
	)
}
