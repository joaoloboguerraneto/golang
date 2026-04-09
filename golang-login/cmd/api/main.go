package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/workshop/golang-login/internal/config"
	"github.com/workshop/golang-login/internal/database"
	"github.com/workshop/golang-login/internal/handler"
	"github.com/workshop/golang-login/internal/repository"
	"github.com/workshop/golang-login/internal/service"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/workshop/golang-login/docs"
)

func main() {
	cfg := config.Load()
	fmt.Println("Iniciando API...")

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Falha ao conectar: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("GET /api/users", authHandler.ListUsers)

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Servir o JSON do Swagger manualmente
	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(swaggerJSON))
	})

	// Arquivos estaticos (UI de login)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	fmt.Println("Conectado ao banco!")
	fmt.Println("")
	fmt.Println("UI Login:  http://localhost:8080")
	fmt.Println("Swagger:   http://localhost:8080/swagger/index.html")
	fmt.Println("")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

const swaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "Golang Login API",
    "description": "API de autenticacao - Workshop Go",
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "schemes": ["http"],
  "paths": {
    "/api/register": {
      "post": {
        "tags": ["Auth"],
        "summary": "Cadastrar usuario",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [{
          "in": "body",
          "name": "body",
          "required": true,
          "schema": {
            "type": "object",
            "required": ["email", "password", "name"],
            "properties": {
              "email": {"type": "string", "example": "maria@email.com"},
              "password": {"type": "string", "example": "123456"},
              "name": {"type": "string", "example": "Maria Santos"}
            }
          }
        }],
        "responses": {
          "201": {"description": "Usuario criado"},
          "400": {"description": "Dados invalidos"}
        }
      }
    },
    "/api/login": {
      "post": {
        "tags": ["Auth"],
        "summary": "Fazer login",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [{
          "in": "body",
          "name": "body",
          "required": true,
          "schema": {
            "type": "object",
            "required": ["email", "password"],
            "properties": {
              "email": {"type": "string", "example": "maria@email.com"},
              "password": {"type": "string", "example": "123456"}
            }
          }
        }],
        "responses": {
          "200": {"description": "Login OK"},
          "401": {"description": "Senha incorreta"},
          "404": {"description": "Usuario nao encontrado"}
        }
      }
    },
    "/api/users": {
      "get": {
        "tags": ["Users"],
        "summary": "Listar usuarios",
        "produces": ["application/json"],
        "responses": {
          "200": {"description": "Lista de usuarios"}
        }
      }
    }
  }
}`
