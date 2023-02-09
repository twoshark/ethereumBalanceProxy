package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/twoshark/alluvial1-1/common"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "balanceProxy",
	Short: "Ethereum Balance Proxy",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// nolint: gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	common.CobraInit()
}
