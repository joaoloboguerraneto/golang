package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect estabelece conexão com o PostgreSQL
func Connect() error {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "authdb")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	// Tentar conectar com retry (útil para Docker)
	for i := 0; i < 30; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				fmt.Println("Conectado ao banco de dados!")
				return nil
			}
		}
		fmt.Printf("Tentando conectar ao banco... (%d/30)\n", i+1)
		time.Sleep(time.Second)
	}

	return fmt.Errorf("não foi possível conectar ao banco: %v", err)
}

// Close fecha a conexão
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// CreateTables cria as tabelas necessárias
func CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		birth_date DATE NOT NULL,
		failed_attempts INT DEFAULT 0,
		is_blocked BOOLEAN DEFAULT FALSE,
		blocked_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		token VARCHAR(255) UNIQUE NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS captcha_tokens (
		id SERIAL PRIMARY KEY,
		token VARCHAR(255) UNIQUE NOT NULL,
		answer VARCHAR(50) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
	CREATE INDEX IF NOT EXISTS idx_captcha_token ON captcha_tokens(token);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("erro ao criar tabelas: %v", err)
	}

	fmt.Println("Tabelas criadas/verificadas com sucesso!")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
