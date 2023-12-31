package diogenes

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHealthClient(t *testing.T) {
	standardTicks := 10 * time.Millisecond
	tick := 10 * time.Millisecond

	t.Run("HealthInfo", func(t *testing.T) {
		fixtures := []string{"health"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)
		sut, err := testClient.Health()
		assert.Nil(t, err)
		assert.True(t, sut)
	})

	t.Run("CheckForHealthStatus", func(t *testing.T) {
		fixtures := []string{"health"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)

		sut := testClient.CheckHealthyStatus(standardTicks, tick)
		assert.Nil(t, err)
		assert.True(t, sut)
	})
}
