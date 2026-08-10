package diogenes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenClient(t *testing.T) {
	t.Run("CreateToken", func(t *testing.T) {
		fixtures := []string{"token"}
		testClient, err := CreateMockVaultClient(fixtures, 200)
		assert.Nil(t, err)
		expected := "s."

		policy := []string{"sdgdfg"}

		sut, err := testClient.CreateOneTimeToken(context.Background(), policy)
		assert.Nil(t, err)
		assert.Contains(t, sut, expected)
	})

}
