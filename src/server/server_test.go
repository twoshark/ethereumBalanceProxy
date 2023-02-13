package server

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/twoshark/balanceproxy/src/common"
)

func RunSuiteWithServer(t *testing.T, testSuite suite.TestingSuite) {
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
			suite.Run(t, testSuite)
			return
		default:
			log.Print("Waiting for Server Startup...")
			time.Sleep(10 * time.Second)
		}
	}
}

type ServerTestSuite struct {
	suite.Suite
	endpoints []string
	port      string
}

func TestBalanceProxyServerTestSuite(t *testing.T) {
	common.CobraInit()
	RunSuiteWithServer(t, new(ServerTestSuite))
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

func (suite *ServerTestSuite) TestLiveHandler() {
	resp, err := http.Get("http://localhost:" + suite.port + "/live")
	assert.Equal(suite.T(), nil, err)
	assert.NotEqual(suite.T(), nil, resp)
	if resp != nil {
		assert.Equal(suite.T(), 200, resp.StatusCode)
	}
}

func (suite *ServerTestSuite) TestReadyHandler() {
	resp, err := http.Get("http://localhost:" + suite.port + "/ready")
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

func (suite *ServerTestSuite) TestMetricHandler() {
	resp, err := http.Get("http://localhost:" + suite.port + "/metrics")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), len(responseData), 0)
	response := string(responseData)

	// Custom Metrics
	assert.True(suite.T(), strings.Contains(response, "startup_time_ms"))
	assert.True(suite.T(), strings.Contains(response, "upstreams_configured"))
	assert.True(suite.T(), strings.Contains(response, "upstreams_healthy"))
	assert.True(suite.T(), strings.Contains(response, "latency_eth_syncing"))
	assert.True(suite.T(), strings.Contains(response, "latency_eth_get_block_number"))
	assert.True(suite.T(), strings.Contains(response, "latency_eth_get_balance"))

	// Go Metrics
	assert.True(suite.T(), strings.Contains(response, "go_"))

	// Echo Metrics
	assert.True(suite.T(), strings.Contains(response, "echo"))
}
