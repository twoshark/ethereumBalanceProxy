package ethereum

import "errors"

var (
	errSyncStatusUnavailable      = errors.New("failed to get sync state of upstream")
	errUpstreamSyncing            = errors.New("upstream is still syncing")
	errBlockHeightZero            = errors.New("node has not synced any blocks and is not syncing. BlockHeight: 0")
	errUpstreamRewind             = errors.New("upstream is rewinding or has experienced a chain re-org. please try again later")
	errBlockHeightIncreaseTooSlow = errors.New("upstream block height is not increasing sufficiently relative to network norms")
)
