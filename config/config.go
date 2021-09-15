package config

import (
	"os"
)

type Config struct {
	ServAddr string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	//FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"${HOME}" envExpand:"true"`
	FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"${HOME}/go/src/github.com/Pashteto/yp_inc1/filed_history//" envExpand:"true"`
}

func (cfg *Config) UpdateByFlags(ServAddr, BaseURL, FStorPath, RedisPtr *string) (changed bool) {
	if *ServAddr != ":8080" {
		cfg.ServAddr = *ServAddr
		changed = true
	}
	if *BaseURL != "http://localhost:8080" {
		changed = true
		cfg.BaseURL = *ServAddr
	}
	if *FStorPath != os.Getenv("HOME") {
		changed = true
		cfg.FStorPath = *FStorPath
	}
	if *RedisPtr != os.Getenv("REDIS_HOST") {
		os.Setenv("REDIS_HOST", *RedisPtr)
	}
	return
}
