package core

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"

	"workflow/heimdall/utils"
)

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
}

type ValkyrieConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// MainConfig contains all system configs
type MainConfig struct {
	DB       *DbConfig       `mapstructure:"database"`
	Server   *ServerConfig   `mapstructure:"server"`
	Valkyrie *ValkyrieConfig `mapstructure:"valkyrie"`
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
