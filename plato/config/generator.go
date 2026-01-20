package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/odysseia-greek/agora/plato/certificates"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
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
		return "", "", "", fmt.Errorf("failed to read directory %s: %w", rootPath, err)
	}

	// Prefer exact directory name matches, then fall back to "contains"
	var candidates []string
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		name := dir.Name()
		lowerName := strings.ToLower(name)

		if lowerName == lowerServiceName {
			// exact match goes first
			candidates = append([]string{name}, candidates...)
			continue
		}

		if strings.Contains(lowerServiceName, lowerName) || strings.Contains(lowerName, lowerServiceName) {
			candidates = append(candidates, name)
		}
	}

	if len(candidates) == 0 {
		return "", "", "", fmt.Errorf("no matching service directory found for %q under %s", serviceName, rootPath)
	}

	// Try each candidate dir until we find a usable set
	for _, dirName := range candidates {
		dirPath := filepath.Join(rootPath, dirName)

		cert := firstExistingFile(dirPath, []string{"tls.crt"})
		key := firstExistingFile(dirPath, []string{"tls.key"})
		ca := firstExistingFile(dirPath, []string{"ca.crt", "tls.pem"})

		if cert != "" && key != "" && ca != "" {
			return cert, key, ca, nil
		}
	}

	// If we get here, we found directories but none had a complete set.
	return "", "", "", fmt.Errorf(
		"found candidate directories for %q under %s (%s) but none contained a complete TLS set; expected cert one of [tls.crt vault.crt], key one of [tls.key vault.key], ca one of [ca.crt tls.pem vault.ca]",
		serviceName,
		rootPath,
		strings.Join(candidates, ", "),
	)
}

func firstExistingFile(dirPath string, candidates []string) string {
	for _, name := range candidates {
		p := filepath.Join(dirPath, name)
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p
		}
	}
	return ""
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
