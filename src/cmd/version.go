/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/twoshark/ethbalanceproxy/src/version"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(version.Get(*verbose))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
