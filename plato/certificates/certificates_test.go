package certificates

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGeneration(t *testing.T) {
	hosts := []string{
		"perikles",
		"perikles.odysseia",
		"perikles.odysseia.svc",
		"perikles.odysseia.svc.cluster.local",
	}
	organizations := []string{"test"}
	validityCa := 3650
	validityCert := 3650

	t.Run("Authority", func(t *testing.T) {
		impl, err := NewCertGeneratorClient(organizations, validityCa)
		assert.Nil(t, err)
		assert.NotNil(t, impl)
		err = impl.InitCa()
		assert.Nil(t, err)

		crt, key, err := impl.GenerateKeyAndCertSet(hosts, validityCert)
		log.Print(string(crt))
		log.Print(string(key))
		assert.Nil(t, err)
		assert.NotNil(t, key)
		assert.NotNil(t, crt)
	})

}
