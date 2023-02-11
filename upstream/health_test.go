package upstream

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type UpstreamHealthTestSuite struct {
	suite.Suite
}

func TestUpstreamHealthTestSuite(t *testing.T) {
	suite.Run(t, new(UpstreamHealthTestSuite))
}

func (suite UpstreamHealthTestSuite) SetupSuite() {}
func (suite UpstreamHealthTestSuite) Teardown()   {}

func (suite UpstreamHealthTestSuite) TestStartHealthCheck() {

}