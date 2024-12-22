package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetConf(t *testing.T) {
	l := zap.NewNop().Sugar()

	server, err := New(l)

	require.NoError(t, err)

	assert.Equal(t, server.Config.Address, "127.0.0.1:8080")
	assert.Equal(t, server.Config.DBdsn, "./stats.db")
}
