// Package config centraliza todas as configurações da aplicação.
// Isso evita "magic strings" espalhadas pelo código.
package config

// Config armazena todas as configurações necessárias
type Config struct {
	DatabaseURL string
	ServerPort  string
}

// Load retorna as configurações da aplicação
// Em produção, isso viria de variáveis de ambiente!
func Load() *Config {
	return &Config{
		// Formato: postgres://usuario:senha@host:porta/database
		DatabaseURL: "postgres://workshop:workshop123@localhost:5432/login_system?sslmode=disable",
		ServerPort:  ":8080",
	}
}
