package aristoteles

import (
	"strings"
	"testing"

	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateDocumentWithIndexClient(t *testing.T) {
	index := "test"
	body := []byte(`{"Greek":"μάχη","English":"battle"}`)

	t.Run("Created", func(t *testing.T) {
		file := "createDocument"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		created, err := testClient.Index().CreateDocument(index, body)
		assert.Nil(t, err)
		assert.Equal(t, index, created.Index)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "createIndex"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		created, err := testClient.Index().CreateDocument(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		created, err := testClient.Index().CreateDocument(index, body)
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

		created, err := testClient.Index().CreateDocument(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)
	})
}

func TestCreateIndexClient(t *testing.T) {
	index := "test"
	expectedMalformed := "invalid character"
	searchWord := "someWord"

	t.Run("Created", func(t *testing.T) {
		file := "createIndex"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().SearchAsYouTypeIndex(searchWord)

		sut, err := testClient.Index().Create(index, body)
		assert.Nil(t, err)
		assert.Equal(t, index, sut.Index)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "createIndex"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().SearchAsYouTypeIndex(searchWord)

		sut, err := testClient.Index().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().SearchAsYouTypeIndex(searchWord)

		sut, err := testClient.Index().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), expectedMalformed)
	})

	t.Run("Unparseable", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		invalidBody := map[string]interface{}{
			"key": make(chan int),
		}

		created, err := testClient.Index().Create(index, invalidBody)
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

		body := testClient.Builder().SearchAsYouTypeIndex(searchWord)
		created, err := testClient.Index().Create(index, body)
		assert.NotNil(t, err)
		assert.False(t, created.Acknowledged)
	})
}

func TestDeleteIndexClient(t *testing.T) {
	index := "test"
	expectedMalformed := "invalid character"

	t.Run("Deleted", func(t *testing.T) {
		file := "deleteIndex"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.Nil(t, err)
		assert.True(t, sut)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "deleteIndex"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.NotNil(t, err)
		assert.False(t, sut)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.NotNil(t, err)
		assert.False(t, sut)
		assert.Contains(t, err.Error(), expectedMalformed)
	})

	t.Run("IndexDoesNotExist", func(t *testing.T) {
		file := "deleteIndex404"
		status := 404
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		sut, err := testClient.Index().Delete(index)
		assert.NotNil(t, err)
		assert.False(t, sut)
		assert.Contains(t, err.Error(), "index_not_found_exception")
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

		sut, err := testClient.Index().Delete(index)
		assert.NotNil(t, err)
		assert.False(t, sut)
	})
}

type seederHandler struct {
	Elastic    Client
	Index      string
	PolicyName string
}

func (s *seederHandler) DeleteIndexAtStartUp() error {
	deleted, err := s.Elastic.Index().Delete(s.Index)
	if err != nil {
		if deleted {
			return nil
		}
		if strings.Contains(err.Error(), "index_not_found_exception") {
			return nil
		}
		return err
	}

	return nil
}

func (s *seederHandler) CreateIndexAtStartup() error {
	indexMapping := s.Elastic.Builder().TextIndex(s.PolicyName)
	_, err := s.Elastic.Index().Create(s.Index, indexMapping)
	if err != nil {
		return err
	}

	return nil
}

func TestSeederDeleteThenCreateIndexFlow(t *testing.T) {
	t.Run("DeletedThenCreated", func(t *testing.T) {
		testClient, err := NewMockClient([]string{"deleteIndex", "createIndex"}, 200)
		assert.Nil(t, err)

		handler := &seederHandler{
			Elastic:    testClient,
			Index:      "parmenides-quiz",
			PolicyName: "parmenides-policy",
		}

		err = handler.DeleteIndexAtStartUp()
		assert.Nil(t, err)

		err = handler.CreateIndexAtStartup()
		assert.Nil(t, err)
	})

	t.Run("NotFoundThenCreated", func(t *testing.T) {
		testClient, err := NewMockClient("deleteIndex404", 404)
		assert.Nil(t, err)

		handler := &seederHandler{
			Elastic:    testClient,
			Index:      "parmenides-quiz",
			PolicyName: "parmenides-policy",
		}

		err = handler.DeleteIndexAtStartUp()
		assert.Nil(t, err)

		createClient, err := NewMockClient("createIndex", 200)
		assert.Nil(t, err)
		handler.Elastic = createClient

		err = handler.CreateIndexAtStartup()
		assert.Nil(t, err)
	})
}
