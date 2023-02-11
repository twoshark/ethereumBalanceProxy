package upstream

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/twoshark/balanceproxy/upstream/ethereum"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	suite.Suite
	endpointSets   map[string][]string
	mockController *gomock.Controller
	dummyErr       error
}

func (suite *ManagerTestSuite) allEndpoints() ([]string, int) {
	endpoints := make([]string, 0)
	for _, endpointSet := range suite.endpointSets {
		endpoints = append(endpoints, endpointSet...)
	}
	return endpoints, len(endpoints)
}

func (suite *ManagerTestSuite) SetupSuite() {
	suite.dummyErr = errors.New("uh oh")
	suite.endpointSets = make(map[string][]string)
	suite.endpointSets =
		map[string][]string{
			"invalid": {
				"",
				" f ^*&* jnkdasik // sad78asgyuhab",
			},
			"nonEthRpc": {
				"https://google.com",
				"https://www.hats.com",
			},
			"validExternal": {
				"https://eth.getblock.io/b33bc13b-2d6b-4112-bd43-d93bb7cf842a/mainnet/",
				"https://mainnet.infura.io/v3/e2edc69a0cef4ff28466331d6d972560",
				"https://fittest-falling-smoke.discover.quiknode.pro/",
			},
		}

	suite.mockController = gomock.NewController(suite.T())
}

func (suite *ManagerTestSuite) TearDown() {
	suite.mockController.Finish()
}

func (suite *ManagerTestSuite) TestNewManager() {
	mgr := NewManager([]string{"", ""})
	assert.IsType(suite.T(), &Manager{}, mgr)
	assert.Equal(suite.T(), 2, len(mgr.Clients))
}

func (suite *ManagerTestSuite) TestLoadClients() {
	endpoints, count := suite.allEndpoints()
	mgr := NewManager(endpoints)
	assert.IsType(suite.T(), &Manager{}, mgr)
	assert.Equal(suite.T(), count, len(mgr.Clients))
	mgr.LoadClients(endpoints)
	for _, client := range mgr.Clients {
		assert.NotNil(suite.T(), client)
		assert.IsType(suite.T(), &ethereum.Client{}, client)
	}
}

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
