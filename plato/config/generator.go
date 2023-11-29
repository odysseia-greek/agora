package config

import (
	"crypto/tls"
	"errors"
	"github.com/odysseia-greek/agora/plato/certificates"
	"github.com/odysseia-greek/agora/plato/helpers"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func CreateOdysseiaClient() (service.OdysseiaClient, error) {
	solonUrl := StringFromEnv(EnvSolonService, DefaultServiceAddress)
	herodotosUrl := StringFromEnv(EnvHerodotosService, DefaultServiceAddress)
	alexandrosUrl := StringFromEnv(EnvAlexandrosService, DefaultServiceAddress)
	dionysiosUrl := StringFromEnv(EnvDionysiosService, DefaultServiceAddress)
	sokratesUrl := StringFromEnv(EnvSokratesService, DefaultServiceAddress)
	tlsEnabled := BoolFromEnv(EnvTlSKey)

	solonParsed, err := url.Parse(solonUrl)
	if err != nil {
		return nil, err
	}

	herodotosParsed, err := url.Parse(herodotosUrl)
	if err != nil {
		return nil, err
	}
	alexandrosParsed, err := url.Parse(alexandrosUrl)
	if err != nil {
		return nil, err
	}
	dionysiosParsed, err := url.Parse(dionysiosUrl)
	if err != nil {
		return nil, err
	}
	sokratesParsed, err := url.Parse(sokratesUrl)
	if err != nil {
		return nil, err
	}

	config := service.ClientConfig{
		Ca: nil,
		Solon: service.OdysseiaApi{
			Url:    solonParsed.Host,
			Scheme: solonParsed.Scheme,
			Cert:   nil,
		},
		Herodotos: service.OdysseiaApi{
			Url:    herodotosParsed.Host,
			Scheme: herodotosParsed.Scheme,
			Cert:   nil,
		},
		Dionysios: service.OdysseiaApi{
			Url:    dionysiosParsed.Host,
			Scheme: dionysiosParsed.Scheme,
			Cert:   nil,
		},
		Alexandros: service.OdysseiaApi{
			Url:    alexandrosParsed.Host,
			Scheme: alexandrosParsed.Scheme,
			Cert:   nil,
		},
		Sokrates: service.OdysseiaApi{
			Url:    sokratesParsed.Host,
			Scheme: sokratesParsed.Scheme,
			Cert:   nil,
		},
	}

	if tlsEnabled {
		log.Print("setting up certs because TLS is enabled")
		rootPath := os.Getenv("CERT_ROOT")
		log.Printf("rootPath: %s", rootPath)
		dirs, err := ioutil.ReadDir(rootPath)
		if err != nil {
			return nil, err
		}

		for _, dir := range dirs {
			if dir.IsDir() {
				dirPath := filepath.Join(rootPath, dir.Name())
				log.Printf("found directory: %s", dirPath)

				certPath := filepath.Join(dirPath, "tls.crt")
				keyPath := filepath.Join(dirPath, "tls.key")

				if _, err := os.Stat(certPath); errors.Is(err, os.ErrNotExist) {
					log.Printf("cannot get file because it does not exist: %s", certPath)
					continue
				}

				if _, err := os.Stat(keyPath); errors.Is(err, os.ErrNotExist) {
					log.Printf("cannot get file because it does not exist: %s", keyPath)
					continue
				}

				loadedCerts, err := tls.LoadX509KeyPair(certPath, keyPath)
				if err != nil {
					return nil, err
				}

				if config.Ca == nil {
					caPath := filepath.Join(rootPath, dir.Name(), "tls.pem")
					if _, err := os.Stat(caPath); errors.Is(err, os.ErrNotExist) {
						log.Printf("cannot get file because it does not exist: %s", caPath)
						continue
					}
					config.Ca, _ = ioutil.ReadFile(caPath)
					log.Printf("writing CA for path %s", caPath)
				}

				switch dir.Name() {
				case "solon":
					config.Solon.Cert = []tls.Certificate{loadedCerts}
				case "dionysios":
					config.Dionysios.Cert = []tls.Certificate{loadedCerts}
				case "herodotos":
					config.Herodotos.Cert = []tls.Certificate{loadedCerts}
				case "alexandros":
					config.Alexandros.Cert = []tls.Certificate{loadedCerts}
				case "sokrates":
					config.Sokrates.Cert = []tls.Certificate{loadedCerts}
				}
			}
		}
	}

	return service.NewClient(config)
}

func RetrieveCertPathLocally(testOverwrite bool, service string) (cert string, key string) {
	keyName := "tls.key"
	certName := "tls.crt"

	if testOverwrite {
		log.Print("trying to read cert file from file")
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

		log.Printf("found certpath: %s - found keypath: %s", cert, key)
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

func PlatoPath(path string) string {
	dirs, _ := ioutil.ReadDir(path)
	var currentModTime time.Time
	var latestPath string
	for _, dir := range dirs {

		if dir.Name() == PLATO {
			latestPath = dir.Name()
			break
		}

		if currentModTime.Before(dir.ModTime()) {
			currentModTime = dir.ModTime()
			latestPath = dir.Name()
		}
	}

	return filepath.Join(path, latestPath)
}

func OdysseiaRootPath(path string) string {
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var odysseiaPath []string
	for i, part := range dirParts {
		if part == path {
			odysseiaPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range odysseiaPath {
		l = filepath.Join(l, path)
	}

	return l
}
