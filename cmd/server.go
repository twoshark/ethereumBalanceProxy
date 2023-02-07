package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Ethereum Balance Proxy Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("server called")
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().String("endpoints", "http://localhost:8545", "A comma separated list of backend endpoints to proxy")

}
