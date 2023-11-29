package aristoteles

import (
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDocument(t *testing.T) {
	index := "test"
	body := []byte(`{"Greek":"μάχη","English":"battle"}`)

	t.Run("Created", func(t *testing.T) {
		file := "createDocument"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		created, err := testClient.Document().Create(index, body)
		assert.Nil(t, err)
		assert.Equal(t, index, created.Index)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "createIndex"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		created, err := testClient.Document().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		created, err := testClient.Document().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})

	t.Run("NoConnection", func(t *testing.T) {
		config := models.Config{
			Service:     "hhttttt://sjdsj.com",
			Username:    "",
			Password:    "",
			ElasticCERT: "",
		}
		testClient, err := NewClient(config)
		assert.Nil(t, err)

		created, err := testClient.Document().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})
}
