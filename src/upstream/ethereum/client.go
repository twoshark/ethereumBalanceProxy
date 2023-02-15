package ethereum

import (
	"context"
	"math/big"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IClient interface {
	IsArchive() bool
	CheckIfArchive()
	BalanceAt(ctx context.Context, account ethCommon.Address, blockNumber *big.Int) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	Dial() error
	EvaluatedHealthCheck()
	EthClient() *ethclient.Client
	GetMaxBlock() uint64
	Healthy() bool
	HealthCheck() error
	SetHealth(bool)
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
}

type Client struct {
	archive       bool
	archiveLock   sync.Mutex
	blockLock     sync.Mutex
	clientLock    sync.Mutex
	endpoint      string
	ethClient     *ethclient.Client
	failureCount  int
	healthy       bool // healthy means the ethClient is connected and available to call
	healthyLock   sync.Mutex
	maxBlock      uint64
	successStreak int
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

// IsArchive returns the current archive state of this endpoint
func (c *Client) IsArchive() bool {
	return c.archive
}

// CheckIfArchive tries to pull the balance of a known wallet at block 15,000,000.
// Only archive instances will be able to return a block that old.
func (c *Client) CheckIfArchive() {
	address := ethCommon.HexToAddress("0x74630370197b4c4795bFEeF6645ee14F8cf8997D")
	_, err := c.BalanceAt(context.Background(), address, big.NewInt(15000000))
	c.archiveLock.Lock()
	defer c.archiveLock.Unlock()
	if err != nil {
		log.Error(err)
		c.archive = false
		return
	}
	c.archive = true
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

func (c *Client) ProposeNewMaxBlock(newBlock uint64) {
	c.blockLock.Lock()
	defer c.blockLock.Unlock()
	if newBlock > c.maxBlock {
		c.maxBlock = newBlock
	}
}

func (c *Client) GetMaxBlock() uint64 {
	c.blockLock.Lock()
	defer c.blockLock.Unlock()
	return c.maxBlock
}
