package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultServiceName string = "eupalinos"
)

func CreateNewConfig(env string) (*Config, error) {
	// POD_INDEX for local development only
	replicasFromEnv := os.Getenv("TOTAL_REPLICAS")
	rootPath := os.Getenv("CERT_ROOT")
	podName := os.Getenv("POD_NAME")
	savePath := os.Getenv("SAVE_PATH")
	if rootPath == "" {
		log.Print("rootpath is empty no certs can be loaded")
	}

	if savePath == "" {
		savePath = "/tmp"
	}

	if podName == "" {
		podName = "eupalinos-0"
	}

	// Get the Namespace from the environment
	namespace := os.Getenv("NAMESPACE")
	serviceName := os.Getenv("SERVICE_NAME")

	if serviceName == "" {
		serviceName = DefaultServiceName
	}

	if replicasFromEnv == "" {
		replicasFromEnv = "1"
	}
	var podNumber string
	podNumber = strings.Split(podName, "-")[1]
	podName = strings.Split(podName, "-")[0]

	log.Printf("podNumber: %s", podNumber)

	podID, err := strconv.Atoi(podNumber)
	if err != nil {
		return nil, err
	}
	replicas, err := strconv.Atoi(replicasFromEnv)
	if err != nil {
		return nil, err
	}

	log.Printf("podName: %s", podName)
	log.Printf("replicas: %d", replicas)
	log.Printf("podID: %d", podID)

	// Calculate the addresses for each replica based on the Pod ID
	addresses := make([]string, replicas-1)
	addrIdx := 0
	if replicas > 1 {
		for i := 0; i < replicas; i++ {
			if i == podID {
				continue
			} else if env == "LOCAL" {
				addresses[addrIdx] = fmt.Sprintf("localhost:5005%d", i+1)
				log.Printf("address added at: %s", addresses[addrIdx])
			} else {
				addresses[addrIdx] = fmt.Sprintf("%s-%d.%s.%s.svc.cluster.local:50051", podName, i, serviceName, namespace)
				log.Printf("address added at: %s", addresses[addrIdx])
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
		log.Print(err)
	}

	subPath := fmt.Sprintf("%s%d", podName, podID)
	backUpPath := filepath.Join(savePath, subPath, "eupalinos_state.json")

	return &Config{
		Addresses: addresses,
		Streaming: streaming,
		TLSConfig: tlsConfig,
		SavePath:  backUpPath,
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

	caCert, err := ioutil.ReadFile(caCertFile)
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
