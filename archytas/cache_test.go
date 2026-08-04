package archytas

import (
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

func TestBadgerClient(t *testing.T) {
	key := "testkey"
	value := "testvalue"

	t.Run("ReadValue", func(t *testing.T) {
		testClient, err := NewInMemoryBadgerClient()
		assert.Nil(t, err)

		err = testClient.Set(key, value)
		assert.Nil(t, err)

		sut, err := testClient.Read(key)
		assert.Equal(t, value, string(sut))
		testClient.Close()
	})

	t.Run("ReadEmptyValue", func(t *testing.T) {
		testClient, err := NewInMemoryBadgerClient()
		assert.Nil(t, err)

		sut, err := testClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, sut)

		testClient.Close()
	})

	t.Run("SetValueWithTTL", func(t *testing.T) {
		testClient, err := NewInMemoryBadgerClient()
		assert.Nil(t, err)

		ttl := time.Millisecond
		err = testClient.SetWithTTL(key, value, ttl)
		assert.Nil(t, err)

		time.Sleep(2 * ttl)
		sut, err := testClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		testClient.Close()
	})

	t.Run("SetValueWithTTL", func(t *testing.T) {
		testClient, err := NewInMemoryBadgerClient()
		assert.Nil(t, err)

		ttl := 1000 * time.Millisecond
		err = testClient.SetWithTTL(key, value, ttl)
		assert.Nil(t, err)
		sut, err := testClient.Read(key)
		assert.NotNil(t, sut)
		assert.Nil(t, err)

		time.Sleep(ttl)
		sut, err = testClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		testClient.Close()
	})

	t.Run("AbilityToOpenTwoDatabases", func(t *testing.T) {
		uuid := uuid2.New().String()
		badgerPath := filepath.Join("/tmp", "badger", uuid)
		testClient, err := NewBadgerClient(badgerPath)
		assert.Nil(t, err)

		sut, err := testClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, sut)

		newUUID := uuid2.New().String()
		newBadgerPath := filepath.Join("/tmp", "badger", newUUID)
		newTestClient, err := NewBadgerClient(newBadgerPath)
		newSut, err := newTestClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, newSut)

		testClient.Close()
		newTestClient.Close()
	})
}

func TestBadgerConvenienceMethods(t *testing.T) {
	client, err := NewInMemoryBadgerClientWithOptions(WithLogging(false))
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, client.Close()) })

	assert.NoError(t, client.SetBytes("users:1", []byte{0, 1, 2}))
	assert.NoError(t, client.Set("users:2", "Homer"))

	exists, err := client.Exists("users:1")
	assert.NoError(t, err)
	assert.True(t, exists)

	value, err := client.GetString("users:2")
	assert.NoError(t, err)
	assert.Equal(t, "Homer", value)

	stats, err := client.Stats()
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), stats.KeyCount)
	assert.Equal(t, stats.LSMSizeBytes+stats.ValueLogSizeBytes, stats.TotalSizeBytes)

	assert.NoError(t, client.Delete("users:1"))
	exists, err = client.Exists("users:1")
	assert.NoError(t, err)
	assert.False(t, exists)

	assert.NoError(t, client.DeletePrefix("users:"))
	stats, err = client.Stats()
	assert.NoError(t, err)
	assert.Zero(t, stats.KeyCount)
}

func TestWithLoggingDisabled(t *testing.T) {
	client, err := NewInMemoryBadgerClientWithOptions(WithLogging(false))
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, client.Close()) })

	implementation := client.(*Badger)
	assert.Nil(t, implementation.db.Opts().Logger)
}
