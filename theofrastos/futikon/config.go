package futikon

import (
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/thales"
	"os"
)

func CreateNewConfig() (*TheofratosHandler, error) {
	logging.Debug("creating config")

	tls := config.BoolFromEnv(config.EnvTlSKey)
	cfg, err := aristoteles.ElasticConfig(tls)
	if err != nil {
		logging.Error(fmt.Sprintf("failed to create Elastic client operations will be interupted, %s", err.Error()))
	}

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	err = aristoteles.HealthCheck(elastic)
	if err != nil {
		return nil, err
	}

	configMapName := os.Getenv("CONFIGMAP_NAME")
	if configMapName == "" {

	}

	ns := os.Getenv(config.EnvNamespace)
	if ns == "" {

	}

	kube, err := thales.CreateKubeClient(false)
	if err != nil {
		return nil, err
	}

	return &TheofratosHandler{
		Elastic:       elastic,
		Namespace:     ns,
		ConfigMapName: configMapName,
		Kube:          kube,
	}, nil
}
