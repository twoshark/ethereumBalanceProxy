package ethereum

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/twoshark/balanceproxy/common"
)

var (
	badEndpoint    = "1'M A R3@L Y0U 4R3 3LL"
	nonEthEndpoint = "https://www.google.com"
	goodEndpoint   = "https://fittest-falling-smoke.discover.quiknode.pro/" // TODO: remove external dep
)

type ClientTestSuite struct {
	suite.Suite
}

func (suite *ClientTestSuite) SetupSuite() {
	common.CobraInit()
}

func (suite *ClientTestSuite) TearDownSuite() {
}

func (suite *ClientTestSuite) TestNewClient() {
	client := NewClient(goodEndpoint)
	assert.NotNil(suite.T(), client)
	assert.Equal(suite.T(), goodEndpoint, client.endpoint)
}

func (suite *ClientTestSuite) TestConnectSuccess() {
	client := NewClient(goodEndpoint)
	// New clients should not be healthy
	assert.Equal(suite.T(), false, client.Healthy())

	err := client.Dial()
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), client.EthClient())
}

func (suite *ClientTestSuite) TestConnectFailure() {
	client := NewClient(badEndpoint)
	assert.Equal(suite.T(), false, client.Healthy())
	err := client.Dial()
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), client.EthClient())
	assert.Equal(suite.T(), false, client.Healthy())
}

func TestBalanceProxyServerTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
