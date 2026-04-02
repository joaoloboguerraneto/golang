//Package conter todas as operacoes de banco de dados relacionadas aos usuarios
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang-login/internal/models"
)

// Erros customizados para facilitar o tratamento de erros na camada de serviço
var (
	ErrUserNotFound    = fmt.Errorf("usuario nao encontrado")
)

// UserRepository define a interface para as operacoes de banco de dados relacionadas aos usuarios
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository criar uma nova instancia do repositorio de usuarios com a conexao ao banco de dados fornecida
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

//Buscar o usuario pelo o email, retornando o usuario ou um erro caso nao seja encontrado
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error){
	//Query SQL para buscar o usuario pelo email (SQL Injection seguro usando parametros)
	query := `
	SELECT id, email, password_hash, created_at 
	FROM users 
	WHERE email = $1
	`
	//$1 e um placeholder para o parametro email, que sera passado de forma segura para evitar SQL Injection

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, 
		&user.Email, 
		&user.PasswordHash, 
		&user.CreatedAt,
	)
	
	//Tratar o erro de usuario nao encontrado (pgx.ErrNoRows) e retornar um erro customizado para a camada de serviço
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuario por email: %w", err)
	}

	return &user, nil	
}

// Criar um novo usuario no banco de dados, retornando o usuario criado ou um erro caso a operacao falhe
func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	//Query SQL para inserir um novo usuario (SQL Injection seguro usando parametros)
	query := `
	INSERT INTO users (email, password_hash, created_at) 
	VALUES ($1, $2, $3) 
	RETURNING id
	`
	//$1, $2 e $3 sao placeholders para os parametros email, password_hash e created_at, que serao passados de forma segura para evitar SQL Injection

	err := r.db.QueryRow(
		ctx, 
		query, 
		user.Email, 
		user.PasswordHash, 
		user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar usuario: %w", err)
	}

	return user, nil
}