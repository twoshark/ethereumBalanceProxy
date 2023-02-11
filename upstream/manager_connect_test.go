package upstream

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	mock_ethereum "github.com/twoshark/balanceproxy/mocks"
	"github.com/twoshark/balanceproxy/upstream/ethereum"
)

type connectParams struct {
	Err   error
	Count int
}

type connectTestCase struct {
	Dial           connectParams
	HealthCheck    connectParams
	expectedOutput bool
}

type connectAllExpected struct {
	ConnectAllErr error
	ClientIndex   int
	ClientErr     error
}

type connectAllTestCase struct {
	expected      connectAllExpected
	clientExpects []connectTestCase
}

func (suite *ManagerTestSuite) TestConnect() {

	testCases := map[string]connectTestCase{
		"Successful Connect": {
			Dial:           connectParams{nil, 1},
			HealthCheck:    connectParams{nil, 1},
			expectedOutput: true,
		},
		"Dial Error": {
			Dial:           connectParams{suite.dummyErr, 1},
			HealthCheck:    connectParams{nil, 0},
			expectedOutput: false,
		},
		"Health Check Error": {
			Dial:           connectParams{nil, 1},
			HealthCheck:    connectParams{suite.dummyErr, 1},
			expectedOutput: false,
		},
	}

	for descriptor, testCase := range testCases {
		log.Print("Testing Connect Test Case: ", descriptor)
		suite.verifyConnectTestCase(testCase)
	}
}

func (suite *ManagerTestSuite) verifyConnectTestCase(testCase connectTestCase) {
	endpoints, count := suite.allEndpoints()
	mgr := NewManager(endpoints)
	assert.IsType(suite.T(), &Manager{}, mgr)
	assert.Equal(suite.T(), count, len(mgr.Clients))

	mockClient := mock_ethereum.NewMockIClient(suite.mockController)
	mockClient.EXPECT().Dial().Return(testCase.Dial.Err).AnyTimes()
	mockClient.EXPECT().HealthCheck().Return(testCase.HealthCheck.Err).AnyTimes()
	connected := mgr.Connect(mockClient)
	assert.Equal(suite.T(), testCase.expectedOutput, connected)
}

func (suite *ManagerTestSuite) TestConnectAllAndGetClient() {
	//the index is the client index in `Manager{}.Clients`
	testCases := map[string]connectAllTestCase{
		"All Clients Connect": {
			expected: connectAllExpected{
				ConnectAllErr: nil,
				ClientIndex:   0,
				ClientErr:     nil,
			},
			clientExpects: []connectTestCase{
				{
					Dial:        connectParams{Err: nil},
					HealthCheck: connectParams{Err: nil},
				},
				{
					Dial:        connectParams{Err: nil},
					HealthCheck: connectParams{Err: nil},
				},
				{
					Dial:        connectParams{Err: nil},
					HealthCheck: connectParams{Err: nil},
				},
			},
		},
		"All Clients Fail to Dial": {
			expected: connectAllExpected{
				ConnectAllErr: errNoUpstreamAvailable,
				ClientErr:     errNoHealthyUpstreamClient,
			},
			clientExpects: []connectTestCase{
				{
					Dial:        connectParams{Err: suite.dummyErr},
					HealthCheck: connectParams{Err: nil},
				},
				{
					Dial:        connectParams{Err: suite.dummyErr},
					HealthCheck: connectParams{Err: nil},
				},
				{
					Dial:        connectParams{Err: suite.dummyErr},
					HealthCheck: connectParams{Err: nil},
				},
			},
		},
		"All Clients Fail HealthCheck": {
			expected: connectAllExpected{
				ConnectAllErr: errNoUpstreamAvailable,
				ClientErr:     errNoHealthyUpstreamClient,
			},
			clientExpects: []connectTestCase{
				{
					Dial:        connectParams{Err: nil},
					HealthCheck: connectParams{Err: suite.dummyErr},
				},
				{
					Dial:        connectParams{Err: nil},
					HealthCheck: connectParams{Err: suite.dummyErr},
				},
				{
					Dial:        connectParams{Err: nil},
					HealthCheck: connectParams{Err: suite.dummyErr},
				},
			},
		},
		"Mixed Connection": {
			expected: connectAllExpected{
				ClientIndex:   2,
				ConnectAllErr: nil,
				ClientErr:     nil,
			},
			clientExpects: []connectTestCase{
				{
					Dial:        connectParams{Err: nil},
					HealthCheck: connectParams{Err: suite.dummyErr},
				},
				{
					Dial:        connectParams{Err: suite.dummyErr},
					HealthCheck: connectParams{Err: nil},
				},
				{
					Dial:        connectParams{Err: nil},
					HealthCheck: connectParams{Err: nil},
				},
			},
		},
	}

	for descriptor, testCase := range testCases {
		log.Print("Testing ConnectAll Test Case: ", descriptor)
		suite.verifyConnectAllTestCase(testCase)
	}
}

func (suite *ManagerTestSuite) verifyConnectAllTestCase(testCase connectAllTestCase) {
	endpoints := []string{"", "", ""}
	mgr := NewManager(endpoints)
	mgr.Clients = make([]ethereum.IClient, len(testCase.clientExpects))
	for i := 0; i < len(testCase.clientExpects); i++ {
		client := mock_ethereum.NewMockIClient(suite.mockController)
		client.EXPECT().Dial().Return(testCase.clientExpects[i].Dial.Err).AnyTimes()
		client.EXPECT().HealthCheck().Return(testCase.clientExpects[i].HealthCheck.Err).AnyTimes()
		//infer `Healthy()` from HeathCheck and Dial errors
		client.EXPECT().Healthy().Return(testCase.clientExpects[i].HealthCheck.Err == nil && testCase.clientExpects[i].Dial.Err == nil).AnyTimes()
		mgr.Clients[i] = client
	}
	err := mgr.ConnectAll()
	assert.Equal(suite.T(), testCase.expected.ConnectAllErr, err)

	client, err := mgr.GetClient()
	assert.Equal(suite.T(), testCase.expected.ClientErr, err)
	if err == nil {
		assert.Equal(suite.T(), mgr.Clients[testCase.expected.ClientIndex], client)
	} else {
		assert.Nil(suite.T(), client)
	}

}
