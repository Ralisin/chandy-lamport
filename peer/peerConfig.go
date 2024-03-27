package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var serviceRegistryAddr string
var serviceRegistryPort string
var peerAddr string

// Config Configuration struct to represent the JSON data
type Config struct {
	Localhost LocalhostConfig `json:"localhost"`
	Docker    DockerConfig    `json:"docker"`
}

type LocalhostConfig struct {
	ServiceRegistryAddr string `json:"serviceRegistryAddr"`
	ServiceRegistryPort string `json:"serviceRegistryPort"`
	PeerAddr            string `json:"peerAddr"`
}

type DockerConfig struct {
	ServiceRegistryAddr string `json:"serviceRegistryAddr"`
	ServiceRegistryPort string `json:"serviceRegistryPort"`
	PeerAddr            string `json:"peerAddr"`
}

func appArgsFetch() error {
	// Call the function to read the configuration
	config, err := ReadConfig("../config.json")
	if err != nil {
		return fmt.Errorf("Error reading configuration: %v\n", err)
	}

	if len(os.Args) == 1 || (len(os.Args) == 2 && os.Args[1] == "-localhost") {
		serviceRegistryAddr = config.Localhost.ServiceRegistryAddr
		serviceRegistryPort = config.Localhost.ServiceRegistryPort

		peerAddr = config.Localhost.PeerAddr
	} else if len(os.Args) == 2 && os.Args[1] == "-docker" {
		serviceRegistryAddr = config.Docker.ServiceRegistryAddr
		serviceRegistryPort = config.Docker.ServiceRegistryPort

		peerAddr = config.Docker.PeerAddr
	} else {
		return fmt.Errorf("\nUsage: go run . [-localhost/-docker]\n\nDefault flag: -localhost")
	}

	return nil
}

func ReadConfig(filename string) (Config, error) {
	var config Config

	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("error reading JSON file: %v", err)
	}

	// Unmarshal JSON data into Config struct
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return config, nil
}
