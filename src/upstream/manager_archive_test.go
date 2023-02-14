package upstream

import (
	_ "github.com/golang/mock/mockgen/model"
	"github.com/stretchr/testify/assert"
	mock_ethereum "github.com/twoshark/balanceproxy/mocks"
)

type ArchiveTestCaseClient struct {
	endpoint    string
	healthCheck error
	healthy     bool
	archive     bool
}

type ArchiveTestCaseOutput struct {
	returnedClientIndex int
	err                 error
}

type ArchiveTestCase struct {
	description string
	inputs      []ArchiveTestCaseClient
	outputs     ArchiveTestCaseOutput
}

func (suite *ManagerTestSuite) TestGetArchiveClient() {
	cases := []ArchiveTestCase{
		{
			description: "No Archive Endpoints Configured",
			inputs: []ArchiveTestCaseClient{
				{
					endpoint:    "https://www.dummyFullNode.com",
					healthCheck: nil,
					healthy:     true,
					archive:     false,
				},
			},
			outputs: ArchiveTestCaseOutput{
				err: errNoHealthyArchiveUpstreamClient,
			},
		},
		{
			description: "No Healthy Archive Endpoints Available",
			inputs: []ArchiveTestCaseClient{
				{
					endpoint:    "https://www.dummyFullNode.com",
					healthCheck: nil,
					healthy:     true,
					archive:     false,
				},
				{
					endpoint:    "https://www.dummyUnhealthyArchiveNode.com",
					healthCheck: suite.dummyErr,
					archive:     true,
				},
			},
			outputs: ArchiveTestCaseOutput{
				err: errNoHealthyArchiveUpstreamClient,
			},
		},
		{
			description: "Healthy Archive Endpoints Available",
			inputs: []ArchiveTestCaseClient{
				{
					endpoint:    "https://www.dummyFullNode.com",
					healthCheck: nil,
					healthy:     true,
					archive:     false,
				},
				{
					endpoint:    "https://www.dummyUnhealthyArchiveNode.com",
					healthCheck: nil,
					healthy:     true,
					archive:     true,
				},
			},
			outputs: ArchiveTestCaseOutput{
				returnedClientIndex: 1,
				err:                 nil,
			},
		},
	}

	for _, testCase := range cases {
		var endpoints []string
		for i := 0; i < len(testCase.inputs); i++ {
			endpoints = append(endpoints, testCase.inputs[i].endpoint)
		}
		mgr := NewManager(endpoints)
		for i := 0; i < len(testCase.inputs); i++ {
			client := mock_ethereum.NewMockIClient(suite.mockController)
			client.EXPECT().Dial().Return(nil).AnyTimes()
			client.EXPECT().CheckIfArchive().AnyTimes()
			client.EXPECT().Healthy().Return(testCase.inputs[i].healthy).AnyTimes()
			client.EXPECT().HealthCheck().Return(testCase.inputs[i].healthCheck).AnyTimes()
			client.EXPECT().IsArchive().Return(testCase.inputs[i].archive).AnyTimes()
			mgr.Clients[i] = client
		}
		archiveClient, err := mgr.GetArchiveClient()
		if testCase.outputs.err != nil {
			if assert.Error(suite.T(), err) {
				assert.Equal(suite.T(), errNoHealthyArchiveUpstreamClient, err)
			}
			continue
		}
		assert.Equal(suite.T(), mgr.Clients[testCase.outputs.returnedClientIndex], archiveClient)
		assert.NoError(suite.T(), err)
	}
}
