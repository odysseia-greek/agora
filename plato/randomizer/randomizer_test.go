package randomizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestRandomizer(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		r, err := NewRandomizerClient()
		assert.Nil(t, err)
		assert.NotNil(t, r)
	})

	t.Run("CreateRandomNumberBaseZero", func(t *testing.T) {
		r, err := NewRandomizerClient()
		assert.Nil(t, err)
		assert.NotNil(t, r)

		for i := 1; i < 10; i++ {
			randomNumber := r.RandomNumberBaseZero(2)
			assert.True(t, randomNumber == 0 || randomNumber == 1)
		}
	})

	t.Run("CreateRandomNumberBaseOne", func(t *testing.T) {
		r, err := NewRandomizerClient()
		assert.Nil(t, err)
		assert.NotNil(t, r)

		for i := 1; i < 10; i++ {
			randomNumber := r.RandomNumberBaseOne(2)
			assert.True(t, randomNumber == 1 || randomNumber == 2)
		}
	})

	t.Run("CreateRandomNumberSetWithVariation", func(t *testing.T) {
		r, err := NewRandomizerClient()
		assert.Nil(t, err)
		assert.NotNil(t, r)

		var numbers []int
		for i := 1; i < 20; i++ {
			randomNumber := r.RandomNumberBaseOne(2000)
			assert.True(t, randomNumber < 2000+1)
			numbers = append(numbers, randomNumber)
		}
		sort.Ints(numbers)
		fmt.Println(numbers)
	})
}
