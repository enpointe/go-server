package server

import (
	"encoding/json"
	"fmt"
	"os"
)

// ConfigFilename the Default configuration filename
const ConfigFilename string = "Configuration.json"

// Config configuration variables that can be set by the operator
type Config struct {
	JWTKey string `json:"jwtKey"`
}

// ReadConfig read the configuration from the specified configuration file
func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	var configuration Config
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Decode Error")
		return nil, err
	}
	return &configuration, nil
}
