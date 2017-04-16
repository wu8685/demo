package config

type Config struct {
	Port int
}

func New() *Config {
	return &Config{
		Port: 8080,
	}
}
