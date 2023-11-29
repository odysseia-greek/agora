package thales

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultKubeConfig = "/.kube/config"
	defaultNamespace  = "odysseia"
	EnvNamespace      = "NAMESPACE"
	EnvKubePath       = "KUBE_PATH"
)

func CreateKubeClient(OutOfClusterKube bool) (KubeClient, error) {
	var kubeManager KubeClient

	namespace := os.Getenv(EnvNamespace)
	if namespace == "" {
		namespace = defaultNamespace
	}

	kubePath := os.Getenv(EnvKubePath)

	if OutOfClusterKube {
		var filePath string
		if kubePath == "" {
			log.Printf("defaulting to %s", defaultKubeConfig)
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Print(err)
			}

			filePath = filepath.Join(homeDir, defaultKubeConfig)
		} else {
			filePath = kubePath
		}

		cfg, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Print("error getting kubeconfig")
		}

		kube, err := NewKubeClient(cfg, namespace)
		if err != nil {
			log.Fatal("error creating kubeclient")
		}

		kubeManager = kube
	} else {
		log.Print("creating in cluster kube client")
		kube, err := NewKubeClient(nil, namespace)
		if err != nil {
			log.Fatal("error creating kubeclient")
		}
		kubeManager = kube
	}

	return kubeManager, nil
}
