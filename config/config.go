package config

type Config struct {
	ServAddr  string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL   string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FStorPath string `env:"FILE_STORAGE_PATH" envDefault:"../URLs" envExpand:"true"`

	PostgresURL string `env:"DATABASE_URL" envDefault:"host=localhost port=5432 user=postgres password=kornkorn dbname=mydb sslmode=disable" envExpand:"true"`
	// PostgresURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"`
}

func (cfg *Config) UpdateByFlags(ServAddr, BaseURL, FStorPath, PostgresURL *string) (bool, error) {
	changed := false
	if *BaseURL != "http://localhost:8080" {
		changed = true
		cfg.BaseURL = *BaseURL
	}
	if *ServAddr != ":8080" {
		changed = true
		cfg.ServAddr = *ServAddr
	}
	if *FStorPath != "../URLs" {
		changed = true
		cfg.FStorPath = *FStorPath
	}
	if *PostgresURL != "host=localhost port=5432 user=postgres password=kornkorn dbname=mydb sslmode=disable" {
		// if *PostgresURL != "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable" {
		cfg.PostgresURL = *PostgresURL
	}
	return changed, nil
}
