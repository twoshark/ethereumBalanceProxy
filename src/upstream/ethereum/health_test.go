package ethereum

import (
	"os"
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/twoshark/balanceproxy/src/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ClientHealthTestSuite struct {
	suite.Suite
}

func TestClientHealthTestSuite(t *testing.T) {
	suite.Run(t, new(ClientHealthTestSuite))
}

func (suite *ClientHealthTestSuite) SetupSuite() {
	common.CobraInit()
}

func (suite *ClientHealthTestSuite) TearDownSuite() {
}

func (suite *ClientHealthTestSuite) TestClient_processHealthCheckCalls() {
	failLimit := 3
	forgiveLimit := 10
	successThreshold := 6
	setEnv("HEALTH_FAILURE_THRESHOLD", strconv.Itoa(failLimit))
	setEnv("HEALTH_FAILURE_FORGIVENESS_THRESHOLD", strconv.Itoa(forgiveLimit))
	setEnv("HEALTH_SUCCESS_THRESHOLD", strconv.Itoa(successThreshold))
	setEnv("HEALTH_CHECK_PERIOD", "1")

	client := NewClient("")
	assert.Equal(suite.T(), 0, client.failureCount)
	assert.Equal(suite.T(), 0, client.successStreak)
	assert.Equal(suite.T(), false, client.Healthy())

	// force client to healthy state
	client.healthy = true
	assert.Equal(suite.T(), true, client.Healthy())

	// verify fail limit
	for i := 0; i < failLimit; i++ {
		client.processHealthCheckFailure()
		assert.Equal(suite.T(), i+1, client.failureCount)
		assert.Equal(suite.T(), 0, client.successStreak)
		assert.Equal(suite.T(), true, client.Healthy())
	}

	client.processHealthCheckFailure() // this should exceed the limit

	assert.Equal(suite.T(), 0, client.failureCount)
	assert.Equal(suite.T(), 0, client.successStreak)
	assert.Equal(suite.T(), false, client.Healthy())

	// verify successThreshold
	for i := 0; i < successThreshold; i++ {
		assert.Equal(suite.T(), 0, client.failureCount)
		assert.Equal(suite.T(), i, client.successStreak)
		assert.Equal(suite.T(), false, client.Healthy())
		client.processHealthCheckSuccess()
	}

	assert.Equal(suite.T(), 0, client.failureCount)
	assert.Equal(suite.T(), successThreshold, client.successStreak)
	assert.Equal(suite.T(), true, client.Healthy())

	// verify forgiveness threshold
	client.processHealthCheckFailure()
	assert.Equal(suite.T(), 1, client.failureCount)
	assert.Equal(suite.T(), 0, client.successStreak)
	assert.Equal(suite.T(), true, client.Healthy())

	for i := 0; i < forgiveLimit; i++ {
		client.processHealthCheckSuccess()
		assert.Equal(suite.T(), 1, client.failureCount)
		assert.Equal(suite.T(), i+1, client.successStreak)
		assert.Equal(suite.T(), true, client.Healthy())
	}

	client.processHealthCheckSuccess()
	assert.Equal(suite.T(), 0, client.failureCount)
	assert.Equal(suite.T(), 0, client.successStreak)
	assert.Equal(suite.T(), true, client.Healthy())
}

func setEnv(envVar string, val string) {
	if err := os.Setenv(envVar, val); err != nil {
		log.Panic(err)
	}
}
