package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Bot   BotMeta   `yaml:"bot"`
	API   APIMeta   `yaml:"api"`
	Redis RedisMeta `yaml:"redis"`
}

type BotMeta struct {
	Mode             string            `yaml:"mode"`
	TokenFile        string            `yaml:"token_file"`
	ClientIDFile     string            `yaml:"client_id_file"`
	ClientSecretFile string            `yaml:"client_secret_file"`
	PermissionFlag   int               `yaml:"permission_flag"`
	ServerID         map[string]string `yaml:"server_id"`
}

// APIMeta is the upstream gRPC API config
type APIMeta struct {
	Host           string `yaml:"host"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

// RedisMeta is the redis configuration
type RedisMeta struct {
	Host string `yaml:"host"`
}

// MustLoadConfig panics on error
func MustLoadConfig(path string) *AppConfig {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	cfg := new(AppConfig)
	err = yaml.UnmarshalStrict(b, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
