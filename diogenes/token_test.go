package diogenes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenClient(t *testing.T) {
	t.Run("CreateToken", func(t *testing.T) {
		fixtures := []string{"token"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)
		expected := "s."

		policy := []string{"sdgdfg"}

		sut, err := testClient.CreateOneTimeToken(policy)
		assert.Nil(t, err)
		assert.Contains(t, sut, expected)
	})

}
