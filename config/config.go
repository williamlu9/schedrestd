package config

import (
	"schedrestd/common"
	"fmt"
	"go.uber.org/fx"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	LogLevel   string `yaml:"log_level"`
	Timeout    string `yaml:"timeout"`
	WebUrlPath string `yaml:"web_url_path"`
	Ssl        string `yaml:"ssl"`
	HttpPort   string `yaml:"http_port"`
	HttpsPort  string `yaml:"https_port"`
	Cert       string `yaml:"cert"`
	Key        string `yaml:"key"`
	LogDir     string `yaml:"logdir"`
	WorkDir    string `yaml:"workdir"`
}

// NewConfig ...
func NewConfig() *Config {
	confDir := GetCBENVDir()
	conf := Config{}

	// Set the viper config file path
	viper.SetConfigName(common.ConfigFileName)
	viper.AddConfigPath(confDir)
	viper.SetConfigType(common.ConfigType)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf(err.Error())
		return &Config{}
	}

	md := mapstructure.Metadata{}
	err = viper.Unmarshal(&conf, func(config *mapstructure.DecoderConfig) {
		config.TagName = common.ConfigType
		config.Metadata = &md
	})
	if err != nil {
		fmt.Printf(err.Error())
		return &Config{}
	}

	if conf.LogDir == "" {
		conf.LogDir = GetLogDir()
	}
	if conf.WorkDir == "" {
		conf.WorkDir = "/var/run"
	}

	if conf.WebUrlPath == "" {
		conf.WebUrlPath = "/"
	}

	return &conf
}

func GetCBENVDir() string {
	return "/etc/schedrestd"
}

func GetCBBinDir() string {
	return "/usr/bin"
}

func GetLogDir() string {
	return "/var/log"
}

// Module Export dependency injection
var Module = fx.Options(
	fx.Provide(NewConfig),
)
