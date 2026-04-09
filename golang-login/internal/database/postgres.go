// Package database gerencia a conexão com o PostgreSQL
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect estabelece conexão com o PostgreSQL
// Retorna um pool de conexões (reutilizáveis e thread-safe)
func Connect(databaseURL string) (*pgxpool.Pool, error) {
	// Context com timeout - não queremos esperar para sempre!
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Libera recursos quando a função terminar

	// pgxpool mantém várias conexões abertas (pool)
	// Isso é MUITO mais eficiente que abrir/fechar a cada query
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar pool de conexões: %w", err)
	}

	// Testar se a conexão realmente funciona
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("erro ao conectar no banco: %w", err)
	}

	fmt.Println("Conectado ao PostgreSQL!")
	return pool, nil
}
