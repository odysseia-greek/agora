package config

import (
	"log"
	"os"
	"strings"
)

func StringFromEnv(envName, defaultValue string) string {
	var value string
	value = os.Getenv(envName)
	if value == "" {
		log.Printf("%s empty set as env variable - defaulting to %s", envName, defaultValue)
		value = defaultValue
	}

	return value
}

func BoolFromEnv(envName string) bool {
	var value bool
	envValue := os.Getenv(envName)
	if envValue == "" || envValue == "no" || envValue == "false" {
		value = false
	} else {
		value = true
	}

	return value
}

func ParsedPodNameFromEnv() string {
	envPodName := os.Getenv(EnvPodName)
	if envPodName == "" {
		log.Printf("%s empty set as env variable - defaulting to %s", EnvPodName, DefaultPodname)
		envPodName = DefaultPodname
	}
	splitPodName := strings.Split(envPodName, "-")
	podName := splitPodName[0]

	return podName
}

func SliceFromEnv(sliceName string) []string {
	slice := os.Getenv(sliceName)

	if slice == "" {
		log.Print("ELASTIC_ROLES or ELASTIC_INDEXES env variables not set!")
	}

	splitSlice := strings.Split(slice, ";")

	return splitSlice
}
