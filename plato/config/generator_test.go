package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)
	defaultPath := "odysseia-greek"

	t.Run("LocalFilePath", func(t *testing.T) {
		filePath := OdysseiaRootPath(defaultPath)
		sut := filepath.Join(homeDir, "/go/src/github.com/", defaultPath)
		assert.Equal(t, sut, filePath)
	})

	t.Run("EmptyPath", func(t *testing.T) {
		filePath := OdysseiaRootPath("")
		sut := "/"
		assert.Equal(t, sut, filePath)
	})

	t.Run("PlatoPathFlat", func(t *testing.T) {
		filePath := OdysseiaRootPath(defaultPath)
		platoPath := PlatoPath(filePath)
		sut := filepath.Join(homeDir, "/go/src/github.com/", defaultPath, PLATO)
		assert.Equal(t, sut, platoPath)
	})

}
