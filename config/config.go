package config

type Config struct {
	SeAd  string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`

	TempFolder string `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`

	//	return cfg.Scheme + "://" + cfg.Host + ":" + cfg.Port
	// Home         string        `env:"HOME"`
	//	Port1 int `env:"PORT" envDefault:"8080"`
	// Password     string        `env:"PASSWORD,unset"`
	// IsProduction bool          `env:"PRODUCTION"`
	//	Hosts []string `env:"HOSTS" envSeparator:":"`
	// Duration     time.Duration `env:"DURATION"`

	// User string `env:"USER"`
}
