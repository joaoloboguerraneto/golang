package models 

import "time"

// Usuario representa a estrutura de um usuario no sistema
// As tags 'db' indicam como os campos devem ser mapeados para o banco de dados

type User struct {
	ID        	  int       `db:"id"`           	// Identificador unico do usuario
	Email	  	  string    `db:"email"`   		   	// Email unico
	PasswordHash  string    `db:"password_hash"`   	// Senha do usuario (deve ser armazenada de forma segura, ex: hash)
	CreatedAt time.Time `db:"created_at"`        	// Data e hora de criacao do usuario
}

// LoginRequest representa a estrutura dos dados recebidos na requisicao de login
type LoginRequest struct {
	Email    string `json:"email"`    // Email do usuario
	Password string `json:"password"` // Senha do usuario
}

// UserPublic representa a estrutura dos dados do usuario que serao expostos na resposta da API
type UserPublic struct {
	ID    int    `json:"id"`    // Identificador unico do usuario
	Email string `json:"email"` // Email do usuario
}
