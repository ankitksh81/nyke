package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

/*
	InitializeViper function initializes viper to read config.yml file
	and environment variables.
*/

func InitializeViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}
