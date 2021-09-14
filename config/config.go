package config

type Config struct {
	ServAddr string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`

	//	FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"${HOME}/filed_history" envExpand:"true"`
	FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"${HOME}/go/src/github.com/Pashteto/yp_inc1/filed_history//" envExpand:"true"`

	//	return cfg.Scheme + "://" + cfg.Host + ":" + cfg.Port
	// Home         string        `env:"HOME"`
	//	Port1 int `env:"PORT" envDefault:"8080"`
	// Password     string        `env:"PASSWORD,unset"`
	// IsProduction bool          `env:"PRODUCTION"`
	//	Hosts []string `env:"HOSTS" envSeparator:":"`
	// Duration     time.Duration `env:"DURATION"`

	// User string `env:"USER"`
}
