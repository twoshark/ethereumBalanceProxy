package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twoshark/balanceproxy/src/common"
	"github.com/twoshark/balanceproxy/src/server"
)

var endpointsFlag *string
var dummyRPRCURL = "http://ethereumrpc1.com, https://localhost:8545"

// serverCmd represents the server command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Ethereum Balance Proxy Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *endpointsFlag == dummyRPRCURL {
			log.Error("Please provide 1 or more ethereum json rpc endpoints with the --upstreams flag. \neth_balance_proxy server --upstreams $ETH_RPC_ENDPOINT,$ETH_RPC_ENDPOINT")
			log.Exit(1)
		}
		port := viper.GetString("PORT")
		config := common.NewAppConfiguration(port, endpointsFlag)
		ready := make(chan bool)
		server.Start(config, ready)
		<-ready
		close(ready)
		log.Print("Server is ready to receive requests")
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(serverCmd)
	endpointsFlag = serverCmd.PersistentFlags().String("upstreams", dummyRPRCURL, "A comma separated list of backend endpoints to proxy")
}
