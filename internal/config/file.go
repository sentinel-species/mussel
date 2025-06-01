package config

import (
	"gopkg.in/yaml.v3"
	"mussel/internal/types"
	"os"
)

var (
	Config = &types.Config{}
)

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, Config)
}
