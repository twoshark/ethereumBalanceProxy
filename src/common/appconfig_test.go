package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppConfiguration(t *testing.T) {
	port := "1848"
	endpoints := "endpoints!,moreEndpoints"
	ac := NewAppConfiguration("1848", &endpoints)
	assert.Equal(t, port, ac.ListenPort)
	assert.Equal(t, 2, len(ac.Endpoints))
	assert.Contains(t, endpoints, ac.Endpoints[0])
	assert.Contains(t, endpoints, ac.Endpoints[1])
}
