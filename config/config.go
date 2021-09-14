package config

//	"net/url"

type Config struct {
	// Host   string `config:"SERVER_HOST"`
	// Port   string `config:"SERVER_PORT"`
	// Scheme string `config:"SERVER_SCHEME"`

	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
	//	return cfg.Scheme + "://" + cfg.Host + ":" + cfg.Port
	// Home         string        `env:"HOME"`
	// Port1        int           `env:"PORT" envDefault:"3000"`
	// Password     string        `env:"PASSWORD,unset"`
	// IsProduction bool          `env:"PRODUCTION"`
	// Hosts        []string      `env:"HOSTS" envSeparator:":"`
	// Duration     time.Duration `env:"DURATION"`
	// TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`

	// User string `env:"USER"`
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
