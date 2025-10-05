// Package config handles loading and parsing of IPVS configuration from YAML files
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the complete IPVS configuration
type Config struct {
	Service  ServiceSpec   `yaml:"service"`  // Virtual service definition
	Backends []BackendSpec `yaml:"backends"` // List of backend servers
}

// ServiceSpec defines the virtual service (VIP) configuration
type ServiceSpec struct {
	VIP       string `yaml:"vip"`       // Virtual IP address
	Port      uint16 `yaml:"port"`      // Virtual port
	Protocol  string `yaml:"protocol"`  // Protocol (tcp/udp)
	Scheduler string `yaml:"scheduler"` // Load balancing scheduler (rr, wrr, lc, etc.)
}

// BackendSpec defines a backend server configuration
type BackendSpec struct {
	IP   string `yaml:"ip"`   // Backend server IP address
	Port uint16 `yaml:"port"` // Backend server port
}

// Load reads and parses a YAML configuration file into a Config struct
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
