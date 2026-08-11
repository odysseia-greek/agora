package diogenes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKubernetesAuthMethodWithAudienceRejectsEmptyAudience(t *testing.T) {
	vault := &Vault{}

	err := vault.KubernetesAuthMethodWithAudience(
		context.Background(),
		"solon",
		"solon-access-sa",
		"delphi",
		"https://kubernetes.default.svc",
		"",
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "audience cannot be empty")
}
