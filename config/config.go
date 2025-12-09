// Package config provides configuration loading utilities for Planx.
package config

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadYAML loads a YAML configuration file into the given struct.
func LoadYAML(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

// LoadJSON loads a JSON configuration file into the given struct.
func LoadJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ParseYAML parses YAML bytes into the given struct.
func ParseYAML(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// ParseJSON parses JSON bytes into the given struct.
func ParseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
