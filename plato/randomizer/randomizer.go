package randomizer

import (
	"crypto/rand"
	"encoding/binary"
	mrand "math/rand"
	"time"
)

type Random interface {
	RandomNumberBaseZero(length int) int
	RandomNumberBaseOne(length int) int
}

type Randomizer struct {
	r1 *mrand.Rand
}

// RandomNumberBaseZero creates a random number starting from 0 to length - 1 (example length = 2 -> 0 or 1)
func (r *Randomizer) RandomNumberBaseZero(length int) int {
	return r.r1.Intn(length)
}

// RandomNumberBaseOne creates a random number starting from 1 to length (example length = 2 -> 1 or 2)
func (r *Randomizer) RandomNumberBaseOne(length int) int {
	return r.r1.Intn(length) + 1
}

func NewRandomizerClient() (*Randomizer, error) {
	var seed int64
	err := binary.Read(rand.Reader, binary.LittleEndian, &seed)
	if err != nil {
		// Fallback to time-based seed if crypto/rand fails
		seed = time.Now().UnixNano()
	}
	localRandomizer := mrand.NewSource(seed)
	r1 := mrand.New(localRandomizer)

	return &Randomizer{r1: r1}, nil
}
