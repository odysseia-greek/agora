package diogenes

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthClient(t *testing.T) {
	standardTicks := 10 * time.Millisecond
	tick := 10 * time.Millisecond

	t.Run("HealthInfo", func(t *testing.T) {
		fixtures := []string{"health"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)
		sut, err := testClient.Health(context.Background())
		assert.Nil(t, err)
		assert.True(t, sut)
	})

	t.Run("CheckForHealthStatus", func(t *testing.T) {
		fixtures := []string{"health"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)

		sut := testClient.CheckHealthyStatus(context.Background(), standardTicks, tick)
		assert.Nil(t, err)
		assert.True(t, sut)
	})
}
