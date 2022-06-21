package core

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/vfluxus/heimdall/utils"
)

// KafkaConfig contains configs for Apache Kafka Client
type KafkaConfig struct {
	Version         string   `mapstructure:"version"`
	ProducerBrokers []string `mapstructure:"producerbrokers"`
	ConsumerBrokers []string `mapstructure:"consumerbrokers"`
	ProducerTopics  []string `mapstructure:"producertopics"`
	ConsumerTopics  []string `mapstructure:"consumertopics"`
	ConsumerGroup   string   `mapstructure:"consumergroup"`
}

// DbConfig database config
type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// ServerConfig start server
type ServerConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	ShutdownTimeout string `mapstructure:"shutdowntimeout"`
	Environment     string `mapstructure:"environment"`
	LogLevel        string `mapstructure:"loglevel"`

	JWKSUrl   string   `mapstructure:"jwks_url"`
	Audiences []string `mapstructure:"audiences"`
	Issuers   []string `mapstructure:"issuers"`
	Purpose   string   `mapstructure:"purpose"`
}

// SchedulerConfig scheduler config
type SchedulerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// TransformerConfig config
type TransformerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type ArboristConfig struct {
	Host string `mapstructure:"host"`
}

type MailServiceConfig struct {
	Address        string `mapstructure:"address"`
	AccountID      int    `mapstructure:"accountID"`
	NotifyTemplate int    `mapstructure:"notifyTemplateID"`
}

type ValkyrieConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// MainConfig contains all system configs
type MainConfig struct {
	DB          *DbConfig          `mapstructure:"database"`
	Server      *ServerConfig      `mapstructure:"server"`
	Scheduler   *SchedulerConfig   `mapstructure:"scheduler"`
	Kafka       *KafkaConfig       `mapstructure:"kafka"`
	Transformer *TransformerConfig `mapstructure:"transformer"`
	Valkyrie    *ValkyrieConfig    `mapstructure:"valkyrie"`
	Arborist    *ArboristConfig    `mapstructure:"arborist"`
	MailService *MailServiceConfig `mapstructure:"mailservice"`
}

var (
	m                   = new(MainConfig)
	initConfigFlag bool = true // true if not initiated, false if initiated
)

// ReadConfig read config file
func ReadConfig(configFile string) {
	// get config from file
	var (
		fileName string = filepath.Base(configFile)
		fileType string = filepath.Ext(configFile)[1:]
		filePath string = filepath.Dir(configFile)
	)

	if err := getConfigFromFile(fileName, fileType, filePath); err != nil {
		log.Fatalln("---- Can not get config from file, error: ", err)
		return
	}

	if err := getConfigFromOSEnv(); err != nil {
		log.Fatalln("---- Can not get config from env, error: ", err)
	}

	utils.PrintJSONFormat("Start with config", m)
	// set flag if ok
	initConfigFlag = false
}

// GetConfig return all system config
func GetConfig() *MainConfig {
	if initConfigFlag {
		log.Println("--- Main config not set ")
		return nil
	}
	return m
}

// -----------------------------------------------------------------------------
// -------------------------- GET FROM FILE ------------------------------------

// getConfigFromFile read file
func getConfigFromFile(fileName string, fileType string, filePath string) (err error) {
	v := viper.New()
	v.SetConfigName(fileName)
	v.SetConfigType(fileType)
	v.AddConfigPath(filePath)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(m); err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// ------------------------------- GET FROM OS ENV -----------------------------

// getConfigFromOSEnv replace config in file with config in os env. If os env
// not exist, keep the config in file
func getConfigFromOSEnv() (err error) {
	//host replace
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if !utils.IsStringEmpty(host) {
		m.Server.Host = host
	}
	if !utils.IsStringEmpty(port) {
		m.Server.Port = port
	}

	//scheduler replace
	schedulerHost := os.Getenv("SCHEDULER_HOST")
	schedulerPort := os.Getenv("SCHEDULER_PORT")
	if !utils.IsStringEmpty(schedulerHost) {
		m.Scheduler.Host = schedulerHost
	}
	if !utils.IsStringEmpty(schedulerPort) {
		m.Scheduler.Port = schedulerPort
	}

	return nil
}
