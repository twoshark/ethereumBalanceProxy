package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
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

	if err = c.isBlockHeightIncreasing(); err != nil {
		return err
	}

	c.SetHealth(true)
	return nil
}

func (c *Client) isBlockHeightIncreasing() error {
	// Verify that the json rpc api responds to a request for block height
	lastBlock, err := c.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	// Verify that block height is > 0
	if lastBlock == 0 {
		return errBlockHeightZero
	}
	// Verify that block height is climbing at a reasonable pace
	increaseObservationWindow := 3
	var block uint64
	period := viper.GetInt("HEALTH_BLOCK_HEIGHT_CHECK_PERIOD_MS")
	ticker := time.NewTicker(time.Duration(period) * time.Millisecond)
	count := 0
	for ; true; <-ticker.C {
		block, err = c.EthClient().BlockNumber(context.Background())
		if err != nil {
			return err
		}
		if block > lastBlock {
			return nil
		} else if block < lastBlock {
			return errUpstreamRewind
		}
		lastBlock = block

		count++
		if count >= increaseObservationWindow {
			return errBlockHeightIncreaseTooSlow
		}
	}
}
