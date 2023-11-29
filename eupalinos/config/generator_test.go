package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateNewConfig(t *testing.T) {
	t.Run("LocalEnv", func(t *testing.T) {
		// Set up test environment variables
		os.Setenv("POD_INDEX", "1")
		os.Setenv("TOTAL_REPLICAS", "2")
		os.Setenv("CERT_ROOT", "/path/to/certs")
		os.Setenv("POD_NAME", "eupalinos-0")
		os.Setenv("SAVE_PATH", "/tmp")
		os.Setenv("NAMESPACE", "mynamespace")
		os.Setenv("SERVICE_NAME", "myservice")

		config, err := CreateNewConfig("LOCAL")

		assert.NoError(t, err, "Unexpected error")
		assert.NotNil(t, config, "Config should not be nil")

		assert.True(t, config.Streaming)
		assert.True(t, len(config.Addresses) == 1)
		// Clean up environment variables after the test
		os.Clearenv()
	})

	t.Run("KubeEnv", func(t *testing.T) {
		// Set up test environment variables
		os.Setenv("POD_INDEX", "1")
		os.Setenv("TOTAL_REPLICAS", "2")
		os.Setenv("CERT_ROOT", "/path/to/certs")
		os.Setenv("POD_NAME", "eupalinos-0")
		os.Setenv("SAVE_PATH", "/tmp")
		os.Setenv("NAMESPACE", "mynamespace")
		os.Setenv("SERVICE_NAME", "myservice")

		config, err := CreateNewConfig("KUBE")

		assert.NoError(t, err, "Unexpected error")
		assert.NotNil(t, config, "Config should not be nil")

		assert.True(t, config.Streaming)
		assert.True(t, len(config.Addresses) == 1)
		// Clean up environment variables after the test
		os.Clearenv()
	})

	t.Run("EmptyEnv", func(t *testing.T) {
		config, err := CreateNewConfig("LOCAL")

		assert.NoError(t, err, "Unexpected error")
		assert.NotNil(t, config, "Config should not be nil")

		assert.False(t, config.Streaming)
		// Clean up environment variables after the test
		os.Clearenv()
	})

}

func TestLoadTLSConfig(t *testing.T) {
	t.Run("TLSConfigValid", func(t *testing.T) {
		tlsConfig, err := loadTLSConfig("fixtures", "certs")

		assert.Nil(t, err)
		assert.NotNil(t, tlsConfig)
	})

	t.Run("TLSConfigInValid", func(t *testing.T) {
		tlsConfig, err := loadTLSConfig("/path/to/certs", "myservice")

		assert.NotNil(t, err)
		assert.Nil(t, tlsConfig)
	})
}
