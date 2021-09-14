package config

type Config struct {
	ServAddr  string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL   string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"${HOME}" envExpand:"true"`
	// FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"${HOME}/go/src/github.com/Pashteto/yp_inc1/filed_history//" envExpand:"true"`

}
