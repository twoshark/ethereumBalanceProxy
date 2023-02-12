package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/twoshark/balanceproxy/src/metrics"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	defer measureLatency(time.Now(), "latency_eth_get_balance")
	return c.EthClient().BalanceAt(ctx, account, blockNumber)
}

func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	defer measureLatency(time.Now(), "latency_eth_get_block_number")
	return c.EthClient().BlockNumber(ctx)
}

func (c *Client) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	defer measureLatency(time.Now(), "latency_eth_syncing")
	return c.EthClient().SyncProgress(ctx)
}

func measureLatency(t time.Time, metric string) {
	duration := time.Since(t)
	switch metric {
	case "latency_eth_get_balance":
		metrics.Metrics().EthGetBalanceLatency.With(nil).Observe(float64(duration.Milliseconds()))
	case "latency_eth_get_block_number":
		metrics.Metrics().EthGetBlockNumberLatency.With(nil).Observe(float64(duration.Milliseconds()))
	case "latency_eth_syncing":
		metrics.Metrics().EthSyncingLatency.With(nil).Observe(float64(duration.Milliseconds()))
	}
}
