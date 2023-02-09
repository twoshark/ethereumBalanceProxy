package server

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/twoshark/balanceproxy/common"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type ServerTestSuite struct {
	suite.Suite
	endpoints []string
	port      string
}

func (suite *ServerTestSuite) SetupSuite() {
	port := viper.GetString("PORT")
	suite.port = port
	suite.endpoints = []string{
		"https://eth.getblock.io/b33bc13b-2d6b-4112-bd43-d93bb7cf842a/mainnet/",
		"https://fittest-falling-smoke.discover.quiknode.pro/",
	}
}

func (suite *ServerTestSuite) TearDownSuite() {
	// suite.stopProxy <- true
}

func (suite *ServerTestSuite) TestRootHandler() {
	resp, err := http.Get("http://localhost:" + suite.port)
	assert.Equal(suite.T(), nil, err)
	assert.NotEqual(suite.T(), nil, resp)
	if resp != nil {
		assert.Equal(suite.T(), 200, resp.StatusCode)
	}
}

func (suite *ServerTestSuite) TestLatestBalanceHandler() {
	resp, err := http.Get("http://localhost:" + suite.port + "/ethereum/balance/0x74630370197b4c4795bFEeF6645ee14F8cf8997D")
	assert.Equal(suite.T(), nil, err)
	assert.NotEqual(suite.T(), nil, resp)
	if resp != nil {
		assert.Equal(suite.T(), 200, resp.StatusCode)
	}
}

func (suite *ServerTestSuite) TestBalanceHandler() {
	testClient, err := ethclient.Dial(suite.endpoints[0])
	assert.NoError(suite.T(), err)
	block, err := testClient.BlockNumber(context.Background())
	assert.NoError(suite.T(), err)
	blockStr := strconv.FormatUint(block-10, 10)
	resp, err := http.Get("http://localhost:" +
		suite.port +
		"/ethereum/balance/0x74630370197b4c4795bFEeF6645ee14F8cf8997D" +
		"/block/" + blockStr)
	assert.Equal(suite.T(), nil, err)
	assert.NotEqual(suite.T(), nil, resp)
	if resp != nil {
		assert.Equal(suite.T(), 200, resp.StatusCode)
	}
}

func TestBalanceProxyServerTestSuite(t *testing.T) {
	common.CobraInit()
	port := viper.GetString("PORT")
	ready := make(chan bool)
	go Start(common.AppConfiguration{
		ListenPort: port,
		Endpoints: []string{
			"https://eth.getblock.io/b33bc13b-2d6b-4112-bd43-d93bb7cf842a/mainnet/",
			"https://fittest-falling-smoke.discover.quiknode.pro/",
		},
	}, ready)
	for {
		select {
		case <-ready:
			close(ready)
			suite.Run(t, new(ServerTestSuite))
			return
		default:
			log.Print("Waiting for Server Startup...")
			time.Sleep(10 * time.Second)
		}
	}
}
