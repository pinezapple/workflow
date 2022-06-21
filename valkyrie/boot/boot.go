package boot

import (
	"context"
	"flag"

	"github.com/spf13/viper"

	"workflow/valkyrie/business"
	"workflow/valkyrie/core"
	"workflow/valkyrie/webserver"
	"workflow/workflow-utils/booting"
)

var (
	first = flag.Bool("first", false, "if we need to start recovering sequence")
)

func Boot() {
	flag.Parse()
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	var ctx = context.Background()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	mainConfig := &core.MainConfig{
		ServiceName: viper.GetString("SERVICE_NAME"),
		Environment: viper.GetString("ENVIRONMENT"),
		LogLevel:    viper.GetString("LOG_LEVEL"),

		HttpServerConfig: &core.HTTPServerConf{
			Port:   viper.GetInt("PORT"),
			JwkURL: viper.GetString("JWL_URL"),

			ExpectJWTConfig: &core.ExpectJWTConf{
				Audiences: viper.GetStringSlice("AUDIENCES"),
				Issuers:   viper.GetStringSlice("ISSUERS"),
				Purpose:   viper.GetString("PURPOSE"),
			},
		},
		MinioConfig: &core.MinioConf{
			AccessKey:      viper.GetString("ACCESS_KEY"),
			SecretKey:      viper.GetString("SECRET_KEY"),
			Region:         viper.GetString("REGION"),
			Endpoint:       viper.GetString("ENDPOINT"),
			BucketMaxSize:  1073741824 * viper.GetInt64("BUCKET_MAX_SIZE"),
			BucketMaxCount: viper.GetInt32("BUCKET_MAX_COUNT"),
		},
		DBConfig: &core.DBConf{
			DBName:   viper.GetString("DB_NAME"),
			Username: viper.GetString("DB_USERNAME"),
			Password: viper.GetString("DB_PASSWORD"),
			Hostname: viper.GetString("DB_HOSTNAME"),
			Port:     viper.GetString("DB_PORT"),
		},
		HeimdallConfig: &core.OutterBoundHTTPOption{
			Addr: viper.GetString("HEIMDALL_ADDR"),
		},

		ExpiredUpdateInterval:       viper.GetInt("EXPIRED_INTERVAL"),
		UpdateStorageChangeInterval: viper.GetInt("DELETE_INTERVAL"),
		NoDelete:                    viper.GetBool("NO_DELETE"),

		InputDirPrefix:      viper.GetString("INPUT"),
		UploadTempDirPrefix: viper.GetString("TEMP_UPLOAD"),

		UploadPenalty:    viper.GetInt("UPLOAD_PENALTY"),
		NormalFileTTL:    viper.GetInt("NORMAL_FILE_TTL"),
		ImportantFileTTL: viper.GetInt("IMPORTANT_FILE_TTL"),

		HardDiskOnly:    viper.GetBool("HARD_DISK_ONLY"),
		APIRetryCount:   viper.GetInt("API_RETRY_COUNT"),
		MinioEndpoint:   viper.GetString("MINIO_ENDPOINT"),
		MinioAuthenFile: viper.GetString("MINIO_AUTHEN_FILE"),
		FUSEMountpoint:  viper.GetString("FUSE_MOUNT_POINT"),
	}

	core.SetMainConfig(mainConfig)
	// start to init core
	core.InitCore(ctx)

	// TODO: if first
	if *first {
		business.RestartDaemon(ctx)
	}

	if mainConfig.HardDiskOnly {
		if mainConfig.NoDelete {
			booting.BootstrapDaemons(context.Background(),
				webserver.WebServer,
			)

		} else {
			booting.BootstrapDaemons(context.Background(),
				business.DailyCheckOnExpiredObjOnHardDisk,
				webserver.WebServer,
			)
		}
	} else {
		if mainConfig.NoDelete {
			booting.BootstrapDaemons(context.Background(),
				webserver.WebServer,
			)
		} else {
			booting.BootstrapDaemons(context.Background(),
				//business.DailyCheckOnExpiredObj,
				//business.CheckOnExpiredUploadObj,
				webserver.WebServer,
			)
		}

	}
}
