package config

import (
	"strings"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetEnvPrefix("BANDA_LUMAKSA")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/banda-lumaksa/")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.ReadInConfig()
}
