package upstream

import (
	"testing"

	mock_ethereum "github.com/twoshark/balanceproxy/mocks"

	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	suite.Suite
	endpoints      []string
	mockController *gomock.Controller
}

func (suite *ManagerTestSuite) SetupSuite() {
	suite.endpoints = []string{"", "", ""}
	suite.mockController = gomock.NewController(suite.T())
	_ = mock_ethereum.NewMockIClient(suite.mockController)
	// client.EXPECT().TryReconnect(gomock.Any()).Return(nil, nil).AnyTimes()
	// client.EXPECT().SyncProgress(gomock.Any()).Return(nil, nil).AnyTimes()
	// suite.U = upstreamClient.NewEthereumClient("https://eth.getblock.io/b33bc13b-2d6b-4112-bd43-d93bb7cf842a/mainnet/", client)
	// err := suite.U.TryReconnect()
	// assert.Equal(suite.T(), nil, err)
	// assert.NotNil(suite.T(), client)
	// assert.Equal(suite.T(), suite.U.Connected, true)
	// assert.Equal(suite.T(), suite.U.Healthy, true)
}

func (suite *ManagerTestSuite) TearDown() {
	suite.mockController.Finish()
}

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
