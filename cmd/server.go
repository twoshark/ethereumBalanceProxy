package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twoshark/alluvial1-1/common"
	"github.com/twoshark/alluvial1-1/server"
)

var endpointsFlag *string

// serverCmd represents the server command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Ethereum Balance Proxy Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("server called")
		port := viper.GetString("PORT")
		config := common.NewAppConfiguration(port, endpointsFlag)
		server.Start(config)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().String("endpoints", "http://localhost:8545", "A comma separated list of backend endpoints to proxy")

}
