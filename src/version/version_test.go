package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VersionTestSuite struct {
	suite.Suite
}

func TestVersionTestSuite(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}

func (suite *VersionTestSuite) TestVersion() {
	Set("a", "b", "c")
	assert.Equal(suite.T(), "a", Get(false))
	assert.Equal(suite.T(), "a-b (c)", Get(true))
}
