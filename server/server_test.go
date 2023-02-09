package server

import (
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/twoshark/alluvial1-1/common"
)

type BalanceProxyServerTestSuite struct {
	suite.Suite
	endpoint string
}

func (suite *BalanceProxyServerTestSuite) SetupSuite() {
	common.CobraInit()
	port := viper.GetString("PORT")
	suite.endpoint = "http://localhost:" + port
	go Start()
}

func (suite *BalanceProxyServerTestSuite) TearDownSuite() {
	// suite.stopProxy <- true
}

func (suite *BalanceProxyServerTestSuite) TestRootHandler() {
	resp, err := http.Get(suite.endpoint)
	assert.Equal(suite.T(), nil, err)
	assert.NotEqual(suite.T(), nil, resp)
	if resp != nil {
		assert.Equal(suite.T(), 200, resp.StatusCode)
	}
}

func TestBalanceProxyServerTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceProxyServerTestSuite))
}
