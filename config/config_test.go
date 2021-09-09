package config

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadFile(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_RecieveEnv(t *testing.T) {
	type fields struct {
		Host   string
		Port   string
		Scheme string
	}
	type args struct {
		envHost string
		envPort string
		envURL  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Host:   tt.fields.Host,
				Port:   tt.fields.Port,
				Scheme: tt.fields.Scheme,
			}
			if err := cfg.RecieveEnv(tt.args.envHost, tt.args.envPort, tt.args.envURL); (err != nil) != tt.wantErr {
				t.Errorf("Config.RecieveEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
