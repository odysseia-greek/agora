package config

import "crypto/tls"

type Config struct {
	Addresses []string // Addresses for each replica
	Streaming bool
	SavePath  string
	TLSConfig *tls.Config
}

const TRACING_KEY string = "aischylos"
