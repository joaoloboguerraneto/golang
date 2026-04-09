# Golang

Sistema de autenticaçao em Golang

## Pré-requisitos

- Go 1.21
- Docker Desktop
- Terminal

### Quick Start

```bash
# 1. Dar permissão ao script
chmod +x setup.sh
 
# 2. Rodar o setup (cria banco, tabelas, baixa deps)
./setup.sh
 
# 3. Executar a aplicação
cd golang-login
go run cmd/api/main.go
```
## Estrutura do Projeto
 
```
golang-login/
├── cmd/api/main.go                 # Ponto de entrada
├── internal/       
│   ├── config/config.go            # Configurações
│   ├── database/postgres.go        # Conexão DB
│   ├── models/user.go              # Estruturas de dados
│   ├── repository/user_repo.go     # Queries
│   └── service/auth_service.go     # Lógica de negócio
├── go.mod
└── go.sum
```

## Cenários de Teste
 
O programa testa automaticamente:
 
| Cenário           | Email                  | Senha            | Resultado       |
|-------------------|------------------------|------------------|-----------------|
| User não existe   | naoexiste@email.com    | qualquer         | Não encontrado  |
| Login OK          | joao@email.com         | senha123         | Sucesso         |
| Senha errada      | joao@email.com         | senhaerrada      | Senha incorreta |

 
## Comandos Úteis
 
```bash
# Conectar ao PostgreSQL
docker exec -it postgres-workshop psql -U workshop -d login_system
 
# Ver usuários cadastrados
SELECT id, email, name, created_at FROM users;
 
# Parar o banco
docker stop postgres-workshop
 
# Reiniciar o banco
docker start postgres-workshop
 
# Remover o container (perde dados!)
docker rm -f postgres-workshop
```
 
## Conceitos Abordados
 
- [x] Estrutura de projetos Go
- [x] Conexão com PostgreSQL (pgx)
- [x] Repository Pattern
- [x] Service Layer
- [x] Hash de senhas (bcrypt)
- [x] Error handling
- [x] Context e timeouts
- [x] Dependency Injection
 
## Próximos Passos
 
1. Adicionar API HTTP com `net/http`
2. Implementar JWT tokens
3. Adicionar middleware de logging
4. Criar testes unitários
5. Dockerizar a aplicação
 
---
 
**Workshop Go**
