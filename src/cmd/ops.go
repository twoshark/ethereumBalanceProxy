package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var endpoint *string

// opsCmd represents the ops command.
var opsCmd = &cobra.Command{
	Use:   "ops",
	Short: "parent call for ops commands for the ethereum balance proxy.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("no-op")
	},
}

func init() {
	rootCmd.AddCommand(opsCmd)
	endpoint = opsCmd.PersistentFlags().String("endpoint", "", "Ethereum JSON RPC Endpoint")
}
