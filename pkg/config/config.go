package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Discord discordMeta  `yaml:"discord"`
	DB      databaseMeta `yaml:"db"`
	Redis   redisMeta    `yaml:"redis"`
	UniteDB uniteDBMeta  `yaml:"unite_db"`
}

func LoadConfig() *AppConfig {
	cfg := new(AppConfig)
	b, err := os.ReadFile("setting/setting.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		panic(err)
	} else {
		return cfg
	}
}

type redisMeta struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type databaseMeta struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Minpool  int    `yaml:"minpool"`
	Maxpool  int    `yaml:"maxpool"`
}

type discordMeta struct {
	Mode              string          `json:"mode" yaml:"mode"`
	Logfile           string          `json:"logfile" yaml:"logfile"`
	PermissionInteger uint64          `json:"permission_integer" yaml:"permission_integer"`
	TokenFile         string          `json:"token_file" yaml:"token_file"`
	ClientIdFile      string          `json:"client_id_file" yaml:"client_id_file"`
	ClientSecretFile  string          `json:"client_secret_file" yaml:"client_secret_file"`
	Servers           discordServerID `json:"servers" yaml:"servers"`
}

type discordServerID struct {
	Pokemon string `yaml:"pokepokemon"`
	Dev     string `yaml:"development"`
}

type uniteDBMeta struct {
	BaseURL            string          `yaml:"base_url"`
	UpdateIntervalDays int             `yaml:"update_interval_days"`
	Endpoints          uniteDBEndpoint `yaml:"endpoints"`
}

type uniteDBEndpoint struct {
	Pokemon     string `yaml:"pokemon"`
	HeldItems   string `yaml:"held_items"`
	BattleItems string `yaml:"battle_items"`
}
