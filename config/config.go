package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mode         string   `yaml:"mode"`
	Backends     []string `yaml:"backends"`
	Strategy     string   `yaml:"strategy"` // round_robin, random, least_conn
	AllowedCIDRs []string `yaml:"allowed_cidrs"`
	RateLimit    int      `yaml:"rate_limit"`
	Burst        int      `yaml:"burst"`
	BlockedPaths []string `yaml:"blocked_paths"`
	TLS          bool     `yaml:"tls"`
	TLSCertPath  string   `yaml:"tls_cert"`
	TLSKeyPath   string   `yaml:"tls_key"`
}

var AppConfig *Config

func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	AppConfig = &cfg
	log.Printf("✅ Loaded config from %s", path)
}
