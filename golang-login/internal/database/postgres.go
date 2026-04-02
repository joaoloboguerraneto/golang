package database 

import (
	"context"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect estabelecer conexao com o  PostgreSQL usando a URL de conexao fornecida
// Retorna um pool de conexoes ou um erro caso a conexao falhe
func Connect(databaseURL string) (*pgxpool.Pool, error) {
	// Context com timeout para evitar conexoes pendentes - nao querer esperar para sempre!
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //Liberar recursos do contexto apos a conexao ser estabelecida ou falhar

	//pool matem varias conexoes abertas (pool)
	// eficiente para lidar com multiplas requisicoes simultaneas sem precisar abrir uma nova conexao a cada vez
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar pool de conexoes: %w", err)
	}

	// Testar se a conexao com o banco de dados esta funcionando corretamente
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("erro ao testar conexao com o banco de dados: %w", err)
	}

	fmt.Print("Conectado com sucesso no PostgreSQL!")
	return pool, nil
}