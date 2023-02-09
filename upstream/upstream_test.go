package upstream

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/twoshark/alluvial1-1/common"
	"testing"
)

type UpstreamTestSuite struct {
	suite.Suite
	endpoint string
}

func (suite *UpstreamTestSuite) SetupSuite() {
	common.CobraInit()
}

func (suite *UpstreamTestSuite) TearDownSuite() {
}

func (suite *UpstreamTestSuite) TestNewUpstream() {
	endpoint := "https://www.totallyreal.gov"
	u := NewUpstream(endpoint)
	assert.NotNil(suite.T(), u)
	assert.Equal(suite.T(), endpoint, u.Endpoint)
	assert.Equal(suite.T(), true, u.Healthy)
}

func TestBalanceProxyServerTestSuite(t *testing.T) {
	suite.Run(t, new(UpstreamTestSuite))
}
