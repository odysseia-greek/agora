package aristoteles

import (
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"os"
	"path/filepath"
	"time"
)

const (
	certPathInPod            = "/app/config/elastic/tls.crt"
	elasticServiceDefault    = "http://localhost:9200"
	elasticServiceDefaultTlS = "https://localhost:9200"
	elasticUsernameDefault   = "elastic"
	elasticPasswordDefault   = "odysseia"
	EnvElasticService        = "ELASTIC_SEARCH_SERVICE"
	EnvElasticUser           = "ELASTIC_SEARCH_USER"
	EnvElasticPassword       = "ELASTIC_SEARCH_PASSWORD"
)

func HealthCheck(client Client) error {
	standardTicks := 120 * time.Second
	tick := 1 * time.Second

	healthy := client.Health().Check(standardTicks, tick)
	if !healthy {
		return fmt.Errorf("elasticClient unhealthy after %s ticks", standardTicks)
	}

	return nil
}

func ElasticService(tls bool) string {
	elasticService := os.Getenv(EnvElasticService)
	if elasticService == "" {
		if tls {
			elasticService = elasticServiceDefaultTlS
		} else {
			elasticService = elasticServiceDefault
		}
	}
	return elasticService
}

func ElasticConfig(env string, testOverwrite, tls bool) models.Config {
	elasticUser := os.Getenv(EnvElasticUser)
	if elasticUser == "" {
		elasticUser = elasticUsernameDefault
	}
	elasticPassword := os.Getenv(EnvElasticPassword)
	if elasticPassword == "" {
		elasticPassword = elasticPasswordDefault
	}

	var elasticCert string
	if tls {
		elasticCert = string(GetCert(env, testOverwrite))
	}

	elasticService := ElasticService(tls)

	esConf := models.Config{
		Service:     elasticService,
		Username:    elasticUser,
		Password:    elasticPassword,
		ElasticCERT: elasticCert,
	}

	return esConf
}

func GetCert(env string, testOverWrite bool) []byte {
	var cert []byte

	if testOverWrite {
		certPath := filepath.Join("eratosthenes", "elastic-test-cert.pem")

		cert, _ = os.ReadFile(certPath)

		return cert
	}

	cert, _ = os.ReadFile(certPathInPod)

	return cert
}
