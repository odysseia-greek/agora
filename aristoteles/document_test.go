package aristoteles

import (
	"errors"
	"testing"

	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/stretchr/testify/assert"
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

	t.Run("AuthenticationUnauthorized401", func(t *testing.T) {
		file := "authenticationError401"
		status := 401
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		created, err := testClient.Document().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)

		var elasticErr *ElasticError
		assert.True(t, errors.As(err, &elasticErr))
		assert.Equal(t, 401, elasticErr.StatusCode)
		assert.Equal(t, "create document", elasticErr.Operation)

		errorDetail, ok := elasticErr.Detail.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "security_exception", errorDetail["type"])
		assert.Contains(t, errorDetail["reason"], "unable to authenticate user")
	})

	t.Run("AuthorizationForbidden403", func(t *testing.T) {
		file := "authorizationError403"
		status := 403
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		created, err := testClient.Document().Create(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)

		var elasticErr *ElasticError
		assert.True(t, errors.As(err, &elasticErr))
		assert.Equal(t, 403, elasticErr.StatusCode)
		assert.Equal(t, "create document", elasticErr.Operation)

		errorDetail, ok := elasticErr.Detail.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "security_exception", errorDetail["type"])
		assert.Contains(t, errorDetail["reason"], "is unauthorized for user")
	})
}

func TestCreateWithId(t *testing.T) {
	index := "test"
	documentID := "my-document-id"
	body := []byte(`{"Greek":"μάχη","English":"battle"}`)

	t.Run("Created", func(t *testing.T) {
		testClient, err := NewMockClient("createWithId", 200)
		assert.Nil(t, err)

		created, err := testClient.Document().CreateWithId(index, documentID, body)
		assert.Nil(t, err)
		assert.Equal(t, "created", created.Result)
	})

	t.Run("Failed", func(t *testing.T) {
		testClient, err := NewMockClient("serviceDown", 502)
		assert.Nil(t, err)

		created, err := testClient.Document().CreateWithId(index, documentID, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)

		var elasticErr *ElasticError
		assert.True(t, errors.As(err, &elasticErr))
		assert.Equal(t, 502, elasticErr.StatusCode)
		assert.Equal(t, "create document with id", elasticErr.Operation)
	})

	t.Run("Malformed", func(t *testing.T) {
		testClient, err := NewMockClient("malformed", 200)
		assert.Nil(t, err)

		created, err := testClient.Document().CreateWithId(index, documentID, body)
		assert.NotNil(t, err)
		assert.Nil(t, created)
		assert.Contains(t, err.Error(), "invalid character")
	})
}

func TestCreateWithIdAndFirstItem(t *testing.T) {
	index := "test"
	documentID := "my-document-id"
	itemBody := `{"word":"μαχη"}`
	paramName := "items"

	t.Run("Created", func(t *testing.T) {
		testClient, err := NewMockClient("createWithId", 200)
		assert.Nil(t, err)

		created, err := testClient.Document().CreateWithIdAndFirstItem(index, documentID, itemBody, paramName)
		assert.Nil(t, err)
		assert.Equal(t, "created", created.Result)
	})

	t.Run("Failed", func(t *testing.T) {
		testClient, err := NewMockClient("serviceDown", 502)
		assert.Nil(t, err)

		created, err := testClient.Document().CreateWithIdAndFirstItem(index, documentID, itemBody, paramName)
		assert.NotNil(t, err)
		assert.Nil(t, created)

		var elasticErr *ElasticError
		assert.True(t, errors.As(err, &elasticErr))
		assert.Equal(t, 502, elasticErr.StatusCode)
		assert.Equal(t, "create document with first item", elasticErr.Operation)
	})

	t.Run("Malformed", func(t *testing.T) {
		testClient, err := NewMockClient("malformed", 200)
		assert.Nil(t, err)

		created, err := testClient.Document().CreateWithIdAndFirstItem(index, documentID, itemBody, paramName)
		assert.NotNil(t, err)
		assert.Nil(t, created)
		assert.Contains(t, err.Error(), "invalid character")
	})
}

func TestUpdateDocument(t *testing.T) {
	index := "test"
	documentID := "my-document-id"
	body := []byte(`{"English":"combat"}`)

	t.Run("Updated", func(t *testing.T) {
		testClient, err := NewMockClient("updated", 200)
		assert.Nil(t, err)

		updated, err := testClient.Document().Update(index, documentID, body)
		assert.Nil(t, err)
		assert.Equal(t, "updated", updated.Result)
	})

	t.Run("Failed", func(t *testing.T) {
		testClient, err := NewMockClient("serviceDown", 502)
		assert.Nil(t, err)

		updated, err := testClient.Document().Update(index, documentID, body)
		assert.NotNil(t, err)
		assert.Nil(t, updated)

		var elasticErr *ElasticError
		assert.True(t, errors.As(err, &elasticErr))
		assert.Equal(t, 502, elasticErr.StatusCode)
		assert.Equal(t, "update document", elasticErr.Operation)
	})

	t.Run("Malformed", func(t *testing.T) {
		testClient, err := NewMockClient("malformed", 200)
		assert.Nil(t, err)

		updated, err := testClient.Document().Update(index, documentID, body)
		assert.NotNil(t, err)
		assert.Nil(t, updated)
		assert.Contains(t, err.Error(), "invalid character")
	})
}

func TestAddItemToDocument(t *testing.T) {
	index := "test"
	documentID := "my-document-id"
	itemBody := `{"word":"αγών"}`
	paramName := "items"

	t.Run("Updated", func(t *testing.T) {
		testClient, err := NewMockClient("updated", 200)
		assert.Nil(t, err)

		updated, err := testClient.Document().AddItemToDocument(index, documentID, itemBody, paramName)
		assert.Nil(t, err)
		assert.Equal(t, "updated", updated.Result)
	})

	t.Run("Failed", func(t *testing.T) {
		testClient, err := NewMockClient("serviceDown", 502)
		assert.Nil(t, err)

		updated, err := testClient.Document().AddItemToDocument(index, documentID, itemBody, paramName)
		assert.NotNil(t, err)
		assert.Nil(t, updated)

		var elasticErr *ElasticError
		assert.True(t, errors.As(err, &elasticErr))
		assert.Equal(t, 502, elasticErr.StatusCode)
		assert.Equal(t, "add item to document", elasticErr.Operation)
	})

	t.Run("Malformed", func(t *testing.T) {
		testClient, err := NewMockClient("malformed", 200)
		assert.Nil(t, err)

		updated, err := testClient.Document().AddItemToDocument(index, documentID, itemBody, paramName)
		assert.NotNil(t, err)
		assert.Nil(t, updated)
		assert.Contains(t, err.Error(), "invalid character")
	})
}
