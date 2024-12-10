package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	URL  string `yaml:"url"`
	File string `yaml:"file"`
}

type Config struct {
	ResponsesDir string                `yaml:"responses_dir"`
	Endpoints    map[string][]Endpoint `yaml:"endpoints"`
}

func MustLoadConfig() Config {
	//  Get the config file from the command line arguments
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <YAML-config-file>", os.Args[0])
	}
	configPath := os.Args[1]

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	// Decode YAML
	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	return config
}
