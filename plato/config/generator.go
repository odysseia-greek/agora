package config

import (
	"crypto/tls"
	"errors"
	"github.com/odysseia-greek/agora/plato/certificates"
	"github.com/odysseia-greek/agora/plato/helpers"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func CreateOdysseiaClient() (service.OdysseiaClient, error) {
	serviceNames := []string{EnvSolonService, EnvHerodotosService, EnvAlexandrosService, EnvDionysiosService, EnvSokratesService}
	serviceURLs := make(map[string]*url.URL)

	for _, serviceName := range serviceNames {
		serviceURL := StringFromEnv(serviceName, DefaultServiceAddress)
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
			certPath, keyPath, caPath := getCertPaths(rootPath, serviceName)

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
		case EnvSokratesService:
			config.Sokrates = apiConfig
		}
	}

	return service.NewClient(config)
}

func getCertPaths(rootPath, serviceName string) (certPath, keyPath, caPath string) {
	dirPath := filepath.Join(rootPath, serviceName)
	return filepath.Join(dirPath, "tls.crt"), filepath.Join(dirPath, "tls.key"), filepath.Join(rootPath, serviceName, "tls.pem")
}

func RetrieveCertPathLocally(testOverwrite bool, service string) (cert string, key string) {
	keyName := "tls.key"
	certName := "tls.crt"

	if testOverwrite {
		rootPath := helpers.OdysseiaRootPath()
		if service == "" {
			service = "solon"
		}
		cert = filepath.Join(rootPath, "eratosthenes", "fixture", service, certName)
		key = filepath.Join(rootPath, "eratosthenes", "fixture", service, keyName)

		return
	} else {
		rootPath := os.Getenv("CERT_ROOT")
		cert = filepath.Join(rootPath, service, certName)
		key = filepath.Join(rootPath, service, keyName)

	}

	return
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

	//org := []string{
	//	"odysseia",
	//}

	certClient, err := certificates.NewCertGeneratorClient(org, caValidity)
	if err != nil {
		return nil, err
	}

	return certClient, nil
}
