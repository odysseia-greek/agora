package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/certificates"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CreateOdysseiaClient() (service.OdysseiaClient, error) {
	serviceNames := []string{EnvSolonService, EnvHerodotosService, EnvAlexandrosService, EnvDionysiosService}
	serviceURLs := make(map[string]*url.URL)

	for _, serviceName := range serviceNames {
		serviceURL := os.Getenv(serviceName)
		if serviceURL == "" {
			continue
		}
		parsedURL, err := url.Parse(serviceURL)
		if err != nil {
			return nil, err
		}
		serviceURLs[serviceName] = parsedURL
	}

	config := service.ClientConfig{}

	rootPath := os.Getenv("CERT_ROOT")
	for serviceName, parsedURL := range serviceURLs {
		apiConfig := service.OdysseiaApi{
			Url:    parsedURL.Host,
			Scheme: parsedURL.Scheme,
			Cert:   nil,
		}

		if parsedURL.Scheme == "https" {
			certPath, keyPath, caPath, err := getCertPaths(rootPath, serviceName)
			if err != nil {
				logging.Error(err.Error())
				continue
			}

			if _, err := os.Stat(certPath); !errors.Is(err, os.ErrNotExist) {
				if _, err := os.Stat(keyPath); !errors.Is(err, os.ErrNotExist) {
					loadedCerts, err := tls.LoadX509KeyPair(certPath, keyPath)
					if err != nil {
						return nil, err
					}

					apiConfig.Cert = []tls.Certificate{loadedCerts}
				}
			}

			if config.Ca == nil {
				if _, err := os.Stat(caPath); !errors.Is(err, os.ErrNotExist) {
					config.Ca, _ = os.ReadFile(caPath)
				}
			}
		}

		// Set the API configuration for each service
		switch serviceName {
		case EnvSolonService:
			config.Solon = apiConfig
		case EnvHerodotosService:
			config.Herodotos = apiConfig
		case EnvAlexandrosService:
			config.Alexandros = apiConfig
		case EnvDionysiosService:
			config.Dionysios = apiConfig
		}
	}

	return service.NewClient(config)
}

func getCertPaths(rootPath, serviceName string) (certPath, keyPath, caPath string, err error) {
	lowerServiceName := strings.ToLower(serviceName)

	dirs, err := os.ReadDir(rootPath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read the directory %s: %w", rootPath, err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			dirName := strings.ToLower(dir.Name())
			if strings.Contains(lowerServiceName, dirName) {
				dirPath := filepath.Join(rootPath, dir.Name())
				certPath := filepath.Join(dirPath, "tls.crt")
				keyPath := filepath.Join(dirPath, "tls.key")
				caPath := filepath.Join(dirPath, "tls.pem")
				return certPath, keyPath, caPath, nil
			}
		}
	}

	return "", "", "", errors.New("no matching service directory found")
}

func CreateNewRandomizer() (randomizer.Random, error) {
	return randomizer.NewRandomizerClient()
}

func CreateCertClient(org []string) (certificates.CertClient, error) {
	envCaValidity := StringFromEnv(EnvCAValidity, DefaultCaValidity)
	caValidity, err := strconv.Atoi(envCaValidity)
	if err != nil {
		return nil, err
	}

	certClient, err := certificates.NewCertGeneratorClient(org, caValidity)
	if err != nil {
		return nil, err
	}

	return certClient, nil
}
