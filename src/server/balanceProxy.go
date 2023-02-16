package server

import (
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/twoshark/balanceproxy/src/upstream/ethereum"

	"github.com/twoshark/balanceproxy/src/common"
	"github.com/twoshark/balanceproxy/src/upstream"

	log "github.com/sirupsen/logrus"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
)

// BalanceProxy contains the handlers for the server.
type BalanceProxy struct {
	UpstreamManager *upstream.Manager
}

func NewBalanceProxy(config common.AppConfiguration) *BalanceProxy {
	bp := &BalanceProxy{
		UpstreamManager: upstream.NewManager(config.Endpoints),
	}
	return bp
}

func (bp *BalanceProxy) InitClients() {
	bp.UpstreamManager.LoadClients()
	var err error
	for {
		err = bp.UpstreamManager.ConnectAll()
		if err != nil {
			log.Error(err)
			time.Sleep(10 * time.Second)
			continue
		}
		return
	}
}

func (bp *BalanceProxy) RootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "server root")
}

// LiveHandler responds to a k8s liveness probe request
// it will always return an unconditional 200 as long as the server is running
func (bp *BalanceProxy) LiveHandler(c echo.Context) error {
	return c.String(http.StatusOK, "live")
}

// ReadyHandler responds to a k8s readiness probe request
// if there are any healthy upstreams, it will return a 200
// otherwise it will return a 503 and disable the ingress
func (bp *BalanceProxy) ReadyHandler(c echo.Context) error {
	counts := bp.UpstreamManager.ClientCounts()
	if counts.Healthy > 0 {
		return c.String(http.StatusOK, "ready")
	}
	return c.String(http.StatusServiceUnavailable, "no healthy upstreams")
}

func (bp *BalanceProxy) GetLatestBalance(c echo.Context) error {
	walletAddress := c.Param("account")
	address := ethcommon.HexToAddress(walletAddress)
	balance, err := bp.GetBalanceFromNode(c, address, nil)
	if err != nil {
		return err
	}

	return c.JSON(200, balance)
}

func (bp *BalanceProxy) GetBalance(c echo.Context) error {
	walletAddress := c.Param("account")
	address := ethcommon.HexToAddress(walletAddress)

	blockParam := c.Param("block")
	block, err := ParseBlockParam(blockParam)
	if err != nil {
		return err
	}

	balance, err := bp.GetBalanceFromNode(c, address, block)
	if err != nil {
		return err
	}

	return c.JSON(200, balance)
}

func (bp *BalanceProxy) GetBalanceFromNode(c echo.Context, address ethcommon.Address, block *big.Int) (map[string]interface{}, error) {
	logger := log.WithFields(log.Fields{
		"block":   block,
		"address": address,
	})
	var client ethereum.IClient
	var err error
	if block != nil && block.Uint64()+128 < bp.UpstreamManager.GetMaxBlock() {
		client, err = bp.UpstreamManager.GetArchiveClient()
		if err != nil {
			logger.Error("Error getting upstream archive client:", err)
			return nil, err
		}
	} else {
		client, err = bp.UpstreamManager.GetClient()
		if err != nil {
			logger.Error("Error getting upstream client:", err)
			return nil, err
		}
	}

	balance, err := client.BalanceAt(c.Request().Context(), address, block)
	if err != nil {
		logger.Error("Balance RPC Error:", err)
		return nil, err
	}

	return map[string]interface{}{"balance": balance.String()}, nil
}

func ParseBlockParam(blockParam string) (*big.Int, error) {
	var block *big.Int
	if blockParam == "" {
		// if the user does not specify, we return the latest block as per the example
		block = nil
	} else {
		parsedBlock, err := strconv.ParseInt(blockParam, 10, 64)
		if err != nil {
			return nil, err
		}
		block = big.NewInt(parsedBlock)
	}
	return block, nil
}
