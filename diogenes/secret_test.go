package diogenes

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretClient(t *testing.T) {
	key := "testKey"
	value := "testValue"
	jsonBody := map[string]string{
		key: value,
	}

	name := "someTest"

	t.Run("CreateSecret", func(t *testing.T) {
		fixtures := []string{"createSecret"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)

		payload, err := json.Marshal(jsonBody)
		assert.Nil(t, err)

		sut, err := testClient.CreateNewSecret(context.Background(), name, payload)
		assert.Nil(t, err)
		assert.True(t, sut)
	})

	t.Run("RetrieveSecret", func(t *testing.T) {
		fixtures := []string{"retrieveSecret"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)

		secret, err := testClient.GetSecret(context.Background(), fixtureSecretName)
		assert.Nil(t, err)
		sut := fmt.Sprintf("%v", secret.Data)
		assert.Contains(t, sut, "value1")
		assert.Contains(t, sut, "key2")
	})
}
