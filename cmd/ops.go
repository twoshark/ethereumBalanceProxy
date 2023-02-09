package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// opsCmd represents the ops command.
var opsCmd = &cobra.Command{
	Use:   "ops",
	Short: "ops commands for the ethereum balance proxy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("ops called")
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(opsCmd)
}
