package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// HealthCheck verifies that the upstream node is not running a sync process
// and that its block height is increasing
func (c *Client) HealthCheck() error {
	c.SetHealth(false)
	// Verify that no sync process is running
	sync, err := c.SyncProgress(context.Background())
	if err != nil {
		log.Error(err)
		return errSyncStatusUnavailable
	}

	if sync != nil {
		var bytes []byte
		bytes, err = json.Marshal(sync)
		if err != nil {
			return err
		}
		log.Error(fmt.Sprintf("upstream sync state: %s", bytes))
		return errUpstreamSyncing
	}

	// Verify that the json rpc api responds to a request for block height
	lastBlock, err := c.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	// Verify that block height is > 0
	if lastBlock == 0 {
		return errBlockHeightZero
	}

	err = c.isBlockHeightIncreasing(lastBlock)
	if err != nil {
		return err
	}

	c.SetHealth(true)
	return nil
}

func (c *Client) isBlockHeightIncreasing(startBlock uint64) error {
	// Verify that block height is climbing at a reasonable pace
	blockHeightIncreased := 0
	increaseObservationWindow := 3

	var block uint64
	lastBlock := startBlock
	var err error
	// period := viper.GetInt64("HEALTH_BLOCK_HEIGHT_CHECK_PERIOD_MS")

	for i := 0; i < increaseObservationWindow; i++ {
		block, err = c.EthClient().BlockNumber(context.Background())
		if err != nil {
			return err
		}
		log.Print(block, ":", lastBlock)
		if block > lastBlock {
			return nil
		} else if block < lastBlock {
			return errUpstreamRewind
		}
		lastBlock = block
		time.Sleep(12 * time.Second) // Eth Avg Block Time is 12.06s
	}

	if blockHeightIncreased <= 1 {
		return errBlockHeightIncreaseTooSlow
	}
	return nil
}
