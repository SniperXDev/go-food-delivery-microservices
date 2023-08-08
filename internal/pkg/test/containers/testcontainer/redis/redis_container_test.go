package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Redis_Container(t *testing.T) {
	redis, err := NewRedisTestContainers().Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, redis)
}
