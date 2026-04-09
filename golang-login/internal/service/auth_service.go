package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/workshop/golang-login/internal/models"
	"github.com/workshop/golang-login/internal/repository"
)

var (
	ErrUserNotFound  = errors.New("usuário não encontrado")
	ErrWrongPassword = errors.New("senha incorreta")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return &models.LoginResponse{Success: false, Message: "❌ Usuário não encontrado"}, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("erro interno: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return &models.LoginResponse{Success: false, Message: "❌ Senha incorreta"}, ErrWrongPassword
	}

	return &models.LoginResponse{
		Success: true,
		Message: "✅ Login realizado!",
		User:    &models.UserPublic{ID: user.ID, Email: user.Email, Name: user.Name},
	}, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password, name string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	return s.userRepo.Create(ctx, &models.User{Email: email, PasswordHash: string(hash), Name: name})
}

func (s *AuthService) ListUsers(ctx context.Context) ([]models.UserPublic, error) {
	users, err := s.userRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var result []models.UserPublic
	for _, u := range users {
		result = append(result, models.UserPublic{ID: u.ID, Email: u.Email, Name: u.Name})
	}
	return result, nil
}
