package aristoteles

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"log"
	"time"
)

type HealthImpl struct {
	es *elasticsearch.Client
}

func NewHealthImpl(suppliedClient *elasticsearch.Client) (*HealthImpl, error) {
	return &HealthImpl{es: suppliedClient}, nil
}

func (h *HealthImpl) Check(ticks, tick time.Duration) bool {
	healthy := false

	ticker := time.NewTicker(tick)
	timeout := time.After(ticks)

	for {
		select {
		case t := <-ticker.C:
			log.Printf("tick: %s", t)
			res := h.Info()
			healthy = res.Healthy
			if !healthy {
				log.Print("Elastic not yet healthy")
				continue
			}

			ticker.Stop()

		case <-timeout:
			ticker.Stop()
		}
		break
	}

	return healthy
}

func (h *HealthImpl) Info() (elasticHealth models.DatabaseHealth) {
	res, err := h.es.Info()

	if err != nil {
		elasticHealth.Healthy = false
		return elasticHealth
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		elasticHealth.Healthy = false
		return elasticHealth
	}

	var r map[string]interface{}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		elasticHealth.Healthy = false
		return elasticHealth
	}

	elasticHealth.ClusterName = fmt.Sprintf("%s", r["cluster_name"])
	elasticHealth.ServerName = fmt.Sprintf("%s", r["name"])
	elasticHealth.ServerVersion = fmt.Sprintf("%s", r["version"].(map[string]interface{})["number"])
	elasticHealth.Healthy = true

	return elasticHealth
}
