package diogenes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsealClient(t *testing.T) {
	t.Run("Unseal", func(t *testing.T) {
		fixtures := []string{"sealed", "sealed", "sealed", "unsealed"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)
		// Define the keys for unsealing
		keys := []string{"key1", "key2", "key3"}

		// Call the Unseal method
		unsealed, err := testClient.Unseal(keys)
		assert.Nil(t, err)

		// Assertions
		assert.NoError(t, err)
		assert.True(t, unsealed)
	})

	t.Run("Status", func(t *testing.T) {
		fixtures := []string{"sealed"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)

		// Call the Unseal method
		status, err := testClient.Status()
		assert.Nil(t, err)

		// Assertions
		assert.NoError(t, err)
		assert.True(t, status.Initialized)
		assert.True(t, status.Sealed)
	})

	t.Run("Initialize", func(t *testing.T) {
		fixtures := []string{"initialized"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)

		// Call the Unseal method
		initialized, err := testClient.Initialize(1, 1)
		assert.Nil(t, err)

		// Assertions
		assert.NoError(t, err)
		assert.Contains(t, initialized.Keys, "one")
	})
}
