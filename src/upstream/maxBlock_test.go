package upstream

import (
	"github.com/stretchr/testify/assert"
	mock_ethereum "github.com/twoshark/balanceproxy/mocks"
)

type MaxBlockTestCase struct {
	clientHeights []uint64
	expectedMax   uint64
}

func (suite *ManagerTestSuite) TestCheckClientMaxBlocks() {

	mgr := NewManager([]string{"", "", ""})
	// Clients { frozen block height, climbing normally, unstable}
	cases := []MaxBlockTestCase{
		{
			clientHeights: []uint64{1, 1, 1},
			expectedMax:   1,
		},
		{
			clientHeights: []uint64{1, 2, 2},
			expectedMax:   2,
		},
		{
			clientHeights: []uint64{1, 3, 1},
			expectedMax:   3,
		},
		{
			clientHeights: []uint64{1, 4, 4},
			expectedMax:   4,
		},
	}

	assert.Zero(suite.T(), mgr.GetMaxBlock())
	for _, testCase := range cases {
		for i := 0; i < len(testCase.clientHeights); i++ {
			client := mock_ethereum.NewMockIClient(suite.mockController)
			client.EXPECT().GetMaxBlock().Return(testCase.clientHeights[i])
			mgr.Clients[i] = client
		}
		mgr.CheckClientMaxBlocks()
		assert.Equal(suite.T(), testCase.expectedMax, mgr.GetMaxBlock())
	}
}
