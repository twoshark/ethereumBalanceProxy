package common

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func CobraInit() {
	log.SetLevel(log.InfoLevel)
	viper.SetDefault("HEALTH_FAILURE_THRESHOLD", 2)
	viper.SetDefault("HEALTH_FAILURE_FORGIVENESS_THRESHOLD", 10)
	viper.SetDefault("HEALTH_CHECK_PERIOD", 60)
	viper.SetDefault("HEALTH_SUCCESS_THRESHOLD", 4)
	viper.SetDefault("HEALTH_BLOCK_HEIGHT_CHECK_PERIOD_MS", 12100)
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("SHUTDOWN_TIMEOUT", 2)
	viper.AutomaticEnv() // read in environment variables that match
}
