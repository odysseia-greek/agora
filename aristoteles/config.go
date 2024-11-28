package aristoteles

import (
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"os"
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

func ElasticConfig(tls bool) (models.Config, error) {
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
		cert, err := GetCert()
		if err != nil {
			return models.Config{}, err
		}
		elasticCert = string(cert)
	}

	elasticService := ElasticService(tls)

	esConf := models.Config{
		Service:     elasticService,
		Username:    elasticUser,
		Password:    elasticPassword,
		ElasticCERT: elasticCert,
	}

	return esConf, nil
}

func GetCert() ([]byte, error) {
	cert, err := os.ReadFile(certPathInPod)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
