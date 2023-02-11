package upstream

import (
	"github.com/stretchr/testify/assert"
	mock_ethereum "github.com/twoshark/balanceproxy/mocks"
)

type connectParams struct {
	Err   error
	Count int
}

type connectTestCase struct {
	dial           connectParams
	healthCheck    connectParams
	expectedOutput bool
}

func (suite *ManagerTestSuite) TestConnect() {

	testCases := []connectTestCase{
		{
			dial:           connectParams{nil, 1},
			healthCheck:    connectParams{nil, 1},
			expectedOutput: true,
		},
		{
			dial:           connectParams{suite.dummyErr, 1},
			healthCheck:    connectParams{nil, 0},
			expectedOutput: false,
		},
		{
			dial:           connectParams{nil, 1},
			healthCheck:    connectParams{suite.dummyErr, 1},
			expectedOutput: false,
		},
	}

	for _, testCase := range testCases {
		suite.verifyConnectTestCase(testCase)
	}
}

func (suite *ManagerTestSuite) verifyConnectTestCase(testCase connectTestCase) {
	endpoints, count := suite.allEndpoints()
	mgr := NewManager(endpoints)
	assert.IsType(suite.T(), &Manager{}, mgr)
	assert.Equal(suite.T(), count, len(mgr.Clients))

	mockClient := mock_ethereum.NewMockIClient(suite.mockController)
	mockClient.EXPECT().Dial().Return(testCase.dial.Err).AnyTimes()
	mockClient.EXPECT().HealthCheck().Return(testCase.healthCheck.Err).AnyTimes()
	connected := mgr.Connect(mockClient)
	assert.Equal(suite.T(), testCase.expectedOutput, connected)
}
