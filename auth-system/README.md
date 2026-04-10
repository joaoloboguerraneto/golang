# AuthSystem - Sistema de Autenticacao em Go

Sistema completo de autenticacao com login, registro, bloqueio de conta e reset de senha.

## Funcionalidades

- Login com email e senha
- Registro de novos usuarios com validacao anti-bot (captcha matematico)
- Bloqueio automatico apos 3 tentativas de senha incorreta
- Reset de senha validando data de nascimento
- Sessoes com cookies HTTP-only
- Interface moderna e responsiva

## Fluxo de Autenticacao

```
+------------------+
|      Login       |
+--------+---------+
         |
         v
+------------------+     NAO     +------------------+
| Usuario existe?  +------------>|    Registro      |
+--------+---------+             +------------------+
         | SIM
         v
+------------------+     NAO     +------------------+
|  Senha correta?  +------------>| Incrementa erro  |
+--------+---------+             +--------+---------+
         | SIM                            |
         v                                v
+------------------+             +------------------+
|    Dashboard     |             |  3 tentativas?   |
+------------------+             +--------+---------+
                                          | SIM
                                          v
                                 +------------------+
                                 | Bloqueia conta   |
                                 +--------+---------+
                                          |
                                          v
                                 +------------------+
                                 |  Reset senha     |
                                 | (data nasc.)     |
                                 +------------------+
```

## Estrutura do Projeto

```
auth-system/
├── main.go                 # Ponto de entrada
├── go.mod                  # Dependencias Go
├── Dockerfile              # Dockerfile da aplicacao
├── Dockerfile.postgres     # Dockerfile do banco de dados
├── docker-compose.yml      # Orquestracao dos containers
├── init.sql                # Script de inicializacao do banco
├── .env.example            # Exemplo de variaveis de ambiente
├── .gitignore
├── database/
│   └── database.go         # Conexao com PostgreSQL
├── models/
│   ├── user.go             # Modelo de usuario
│   └── captcha.go          # Modelo de captcha
├── handlers/
│   └── auth.go             # Handlers HTTP
├── templates/
│   ├── login.html          # Pagina de login
│   ├── register.html       # Pagina de registro
│   ├── reset.html          # Pagina de reset de senha
│   └── dashboard.html      # Dashboard do usuario
└── static/
    └── style.css           # Estilos CSS
```

## Como executar

### Com Docker (recomendado)

1. Certifique-se de ter Docker e Docker Compose instalados

2. Clone o repositorio e entre na pasta:
   ```bash
   cd auth-system
   ```

3. Execute com Docker Compose:
   ```bash
   docker-compose up --build
   ```

4. Acesse http://localhost:8080

### Sem Docker (desenvolvimento)

1. Instale Go 1.21 ou superior

2. Instale e configure PostgreSQL

3. Crie o banco de dados:
   ```sql
   CREATE DATABASE authdb;
   ```

4. Configure as variaveis de ambiente:
   ```bash
   cp .env.example .env
   # Edite .env com suas configuracoes
   ```

5. Execute a aplicacao:
   ```bash
   go mod download
   go run main.go
   ```

## Variaveis de Ambiente

| Variavel     | Descricao              | Padrao    |
|--------------|------------------------|-----------|
| PORT         | Porta da aplicacao     | 8080      |
| DB_HOST      | Host do PostgreSQL     | localhost |
| DB_PORT      | Porta do PostgreSQL    | 5432      |
| DB_USER      | Usuario do banco       | postgres  |
| DB_PASSWORD  | Senha do banco         | postgres  |
| DB_NAME      | Nome do banco          | authdb    |

## Seguranca

- Senhas armazenadas com bcrypt (hash + salt)
- Sessoes com tokens aleatorios de 256 bits
- Cookies HTTP-only para prevenir XSS
- Captcha matematico para prevenir bots
- Bloqueio de conta apos tentativas falhas
- Validacao de data de nascimento para reset

## Tecnologias

- Go 1.21
- PostgreSQL 15
- Docker & Docker Compose
- HTML/CSS (sem frameworks)
- bcrypt para hash de senhas

## Comandos Docker uteis

```bash
# Iniciar containers
docker-compose up -d

# Parar containers
docker-compose down

# Ver logs
docker-compose logs -f

# Reiniciar apenas a aplicacao
docker-compose restart app

# Limpar tudo (incluindo dados)
docker-compose down -v

# Rebuild apos mudancas
docker-compose up --build
```

## Licenca

MIT License
