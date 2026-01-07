package aristoteles

import (
	"context"
	"fmt"
	"testing"

	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/stretchr/testify/assert"
)

func TestQueryClientMatch(t *testing.T) {
	index := "test"
	match := "elastic"
	word := "isGreat"
	expectedMalformed := "invalid character"

	t.Run("MatchPass", func(t *testing.T) {
		file := "match"
		status := 200
		expected := int64(1)
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().Match(index, body)
		assert.Nil(t, err)
		assert.Equal(t, expected, sut.Hits.Total.Value)
	})

	t.Run("SearchMultiple", func(t *testing.T) {
		files := []string{"match", "searchWord"}
		status := 200
		expectedFirstCall := int64(1)
		expectedSecondCall := int64(2)
		testClient, err := NewMockClient(files, status)
		assert.Nil(t, err)

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().Match(index, body)
		assert.Nil(t, err)
		assert.Equal(t, expectedFirstCall, sut.Hits.Total.Value)
		sut2, err := testClient.Query().Match(index, body)
		assert.Nil(t, err)
		assert.Equal(t, expectedSecondCall, sut2.Hits.Total.Value)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "serviceDown"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().Match(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), errorMessage)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().Match(index, body)
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

		sut, err := testClient.Query().Match(index, invalidBody)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
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

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().Match(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
	})
}

func TestQueryClientCount(t *testing.T) {
	index := "test"

	t.Run("CountPass", func(t *testing.T) {
		file := "count"
		status := 200
		expected := int64(1132)
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		ctx := context.Background()
		sut, err := testClient.Query().CountRaw(ctx, index, nil)
		assert.Nil(t, err)
		assert.Equal(t, expected, sut.Count)
	})
}

func TestQueryClientMatchWithScroll(t *testing.T) {
	index := "test"
	match := "elastic"
	word := "isGreat"
	expectedMalformed := "invalid character"

	t.Run("MatchPass", func(t *testing.T) {
		file := "createQuestionSokrates"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().MatchWithScroll(index, body)
		assert.Nil(t, err)
		assert.Equal(t, len(sut.Hits.Hits), 5)
	})

	t.Run("Failed", func(t *testing.T) {
		file := "serviceDown"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().MatchWithScroll(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), errorMessage)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().MatchWithScroll(index, body)
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

		sut, err := testClient.Query().MatchWithScroll(index, invalidBody)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
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

		body := testClient.Builder().MatchQuery(match, word)

		sut, err := testClient.Query().MatchWithScroll(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
	})
}

func TestQueryClientMatchAggregate(t *testing.T) {
	index := "test"
	expectedMalformed := "invalid character"
	aggregate := "authors"
	field := "author.keyword"
	expectedAuthors := [3]string{"herodotos", "ploutarchos", "thucydides"}

	t.Run("MatchPass", func(t *testing.T) {
		file := "match"
		status := 200
		expected := int64(1)
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().Aggregate(aggregate, field)
		sut, err := testClient.Query().MatchAggregate(index, body)
		assert.Nil(t, err)
		assert.Equal(t, expected, sut.Hits.Total.Value)
		for _, bucket := range sut.Aggregations.AuthorAggregation.Buckets {
			assert.Contains(t, expectedAuthors, bucket.Key)
		}
	})

	t.Run("Failed", func(t *testing.T) {
		file := "serviceDown"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().Aggregate(aggregate, field)
		sut, err := testClient.Query().MatchAggregate(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), errorMessage)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().Aggregate(aggregate, field)
		sut, err := testClient.Query().MatchAggregate(index, body)
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

		sut, err := testClient.Query().MatchAggregate(index, invalidBody)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
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

		body := testClient.Builder().Aggregate(aggregate, field)
		sut, err := testClient.Query().MatchAggregate(index, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
	})
}

func TestQueryClientSort(t *testing.T) {
	index := "test"
	expectedMalformed := "invalid character"
	sort := "author"
	mode := "desc"
	expectedKey := "firstKey"
	expectedValue := "firstValue"
	size := 1
	query := []map[string]string{
		{
			expectedKey: expectedValue,
		},
		{
			fmt.Sprintf("%s-2", expectedKey): fmt.Sprintf("%s-2", expectedValue),
		},
	}

	t.Run("MatchPass", func(t *testing.T) {
		file := "sorted"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MultipleMatch(query)
		sut, err := testClient.Query().MatchWithSort(index, mode, sort, size, body)
		assert.Nil(t, err)
		assert.Equal(t, size, len(sut.Hits.Hits))
	})

	t.Run("Failed", func(t *testing.T) {
		file := "serviceDown"
		status := 502
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MultipleMatch(query)
		sut, err := testClient.Query().MatchWithSort(index, mode, sort, size, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), errorMessage)
	})

	t.Run("Malformed", func(t *testing.T) {
		file := "malformed"
		status := 200
		testClient, err := NewMockClient(file, status)
		assert.Nil(t, err)

		body := testClient.Builder().MultipleMatch(query)
		sut, err := testClient.Query().MatchWithSort(index, mode, sort, size, body)
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

		sut, err := testClient.Query().MatchWithSort(index, mode, sort, size, invalidBody)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
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

		body := testClient.Builder().MultipleMatch(query)
		sut, err := testClient.Query().MatchWithSort(index, mode, sort, size, body)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
	})
}
