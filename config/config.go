package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host   string `config:"SERVER_HOST"`
	Port   string `config:"SERVER_PORT"`
	Scheme string `config:"SERVER_SCHEME"`
}

func ReadFile(cfg *Config) error {
	f, err := os.Open("conf.json")
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}
	return nil
}
func String(cfg *Config) string {
	return cfg.Scheme + "://" + cfg.Host + ":" + cfg.Port
}
