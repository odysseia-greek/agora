package futikon

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newKubeClient() (kubernetes.Interface, error) {
	config, err := kubeConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func newFakeKubeClient() kubernetes.Interface {
	return fake.NewSimpleClientset()
}

func kubeConfig() (*rest.Config, error) {
	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}

	path := os.Getenv("KUBECONFIG")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".kube", "config")
	}

	return clientcmd.BuildConfigFromFlags("", path)
}
