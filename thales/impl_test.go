package thales

import (
	"github.com/odysseia-greek/thales/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestClientCreation(t *testing.T) {
	ns := "test"
	fixturePath := "./eratosthenes/"

	t.Run("ConfigHappyPath", func(t *testing.T) {
		config := "config.yaml"
		fp := filepath.Join(fixturePath, config)
		var cfg []byte
		cfg, _ = ioutil.ReadFile(fp)

		kube, err := NewKubeClient(cfg, ns)
		assert.NotNil(t, kube)
		assert.Nil(t, err)
	})

	t.Run("ConfigWithoutContextSetsFirst", func(t *testing.T) {
		config := "no-context.yaml"
		fp := filepath.Join(fixturePath, config)
		var cfg []byte
		cfg, _ = ioutil.ReadFile(fp)

		c, err := models.UnmarshalKubeConfig(cfg)
		assert.Nil(t, err)

		kube, err := NewKubeClient(cfg, ns)
		assert.NotNil(t, kube)
		assert.Nil(t, err)

		kubeContext, err := kube.Cluster().GetCurrentContext()
		assert.Equal(t, c.Contexts[0].Name, kubeContext)
		assert.Nil(t, err)
	})

	t.Run("ConfigWithoutContextSetsToDesktopWithMultiple", func(t *testing.T) {
		config := "no-context-multiple.yaml"
		fp := filepath.Join(fixturePath, config)
		var cfg []byte
		cfg, _ = ioutil.ReadFile(fp)

		kube, err := NewKubeClient(cfg, ns)
		assert.NotNil(t, kube)
		assert.Nil(t, err)

		kubeContext, err := kube.Cluster().GetCurrentContext()
		assert.Contains(t, kubeContext, "desktop")
		assert.Nil(t, err)
	})

	t.Run("EmptyConfig", func(t *testing.T) {
		config := "empty.yaml"
		fp := filepath.Join(fixturePath, config)
		var cfg []byte
		cfg, _ = ioutil.ReadFile(fp)

		kube, err := NewKubeClient(cfg, ns)
		assert.Nil(t, kube)
		assert.NotNil(t, err)
	})

	t.Run("MalformedCerts", func(t *testing.T) {
		config := "bad-certs.yaml"
		fp := filepath.Join(fixturePath, config)
		var cfg []byte
		cfg, _ = ioutil.ReadFile(fp)

		kube, err := NewKubeClient(cfg, ns)
		assert.Nil(t, kube)
		assert.NotNil(t, err)
	})
}
