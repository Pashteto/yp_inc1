package config

import (
	"fmt"
	"log"
	"os"
)

//	"net/url"

type Config struct {
	// Host   string `config:"SERVER_HOST"`
	// Port   string `config:"SERVER_PORT"`
	// Scheme string `config:"SERVER_SCHEME"`
	Host string `env:"SERVER_HOST"`
	Port string `env:"SERVER_PORT"`
	// Scheme string `config:"SERVER_SCHEME"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
	//	return cfg.Scheme + "://" + cfg.Host + ":" + cfg.Port
	// Home         string        `env:"HOME"`
	Port1 int `env:"PORT" envDefault:"8080"`
	// Password     string        `env:"PASSWORD,unset"`
	// IsProduction bool          `env:"PRODUCTION"`
	Hosts []string `env:"HOSTS" envSeparator:":"`
	// Duration     time.Duration `env:"DURATION"`
	// TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`

	// User string `env:"USER"`
}

func (cfg *Config) CheckEnv() {

	if len(cfg.Hosts) > 0 && cfg.Port1 != 0 {
		cfg.ServerAddress = "http://" + cfg.Hosts[0] + ":" + fmt.Sprint(cfg.Port1)
		return
	}
	if cfg.Port1 != 8080 {
		cfg.ServerAddress = "http://localhost:" + fmt.Sprint(cfg.Port1)
		return
	}

	if os.Getenv("SERVER_HOST") != "" && os.Getenv("SERVER_PORT") != "" {
		cfg.ServerAddress = "http://" + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
		log.Println("here is\t", cfg.ServerAddress)
		return
	}
	if os.Getenv("SERVER_PORT") != "" {
		cfg.ServerAddress = "http://localhost:" + os.Getenv("SERVER_PORT")
		return
	}
	if cfg.Host != "" && cfg.Port != "" {
		cfg.ServerAddress = "http://" + cfg.Host + ":" + cfg.Port
		return
	}
	if cfg.Port != "" {
		cfg.ServerAddress = "http://localhost:" + cfg.Port
		return
	}
}

/*
func (cfg *Config) RecieveEnv(envHost, envPort, envURL string) error {
	if envHost != "" {
		cfg.Host = envHost
	}
	if envPort != "" {
		cfg.Port = envPort
	}
	if envURL != "" {
		envURLParsed, err := url.Parse(envURL)
		if err != nil {
			return err
		}
		cfg.Scheme = envURLParsed.Scheme
		cfg.Host = envURLParsed.Hostname()
		cfg.Port = envURLParsed.Path
	}
	return nil
}

func String(cfg *Config) string {
	return cfg.Scheme + "://" + cfg.Host + ":" + cfg.Port
}
*/
