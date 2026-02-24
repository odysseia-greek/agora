package stomion

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
)

const (
	DefaultServiceName string = "eupalinos"
)

func CreateNewConfig() (*QueueServiceImpl, error) {

	replicasFromEnv := config.StringFromEnv("TOTAL_REPLICAS", "1")
	rootPath := config.StringFromEnv(config.EnvRootTlSDir, "")
	podName := config.StringFromEnv(config.EnvPodName, "eupalinos-0")
	namespace := config.StringFromEnv(config.EnvNamespace, "agora")
	serviceName := config.StringFromEnv("SERVICE_NAME", DefaultServiceName)

	var podNumber string
	podNumber = strings.Split(podName, "-")[1]
	podName = strings.Split(podName, "-")[0]

	logging.Debug(fmt.Sprintf("podNumber: %s", podNumber))

	podID, err := strconv.Atoi(podNumber)
	if err != nil {
		return nil, err
	}
	replicas, err := strconv.Atoi(replicasFromEnv)
	if err != nil {
		return nil, err
	}

	logging.Debug(fmt.Sprintf("podName: %s", podName))
	logging.Debug(fmt.Sprintf("replicas: %d", replicas))
	logging.Debug(fmt.Sprintf("podID: %d", podID))

	// Calculate the addresses for each replica based on the Pod ID
	addresses := make([]string, replicas-1)
	addrIdx := 0
	if replicas > 1 {
		for i := 0; i < replicas; i++ {
			if i == podID {
				continue
			} else {
				addresses[addrIdx] = fmt.Sprintf("%s-%d.%s.%s.svc.cluster.local:50051", podName, i, serviceName, namespace)
				logging.Debug(fmt.Sprintf("address added at: %s", addresses[addrIdx]))
			}
			addrIdx++
		}
	}

	var streaming bool
	if replicas > 1 {
		streaming = true
	}

	tlsConfig, err := loadTLSConfig(rootPath, serviceName)
	if err != nil {
		logging.Error(err.Error())
	}

	version := os.Getenv(config.EnvVersion)

	return &QueueServiceImpl{
		Version:     version,
		DiexodosMap: make([]*Diexodos, 0),
		mu:          sync.Mutex{},
		Addresses:   addresses,
		Streaming:   streaming,
		TLSConfig:   tlsConfig,
	}, nil
}

// Function to load TLS configuration for HTTPS mode
func loadTLSConfig(rootPath, serviceName string) (*tls.Config, error) {
	certFile := filepath.Join(rootPath, serviceName, "tls.crt")
	keyFile := filepath.Join(rootPath, serviceName, "tls.key")
	caCertFile := filepath.Join(rootPath, serviceName, "tls.pem")

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
