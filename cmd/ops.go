package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// opsCmd represents the ops command
var opsCmd = &cobra.Command{
	Use:   "ops",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("ops called")
	},
}

func init() {
	rootCmd.AddCommand(opsCmd)
}
