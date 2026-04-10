package main

import (
	"log"
	"net/http"
	"os"

	"auth-system/database"
	"auth-system/handlers"
)

func main() {
	// Conectar ao banco de dados
	if err := database.Connect(); err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer database.Close()

	// Criar tabelas se não existirem
	if err := database.CreateTables(); err != nil {
		log.Fatal("Erro ao criar tabelas:", err)
	}

	// Servir arquivos estáticos (CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Rotas
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/reset-password", handlers.ResetPasswordHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/captcha", handlers.CaptchaHandler)
	http.HandleFunc("/verify-captcha", handlers.VerifyCaptchaHandler)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando em http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
