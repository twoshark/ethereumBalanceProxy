package ethereum

import (
	"context"
	"math/big"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IClient interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	CountHealthCheckFailure()
	CountHealthCheckSuccess()
	Dial() error
	EthClient() *ethclient.Client
	Healthy() bool
	HealthCheck() error
	SetHealth(bool)
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
}

type Client struct {
	endpoint      string
	ethClient     *ethclient.Client
	healthy       bool // healthy means the ethClient is connected and available to call
	successStreak int
	failureCount  int
	healthyLock   sync.Mutex
	clientLock    sync.Mutex
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// Dial connects to the rpc endpoint and sets c.ethClient.
func (c *Client) Dial() error {
	c.clientLock.Lock()
	defer c.clientLock.Unlock()
	// Clear existing clients for redials
	c.ethClient = nil
	ethClient, err := ethclient.Dial(c.endpoint)
	if err != nil {
		log.Error("failed to dial eth json rpc api: ", err)
		return err
	}

	c.ethClient = ethClient
	return nil
}

func (c *Client) Healthy() bool {
	c.healthyLock.Lock()
	defer c.healthyLock.Unlock()
	return c.healthy
}

func (c *Client) SetHealth(healthy bool) {
	c.healthyLock.Lock()
	defer c.healthyLock.Unlock()
	c.healthy = healthy
}

func (c *Client) EthClient() *ethclient.Client {
	c.clientLock.Lock()
	defer c.clientLock.Unlock()
	return c.ethClient
}
