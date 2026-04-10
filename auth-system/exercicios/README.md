# Exercicios Praticos

Esta pasta contem exercicios para praticar os conceitos do sistema de autenticacao.

## Como executar

Cada exercicio pode ser executado individualmente:

```bash
cd exercicios

# Exercicio 1: Structs
go run 01_structs.go

# Exercicio 2: Ponteiros
go run 02_ponteiros.go

# Exercicio 3: Hash de Senhas (precisa instalar biblioteca)
go get golang.org/x/crypto/bcrypt
go run 03_hash_senhas.go

# Exercicio 4: Sistema de Bloqueio (interativo)
go run 04_bloqueio.go

# Exercicio 5: Captcha (interativo)
go run 05_captcha.go

# Exercicio 6: HTTP Handlers (servidor web)
go run 06_http_handlers.go
# Depois acesse: http://localhost:8081
```

## Ordem recomendada

1. **01_structs.go** - Entenda como criar "fichas de dados"
2. **02_ponteiros.go** - Entenda referencias e enderecos de memoria
3. **03_hash_senhas.go** - Entenda por que e como criptografar senhas
4. **04_bloqueio.go** - Simule o sistema de bloqueio de conta
5. **05_captcha.go** - Simule o sistema anti-bot
6. **06_http_handlers.go** - Entenda como funcionam rotas web

## Relacao com o projeto

| Exercicio | Arquivo do projeto | Conceitos |
|-----------|-------------------|-----------|
| 01_structs | models/user.go | type User struct |
| 02_ponteiros | models/user.go | *User, &usuario |
| 03_hash_senhas | models/user.go | bcrypt.GenerateFromPassword |
| 04_bloqueio | models/user.go | Authenticate, BlockUser |
| 05_captcha | models/captcha.go | GenerateCaptcha, ValidateCaptcha |
| 06_http_handlers | handlers/auth.go | LoginHandler, RegisterHandler |

## Dicas

- Leia os comentarios nos arquivos - eles explicam cada linha
- Modifique o codigo e veja o que acontece
- Tente resolver os desafios no final de cada exercicio
- Use `fmt.Println()` para "debugar" e entender o fluxo
