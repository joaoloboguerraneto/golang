// Package models define as estruturas de dados da aplicação
package models

import "time"

// User representa um usuário no sistema
// As tags `db` mapeiam para as colunas do banco
type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"` // NUNCA armazenar senha em texto!
	Name         string    `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
}

// LoginRequest representa os dados enviados no login
// Separamos para não expor campos internos como ID ou CreatedAt
type LoginRequest struct {
	Email    string
	Password string // Senha em texto (só existe na memória, temporariamente)
}

// LoginResponse é o que retornamos após login bem-sucedido
type LoginResponse struct {
	Success bool
	Message string
	User    *UserPublic // Dados públicos do usuário
}

// UserPublic contém apenas dados seguros para expor
// NUNCA incluir password_hash aqui!
type UserPublic struct {
	ID    int
	Email string
	Name  string
}
