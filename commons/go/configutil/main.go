package configutil

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
	"strings"
)

func LoadConfig(configPath string, logger log.Logger, config interface{}) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/")
	viper.AddConfigPath(".")
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		level.Error(logger).Log("msg", "error reading config file", "err", err)
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		level.Error(logger).Log("msg", "error unmarshalling into struct", "err", err)
		return err
	}

	return nil
}
