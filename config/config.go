package config

import "fmt"

type DBConfig struct {
	Host string
	Port string
	Database string
	Username string
	Password string
	Driver string
}

type ApiConfig struct {
	ApiPort string
}

type Config struct {
	DBConfig
	ApiConfig
}

func (c *Config) readConfig() error {
	c.DBConfig = DBConfig{
		Host: "localhost",
		Port: "5432",
		Database: "enigma_laundry_db",
		Username: "postgres",
		Password: "root",
		Driver: "postgres",
	}
	c.ApiConfig = ApiConfig{
		ApiPort: "8080",
	}

	if c.Host == "" || c.Port == "" || c.Username == "" || c.Password == "" || c.ApiPort == "" {
		return fmt.Errorf("required config")
	}
	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cfg.readConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}