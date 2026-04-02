package config

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() *Config {
	return &Config{
		DatabaseURL: "postgres://${username}:${password}@${host}:${port}/${database}?sslmode=disable",
		Port:        "8080",
	}
}