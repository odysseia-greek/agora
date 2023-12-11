package config

import (
	"os"
	"strings"
)

func StringFromEnv(envName, defaultValue string) string {
	var value string
	value = os.Getenv(envName)
	if value == "" {
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
		envPodName = DefaultPodname
	}
	splitPodName := strings.Split(envPodName, "-")
	podName := splitPodName[0]

	return podName
}

func SliceFromEnv(sliceName string) []string {
	slice := os.Getenv(sliceName)
	splitSlice := strings.Split(slice, ";")

	return splitSlice
}
