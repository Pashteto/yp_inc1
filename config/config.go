package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	ServAddr string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	//FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"${HOME}" envExpand:"true"`
	FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"${HOME}/go/src/github.com/Pashteto/yp_inc1/tmp/URLs" envExpand:"true"`
}

func (cfg *Config) UpdateByFlags(ServAddr, BaseURL, FStorPath *string) (bool, error) {
	changed := false
	if *BaseURL != "http://localhost:8080" {
		changed = true
		cfg.BaseURL = *BaseURL
	}
	if *ServAddr != ":8080" {
		changed = true
		cfg.ServAddr = *ServAddr
		fjnv := strings.SplitAfter(cfg.ServAddr, ":")
		var base string
		var port string
		if len(fjnv) > 0 {
			port = fjnv[len(fjnv)-1]
		}
		fjnv = strings.SplitAfter(cfg.BaseURL, ":")
		if len(fjnv) > 0 {
			base = strings.Join(fjnv[:len(fjnv)-1], "")
		}
		if base == "" || port == "" {
			return changed, errors.New("SERVER_ADDRESS, BASE_URL flags error")
		}
		cfg.BaseURL = base + port
	}
	if *FStorPath != os.Getenv("HOME") {
		changed = true
		cfg.FStorPath = *FStorPath
	}
	return changed, nil
}
