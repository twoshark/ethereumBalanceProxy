package cmd

import (
	"fmt"

	"github.com/twoshark/balanceproxy/src/upstream/ethereum"

	"github.com/spf13/cobra"
)

// healthCheckCmd represents the healthCheck command
var healthCheckCmd = &cobra.Command{
	Use:   "healthCheck",
	Short: "run the proxy server's upstream health check against a single endpoint",
	Long:  ``,
	RunE:  HealthCheckCommand,
}

func init() {
	opsCmd.AddCommand(healthCheckCmd)
}

func HealthCheckCommand(*cobra.Command, []string) error {
	client := ethereum.NewClient(*endpoint)
	err := client.Dial()
	if err != nil {
		return err
	}
	err = client.HealthCheck()
	if err != nil {
		return err
	}
	fmt.Print("Healthy")
	return nil
}
