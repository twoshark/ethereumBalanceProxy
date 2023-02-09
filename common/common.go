package common

import (
	"github.com/spf13/viper"
)

func CobraInit() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("SHUTDOWN_TIMEOUT", 2)
	viper.AutomaticEnv() // read in environment variables that match
}
