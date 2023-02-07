package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.EthClient().BalanceAt(ctx, account, blockNumber)
}

func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	return c.EthClient().BlockNumber(ctx)
}

func (c *Client) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	return c.EthClient().SyncProgress(ctx)
}
