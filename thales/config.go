package thales

import (
	"os"
	"path/filepath"
)

const (
	defaultKubeConfig = "/.kube/config"
	defaultNamespace  = "odysseia"
	EnvNamespace      = "NAMESPACE"
	EnvKubePath       = "KUBE_PATH"
)

func CreateKubeClient(OutOfClusterKube bool) (*KubeClient, error) {
	var kubeManager *KubeClient

	namespace := os.Getenv(EnvNamespace)
	if namespace == "" {
		namespace = defaultNamespace
	}

	kubePath := os.Getenv(EnvKubePath)

	if OutOfClusterKube {
		var filePath string
		if kubePath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}

			filePath = filepath.Join(homeDir, defaultKubeConfig)
		} else {
			filePath = kubePath
		}

		cfg, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		kube, err := NewFromConfig(cfg)
		if err != nil {
			return nil, err
		}

		kubeManager = kube
	} else {
		kube, err := NewInClusterKube()
		if err != nil {
			return nil, err
		}
		kubeManager = kube
	}

	return kubeManager, nil
}
