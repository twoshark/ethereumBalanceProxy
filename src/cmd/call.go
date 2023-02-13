/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/twoshark/balanceproxy/src/upstream/ethereum"
)

var account, block, method *string

// callCmd represents the call command
var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Execute eth rpc calls using the balance proxy's client",
	Long:  `Execute an ethereum json rpc call against an endpoint using balance proxy's client directly, bypassing the Server and Upstream Manager'`,
	RunE:  CallCommand,
}

func init() {
	opsCmd.AddCommand(callCmd)

	account = callCmd.PersistentFlags().String("account", "", "Ethereum Wallet Address")
	block = callCmd.PersistentFlags().String("block", "latest", "Ethereum Block Number or `latest`")
	method = callCmd.PersistentFlags().String("method", "", "Ethereum JSON RPC Method")
}

func CallCommand(*cobra.Command, []string) error {
	client := ethereum.NewClient(*endpoint)
	err := client.Dial()
	if err != nil {
		return err
	}

	switch *method {
	case "eth_syncing":
		sync, err := client.SyncProgress(context.Background())
		if err != nil {
			return err
		}
		if sync != nil {
			fmt.Printf("%#v", *sync)
		} else {
			fmt.Print("No Sync Process in Progress")
		}
	case "eth_getBalance":
		address := ethcommon.HexToAddress(*account)
		var blockInt *big.Int
		if *block == "latest" {
			blockInt = nil
		} else {
			blockInt64, err := strconv.ParseInt(*block, 10, 64)
			if err != nil {
				return err
			}
			blockInt = big.NewInt(blockInt64)
		}

		balance, err := client.BalanceAt(context.Background(), address, blockInt)
		if err != nil {
			return err
		}
		if balance != nil {
			fmt.Print(balance.String())
		} else {
			fmt.Print("Nil Balance Response")
		}
	case "eth_getBlockNumber":
		block, err := client.BlockNumber(context.Background())
		if err != nil {
			return err
		}
		fmt.Print(block)
	default:
		return errors.New("provided method unsupported")
	}
	return nil
}
