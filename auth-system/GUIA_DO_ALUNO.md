# GUIA DO ALUNO - Sistema de Autenticacao em Go

## Indice

1. [Introducao: O que vamos construir?](#1-introducao)
2. [Conceitos Basicos](#2-conceitos-basicos)
3. [Entendendo o user.go](#3-entendendo-o-usergo)
4. [Entendendo o captcha.go](#4-entendendo-o-captchago)
5. [Entendendo o auth.go](#5-entendendo-o-authgo)
6. [Exercicios Praticos](#6-exercicios-praticos)
7. [Testes no Terminal](#7-testes-no-terminal)

---

## 1. Introducao

### O que e um Sistema de Autenticacao?

Pense em um clube exclusivo:
- Para ENTRAR, voce precisa mostrar sua carteirinha (LOGIN)
- Para ter carteirinha, precisa se CADASTRAR primeiro (REGISTRO)
- Se esquecer a carteirinha, precisa provar quem voce e (RESET SENHA)
- Se tentar entrar com carteirinha falsa varias vezes, e BLOQUEADO

Nosso sistema faz exatamente isso, mas para websites!

### Arquivos que vamos estudar

```
models/
  ├── user.go      -> Tudo sobre usuarios (criar, buscar, logar)
  └── captcha.go   -> Sistema anti-robo (prova que e humano)

handlers/
  └── auth.go      -> Recebe requisicoes e responde (tipo um garcom)
```

---

## 2. Conceitos Basicos

### 2.1 O que e uma Struct?

Uma struct e como uma FICHA DE CADASTRO. Define quais informacoes guardar.

```
FICHA DE ALUNO (no papel)          STRUCT EM GO
┌─────────────────────────┐        type Aluno struct {
│ Nome: ______________    │            Nome     string
│ Idade: _____________    │            Idade    int
│ Email: _____________    │            Email    string
│ Ativo: [ ] Sim [ ] Nao  │            Ativo    bool
└─────────────────────────┘        }
```

**Teste no Go Playground (https://go.dev/play/):**

```go
package main

import "fmt"

// Definindo a struct (a ficha)
type Aluno struct {
    Nome  string
    Idade int
    Email string
}

func main() {
    // Criando um aluno (preenchendo a ficha)
    joao := Aluno{
        Nome:  "Joao Silva",
        Idade: 25,
        Email: "joao@email.com",
    }
    
    // Acessando os dados
    fmt.Println("Nome:", joao.Nome)
    fmt.Println("Idade:", joao.Idade)
    fmt.Println("Email:", joao.Email)
}
```

### 2.2 O que e um Ponteiro?

Ponteiro e um ENDERECO. Em vez de copiar o dado, voce aponta para onde ele esta.

```
SEM PONTEIRO (copia):              COM PONTEIRO (referencia):
┌─────────┐    ┌─────────┐         ┌─────────┐
│ Joao    │───>│ Joao    │         │ Joao    │<──────┐
│ (copia) │    │(original)│        │(original)│       │
└─────────┘    └─────────┘         └─────────┘       │
                                                     │
                                   ┌─────────┐       │
                                   │ &Joao   │───────┘
                                   │(endereco)│
                                   └─────────┘
```

**Teste pratico:**

```go
package main

import "fmt"

func main() {
    nome := "Maria"
    
    // & = "me da o endereco de"
    endereco := &nome
    
    fmt.Println("Valor:", nome)           // Maria
    fmt.Println("Endereco:", endereco)    // 0xc0000... (endereco de memoria)
    
    // * = "vai no endereco e pega/muda o valor"
    *endereco = "Ana"
    
    fmt.Println("Novo valor:", nome)      // Ana (mudou o original!)
}
```

### 2.3 O que e um Erro (error)?

Em Go, funcoes podem retornar DOIS valores: o resultado E um possivel erro.

```
FUNCAO QUE PODE FALHAR:

resultado, erro := FazerAlgo()

if erro != nil {
    // Algo deu errado!
    fmt.Println("Erro:", erro)
} else {
    // Tudo certo!
    fmt.Println("Resultado:", resultado)
}
```

**Analogia do Delivery:**

```
Voce pede uma pizza (chama funcao)

CENARIO 1 - Sucesso:
  pizza = Pizza de Calabresa
  erro  = nil (nada)
  
CENARIO 2 - Falha:
  pizza = nil (nada)
  erro  = "Restaurante fechado"
```

---

## 3. Entendendo o user.go

### 3.1 A Struct User (Ficha do Usuario)

```go
type User struct {
    ID             int        // Numero unico (1, 2, 3...)
    Email          string     // "joao@email.com"
    Password       string     // Senha CRIPTOGRAFADA (nunca texto puro!)
    Name           string     // "Joao Silva"
    BirthDate      time.Time  // Data de nascimento
    FailedAttempts int        // Quantas vezes errou a senha
    IsBlocked      bool       // Esta bloqueado? true/false
    BlockedAt      *time.Time // Quando foi bloqueado (pode ser vazio)
    CreatedAt      time.Time  // Quando criou a conta
    UpdatedAt      time.Time  // Ultima atualizacao
}
```

**Visualizando um usuario:**

```
┌────────────────────────────────────────────┐
│                 USUARIO                     │
├────────────────────────────────────────────┤
│ ID:             42                          │
│ Email:          joao@email.com              │
│ Password:       $2a$10$xK8f... (HASH)       │
│ Name:           Joao Silva                  │
│ BirthDate:      1990-05-15                  │
│ FailedAttempts: 0                           │
│ IsBlocked:      false                       │
│ BlockedAt:      nil (vazio)                 │
│ CreatedAt:      2024-01-10 10:30:00         │
│ UpdatedAt:      2024-01-10 10:30:00         │
└────────────────────────────────────────────┘
```

### 3.2 Por que NUNCA salvar senha em texto puro?

**ERRADO (perigo!):**
```
Banco de dados:
| email            | senha      |
|------------------|------------|
| joao@email.com   | minhasenha |  <- QUALQUER UM QUE ACESSAR O BANCO VE!
```

**CORRETO (com hash):**
```
Banco de dados:
| email            | senha                                          |
|------------------|------------------------------------------------|
| joao@email.com   | $2a$10$N9qo8uLOickgx2ZMRZoMy.somethinglong...   |
```

**O que e HASH?**

Hash e uma funcao UNIDIRECIONAL. Transforma texto em codigo, mas nao da pra voltar.

```
ENTRADA          FUNCAO HASH         SAIDA
"senha123"   -->  bcrypt()    -->  "$2a$10$xK8f..."

NAO EXISTE FUNCAO PARA VOLTAR!
"$2a$10$xK8f..."  -->  ???     -->  "senha123"  (IMPOSSIVEL)
```

**Teste pratico - Veja o hash:**

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    senha := "minhasenha123"
    
    // Gerar hash
    hash, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
    
    fmt.Println("Senha original:", senha)
    fmt.Println("Hash gerado:", string(hash))
    
    // Gerar outro hash da MESMA senha
    hash2, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
    
    fmt.Println("Outro hash:", string(hash2))
    fmt.Println("Sao iguais?", string(hash) == string(hash2)) // false!
}
```

**Por que hashes diferentes para mesma senha?**

O bcrypt adiciona um "SALT" (sal) aleatorio. E como temperar comida diferente cada vez.

```
senha + salt1 = hash1
senha + salt2 = hash2

hash1 != hash2, mas ambos "batem" com a senha original
```

### 3.3 Funcao CreateUser (Criar Usuario)

**Fluxograma:**

```
INICIO
   │
   ▼
┌─────────────────────┐
│ Recebe: email,      │
│ senha, nome, data   │
└──────────┬──────────┘
           │
           ▼
    ┌──────────────┐     SIM    ┌─────────────┐
    │ Email ja     │──────────>│ ERRO:       │
    │ existe?      │           │ Usuario ja  │
    └──────┬───────┘           │ existe      │
           │ NAO               └─────────────┘
           ▼
┌─────────────────────┐
│ Criptografar senha  │
│ (criar hash)        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Inserir no banco    │
│ de dados            │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Retornar usuario    │
│ criado              │
└─────────────────────┘
           │
           ▼
         FIM
```

**Codigo explicado linha por linha:**

```go
func CreateUser(email, password, name string, birthDate time.Time) (*User, error) {
    
    // PASSO 1: Verificar se email ja existe
    // -------------------------------------
    // Antes de criar, checamos se ja tem alguem com esse email
    exists, err := UserExists(email)
    
    if err != nil {
        return nil, err  // Erro ao verificar? Retorna o erro
    }
    
    if exists {
        return nil, ErrUserExists  // Ja existe? Retorna erro especifico
    }
    
    // PASSO 2: Criptografar a senha
    // -----------------------------
    // NUNCA salvamos "minhasenha123" direto!
    // Transformamos em "$2a$10$xK8f..."
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password),    // Converte string pra bytes
        bcrypt.DefaultCost,  // Nivel de seguranca (10 = padrao)
    )
    
    if err != nil {
        return nil, err
    }
    
    // PASSO 3: Montar a query SQL
    // ---------------------------
    // $1, $2, $3, $4 sao PARAMETROS (evita SQL Injection!)
    query := `
        INSERT INTO users (email, password, name, birth_date)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
    
    // PASSO 4: Criar struct para receber dados
    // ----------------------------------------
    user := &User{
        Email:     email,
        Password:  string(hashedPassword),
        Name:      name,
        BirthDate: birthDate,
    }
    
    // PASSO 5: Executar query e preencher campos retornados
    // -----------------------------------------------------
    err = database.DB.QueryRow(
        query, 
        email, 
        string(hashedPassword), 
        name, 
        birthDate,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return nil, err
    }
    
    // PASSO 6: Retornar usuario criado
    // --------------------------------
    return user, nil  // nil = sem erro
}
```

### 3.4 Funcao Authenticate (Fazer Login)

**Fluxograma:**

```
INICIO
   │
   ▼
┌─────────────────────┐
│ Recebe: email,      │
│ senha               │
└──────────┬──────────┘
           │
           ▼
    ┌──────────────┐     NAO     ┌─────────────┐
    │ Usuario      │───────────>│ ERRO:       │
    │ existe?      │            │ Nao encontr.│
    └──────┬───────┘            └─────────────┘
           │ SIM
           ▼
    ┌──────────────┐     SIM     ┌─────────────┐
    │ Esta         │───────────>│ ERRO:       │
    │ bloqueado?   │            │ Bloqueado   │
    └──────┬───────┘            └─────────────┘
           │ NAO
           ▼
    ┌──────────────┐     NAO     ┌─────────────┐
    │ Senha        │───────────>│ Incrementar │
    │ correta?     │            │ tentativas  │
    └──────┬───────┘            └──────┬──────┘
           │ SIM                       │
           │                           ▼
           │                    ┌──────────────┐
           │                    │ Tentativas   │ SIM
           │                    │ >= 3?        │────> BLOQUEAR!
           │                    └──────┬───────┘
           │                           │ NAO
           │                           ▼
           │                    ┌─────────────┐
           │                    │ ERRO: Senha │
           │                    │ incorreta   │
           │                    └─────────────┘
           ▼
┌─────────────────────┐
│ Resetar contador    │
│ de tentativas       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ SUCESSO! Retornar   │
│ dados do usuario    │
└─────────────────────┘
```

**Codigo explicado:**

```go
func Authenticate(email, password string) (*User, error) {
    
    // PASSO 1: Buscar usuario pelo email
    // ----------------------------------
    user, err := GetUserByEmail(email)
    
    if err != nil {
        // Email nao cadastrado
        return nil, ErrUserNotFound
    }
    
    // PASSO 2: Verificar bloqueio
    // ---------------------------
    if user.IsBlocked {
        // Conta bloqueada, nem olha a senha
        return nil, ErrUserBlocked
    }
    
    // PASSO 3: Comparar senha
    // -----------------------
    // NAO comparamos strings diretamente!
    // bcrypt faz a comparacao de forma segura
    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password),  // Hash salvo no banco
        []byte(password),       // Senha digitada
    )
    
    if err != nil {
        // SENHA ERRADA!
        
        // Incrementar contador de erros
        user.FailedAttempts++
        
        // Atingiu limite de 3?
        if user.FailedAttempts >= 3 {
            // BLOQUEAR CONTA!
            BlockUser(user.ID)
            return nil, ErrUserBlocked
        }
        
        // Salvar novo numero de tentativas
        UpdateFailedAttempts(user.ID, user.FailedAttempts)
        
        return nil, ErrInvalidPassword
    }
    
    // PASSO 4: Login OK!
    // ------------------
    // Zerar contador de erros (usuario acertou)
    ResetFailedAttempts(user.ID)
    
    return user, nil
}
```

### 3.5 Por que usar $1, $2, $3 nas queries?

**SQL Injection - O PERIGO:**

Imagine que o usuario digita isso no campo email:

```
'; DROP TABLE users; --
```

**SEM parametros (PERIGOSO!):**

```go
// NUNCA FACA ISSO!
query := "SELECT * FROM users WHERE email = '" + email + "'"

// A query fica:
// SELECT * FROM users WHERE email = ''; DROP TABLE users; --'
// ISSO APAGA A TABELA INTEIRA!
```

**COM parametros (SEGURO):**

```go
// SEMPRE USE ASSIM!
query := "SELECT * FROM users WHERE email = $1"
db.QueryRow(query, email)

// O banco trata o email como TEXTO, nao como codigo SQL
// Mesmo com '; DROP TABLE users; -- ele busca literalmente esse texto
```

---

## 4. Entendendo o captcha.go

### 4.1 O que e CAPTCHA?

CAPTCHA = Teste para saber se e humano ou robo

**Por que precisamos disso?**

Sem captcha, um robo (bot) poderia:
- Criar milhares de contas falsas
- Tentar senhas automaticamente
- Sobrecarregar o servidor

### 4.2 Como nosso Captcha funciona

```
FLUXO DO CAPTCHA:

1. Usuario acessa pagina de registro
         │
         ▼
┌─────────────────────────────────┐
│ Servidor GERA pergunta:         │
│ "Quanto e 7 + 4?"               │
│                                 │
│ Servidor GUARDA no banco:       │
│ Token: abc123                   │
│ Resposta: 11                    │
└─────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│ Navegador EXIBE:                │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ Quanto e 7 + 4?             │ │
│ │ Sua resposta: [____]        │ │
│ │ <input hidden="abc123">     │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
         │
         ▼ (usuario digita 11)
         │
┌─────────────────────────────────┐
│ Servidor VERIFICA:              │
│                                 │
│ Token abc123 tem resposta 11?   │
│ SIM! Usuario e humano.          │
└─────────────────────────────────┘
```

### 4.3 Funcao GenerateCaptcha

**Codigo explicado:**

```go
func GenerateCaptcha() (*Captcha, error) {
    
    // PASSO 1: Gerar dois numeros aleatorios
    // --------------------------------------
    // rand.Int gera numero entre 0 e (maximo - 1)
    // big.NewInt(20) = maximo 20, entao gera 0-19
    
    num1, _ := rand.Int(rand.Reader, big.NewInt(20))
    num2, _ := rand.Int(rand.Reader, big.NewInt(20))
    
    // Converter pra int normal e somar 1 (evitar zero)
    n1 := int(num1.Int64()) + 1  // 1 a 20
    n2 := int(num2.Int64()) + 1  // 1 a 20
    
    // PASSO 2: Escolher operacao aleatoria
    // ------------------------------------
    operations := []string{"+", "-", "*"}
    
    opIndex, _ := rand.Int(rand.Reader, big.NewInt(3))
    op := operations[opIndex.Int64()]
    
    // PASSO 3: Calcular resposta
    // --------------------------
    var answer int
    var question string
    
    switch op {
    case "+":
        answer = n1 + n2
        question = fmt.Sprintf("Quanto e %d + %d?", n1, n2)
        
    case "-":
        // Garantir resultado positivo
        if n1 < n2 {
            n1, n2 = n2, n1  // Trocar valores
        }
        answer = n1 - n2
        question = fmt.Sprintf("Quanto e %d - %d?", n1, n2)
        
    case "*":
        // Usar numeros menores (1-9) pra nao complicar
        n1 = (n1 % 10) + 1
        n2 = (n2 % 10) + 1
        answer = n1 * n2
        question = fmt.Sprintf("Quanto e %d x %d?", n1, n2)
    }
    
    // PASSO 4: Gerar token unico
    // --------------------------
    // Token e um codigo aleatorio que identifica este captcha
    tokenBytes := make([]byte, 16)
    rand.Read(tokenBytes)
    token := hex.EncodeToString(tokenBytes)
    
    // PASSO 5: Definir expiracao (5 minutos)
    // --------------------------------------
    expiresAt := time.Now().Add(5 * time.Minute)
    
    // PASSO 6: Salvar no banco
    // ------------------------
    query := `INSERT INTO captcha_tokens (token, answer, expires_at) 
              VALUES ($1, $2, $3)`
    
    _, err := database.DB.Exec(query, token, fmt.Sprintf("%d", answer), expiresAt)
    
    if err != nil {
        return nil, err
    }
    
    // PASSO 7: Retornar (sem expor a resposta!)
    // -----------------------------------------
    return &Captcha{
        Token:     token,
        Question:  question,
        Answer:    fmt.Sprintf("%d", answer),
        ExpiresAt: expiresAt,
    }, nil
}
```

### 4.4 Funcao ValidateCaptcha

**Fluxograma:**

```
INICIO
   │
   ▼
┌─────────────────────┐
│ Recebe: token,      │
│ resposta do usuario │
└──────────┬──────────┘
           │
           ▼
    ┌──────────────┐     NAO     ┌─────────────┐
    │ Token existe │───────────>│ ERRO:       │
    │ no banco?    │            │ Invalido    │
    └──────┬───────┘            └─────────────┘
           │ SIM
           ▼
    ┌──────────────┐     SIM     ┌─────────────┐
    │ Ja foi       │───────────>│ ERRO:       │
    │ usado?       │            │ Ja usado    │
    └──────┬───────┘            └─────────────┘
           │ NAO
           ▼
    ┌──────────────┐     SIM     ┌─────────────┐
    │ Expirou?     │───────────>│ ERRO:       │
    │              │            │ Expirado    │
    └──────┬───────┘            └─────────────┘
           │ NAO
           ▼
    ┌──────────────┐     NAO     ┌─────────────┐
    │ Resposta     │───────────>│ ERRO:       │
    │ correta?     │            │ Incorreta   │
    └──────┬───────┘            └─────────────┘
           │ SIM
           ▼
┌─────────────────────┐
│ Marcar como usado   │
└──────────┬──────────┘
           │
           ▼
       SUCESSO!
```

---

## 5. Entendendo o auth.go

### 5.1 O que e um Handler?

Handler e uma funcao que RECEBE uma requisicao e ENVIA uma resposta.

**Analogia do Restaurante:**

```
CLIENTE (Navegador)              GARCOM (Handler)              COZINHA (Banco)
      │                               │                              │
      │  "Quero fazer login"          │                              │
      │ ─────────────────────────────>│                              │
      │                               │  "Busca usuario X"           │
      │                               │ ────────────────────────────>│
      │                               │                              │
      │                               │        "Aqui esta"           │
      │                               │ <────────────────────────────│
      │    "Login OK! Aqui seu        │                              │
      │     cookie de sessao"         │                              │
      │ <─────────────────────────────│                              │
```

### 5.2 Estrutura de um Handler

```go
func MeuHandler(w http.ResponseWriter, r *http.Request) {
    // w = ResponseWriter = onde escrevemos a RESPOSTA
    // r = Request = contem dados da REQUISICAO
    
    // Ler dados da requisicao
    email := r.FormValue("email")      // Campo do formulario
    metodo := r.Method                  // "GET" ou "POST"
    url := r.URL.Path                   // "/login"
    
    // Escrever resposta
    w.Write([]byte("Ola!"))            // Texto simples
    // OU
    templates.ExecuteTemplate(w, "pagina.html", dados)  // HTML
    // OU
    http.Redirect(w, r, "/outra-pagina", http.StatusSeeOther)  // Redirect
}
```

### 5.3 Diferenca entre GET e POST

```
GET = BUSCAR/EXIBIR                    POST = ENVIAR/CRIAR
─────────────────                      ─────────────────
Usuario digita URL                     Usuario clica em "Enviar"
no navegador                           no formulario

/login (GET)                           /login (POST)
   │                                      │
   ▼                                      ▼
Exibe formulario                       Processa dados
de login                               do formulario
```

**No codigo:**

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    
    if r.Method == "GET" {
        // Usuario acessou a pagina
        // Mostrar o formulario vazio
        templates.ExecuteTemplate(w, "login.html", PageData{})
        return
    }
    
    // Se chegou aqui, e POST
    // Usuario enviou o formulario
    email := r.FormValue("email")
    password := r.FormValue("password")
    
    // Processar login...
}
```

### 5.4 O que sao Cookies?

Cookie e um pequeno arquivo que o SERVIDOR pede pro NAVEGADOR guardar.

```
SERVIDOR                              NAVEGADOR
    │                                     │
    │  "Guarda esse cookie:               │
    │   session_token=abc123"             │
    │ ───────────────────────────────────>│
    │                                     │  (guarda no computador)
    │                                     │
    │                                     │
    │  (proxima requisicao)               │
    │  "Aqui esta o cookie:               │
    │   session_token=abc123"             │
    │ <───────────────────────────────────│
    │                                     │
    │  "Ah, voce e o usuario abc123!"     │
```

**Criando um cookie:**

```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_token",           // Nome do cookie
    Value:    "abc123xyz...",            // Valor (nosso token)
    Path:     "/",                       // Disponivel em todas as paginas
    Expires:  time.Now().Add(24*time.Hour),  // Expira em 24h
    HttpOnly: true,                      // JavaScript NAO pode acessar
})
```

**Lendo um cookie:**

```go
cookie, err := r.Cookie("session_token")

if err != nil {
    // Cookie nao existe (usuario nao logado)
}

token := cookie.Value  // "abc123xyz..."
```

### 5.5 Fluxo Completo do Login

```
┌─────────────────────────────────────────────────────────────────────┐
│                         FLUXO DO LOGIN                              │
└─────────────────────────────────────────────────────────────────────┘

USUARIO                    SERVIDOR                      BANCO DE DADOS
   │                          │                               │
   │  GET /login              │                               │
   │ ────────────────────────>│                               │
   │                          │                               │
   │  <pagina de login>       │                               │
   │ <────────────────────────│                               │
   │                          │                               │
   │  (preenche formulario)   │                               │
   │                          │                               │
   │  POST /login             │                               │
   │  email=joao@email.com    │                               │
   │  senha=minhasenha        │                               │
   │ ────────────────────────>│                               │
   │                          │                               │
   │                          │  SELECT * FROM users          │
   │                          │  WHERE email = ?              │
   │                          │ ─────────────────────────────>│
   │                          │                               │
   │                          │  {id: 1, email: joao,         │
   │                          │   password: $2a$10$...}       │
   │                          │ <─────────────────────────────│
   │                          │                               │
   │                          │  (compara senha com hash)     │
   │                          │  (cria sessao)                │
   │                          │                               │
   │                          │  INSERT INTO sessions         │
   │                          │  (user_id, token, expires)    │
   │                          │ ─────────────────────────────>│
   │                          │                               │
   │  Set-Cookie: session=xyz │                               │
   │  Redirect: /dashboard    │                               │
   │ <────────────────────────│                               │
   │                          │                               │
   │  GET /dashboard          │                               │
   │  Cookie: session=xyz     │                               │
   │ ────────────────────────>│                               │
   │                          │                               │
   │                          │  SELECT * FROM sessions       │
   │                          │  WHERE token = xyz            │
   │                          │ ─────────────────────────────>│
   │                          │                               │
   │                          │  (sessao valida!)             │
   │                          │ <─────────────────────────────│
   │                          │                               │
   │  <dashboard do usuario>  │                               │
   │ <────────────────────────│                               │
```

---

## 6. Exercicios Praticos

### Exercicio 1: Testando Hash de Senha

Crie um arquivo `teste_hash.go`:

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    // Senha que vamos testar
    senha := "minhasenha123"
    
    fmt.Println("=== TESTE DE HASH ===")
    fmt.Println()
    
    // 1. Gerar hash
    hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    
    fmt.Println("Senha original:", senha)
    fmt.Println("Hash gerado:", string(hash))
    fmt.Println()
    
    // 2. Verificar senha CORRETA
    err = bcrypt.CompareHashAndPassword(hash, []byte("minhasenha123"))
    if err == nil {
        fmt.Println("Teste 1: Senha CORRETA validada com sucesso!")
    } else {
        fmt.Println("Teste 1: FALHOU")
    }
    
    // 3. Verificar senha ERRADA
    err = bcrypt.CompareHashAndPassword(hash, []byte("senhaerrada"))
    if err != nil {
        fmt.Println("Teste 2: Senha ERRADA rejeitada corretamente!")
    } else {
        fmt.Println("Teste 2: FALHOU - aceitou senha errada!")
    }
    
    fmt.Println()
    
    // 4. Mostrar que mesma senha gera hashes diferentes
    fmt.Println("=== DEMONSTRACAO DE SALT ===")
    for i := 1; i <= 3; i++ {
        h, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
        fmt.Printf("Hash %d: %s\n", i, string(h)[:30]+"...")
    }
    fmt.Println("(Todos diferentes, mas todos validam a mesma senha!)")
}
```

Execute:
```bash
go run teste_hash.go
```

### Exercicio 2: Simulando Bloqueio de Conta

Crie um arquivo `teste_bloqueio.go`:

```go
package main

import "fmt"

// Simulacao simples (sem banco de dados)
type Usuario struct {
    Email          string
    Senha          string
    Tentativas     int
    Bloqueado      bool
}

var usuario = Usuario{
    Email:      "joao@email.com",
    Senha:      "senha123",
    Tentativas: 0,
    Bloqueado:  false,
}

func tentarLogin(email, senha string) string {
    // Verificar bloqueio
    if usuario.Bloqueado {
        return "ERRO: Conta bloqueada! Redefina sua senha."
    }
    
    // Verificar email
    if email != usuario.Email {
        return "ERRO: Usuario nao encontrado"
    }
    
    // Verificar senha
    if senha != usuario.Senha {
        usuario.Tentativas++
        
        if usuario.Tentativas >= 3 {
            usuario.Bloqueado = true
            return "ERRO: Conta BLOQUEADA apos 3 tentativas!"
        }
        
        return fmt.Sprintf("ERRO: Senha incorreta. Tentativas: %d/3", usuario.Tentativas)
    }
    
    // Sucesso!
    usuario.Tentativas = 0
    return "SUCESSO: Login realizado!"
}

func main() {
    fmt.Println("=== TESTE DE BLOQUEIO ===")
    fmt.Println()
    
    // Tentativa 1: senha errada
    fmt.Println("Tentativa 1 (senha errada):")
    fmt.Println(tentarLogin("joao@email.com", "errada"))
    fmt.Println()
    
    // Tentativa 2: senha errada
    fmt.Println("Tentativa 2 (senha errada):")
    fmt.Println(tentarLogin("joao@email.com", "errada"))
    fmt.Println()
    
    // Tentativa 3: senha errada (vai bloquear)
    fmt.Println("Tentativa 3 (senha errada):")
    fmt.Println(tentarLogin("joao@email.com", "errada"))
    fmt.Println()
    
    // Tentativa 4: mesmo com senha certa, esta bloqueado
    fmt.Println("Tentativa 4 (senha CERTA, mas bloqueado):")
    fmt.Println(tentarLogin("joao@email.com", "senha123"))
}
```

### Exercicio 3: Criando seu proprio Captcha

Crie um arquivo `teste_captcha.go`:

```go
package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
)

func gerarCaptcha() (string, int) {
    rand.Seed(time.Now().UnixNano())
    
    n1 := rand.Intn(10) + 1
    n2 := rand.Intn(10) + 1
    
    ops := []string{"+", "-", "*"}
    op := ops[rand.Intn(3)]
    
    var resposta int
    switch op {
    case "+":
        resposta = n1 + n2
    case "-":
        if n1 < n2 {
            n1, n2 = n2, n1
        }
        resposta = n1 - n2
    case "*":
        resposta = n1 * n2
    }
    
    pergunta := fmt.Sprintf("Quanto e %d %s %d?", n1, op, n2)
    return pergunta, resposta
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Println("=== TESTE DE CAPTCHA ===")
    fmt.Println()
    
    pergunta, resposta := gerarCaptcha()
    
    fmt.Println(pergunta)
    fmt.Print("Sua resposta: ")
    
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(input)
    
    respostaUsuario, err := strconv.Atoi(input)
    
    if err != nil {
        fmt.Println("Entrada invalida!")
        return
    }
    
    if respostaUsuario == resposta {
        fmt.Println("CORRETO! Voce e humano.")
    } else {
        fmt.Printf("ERRADO! A resposta era %d\n", resposta)
    }
}
```

---

## 7. Testes no Terminal

### 7.1 Testando a API com curl

Depois de subir a aplicacao com Docker:

```bash
# Ver pagina de login (HTML)
curl http://localhost:8080/login

# Tentar login (vai redirecionar)
curl -X POST http://localhost:8080/login \
  -d "email=teste@email.com" \
  -d "password=senha123" \
  -v

# Obter novo captcha (JSON)
curl http://localhost:8080/captcha
```

### 7.2 Testando no Navegador

1. Acesse http://localhost:8080/register
2. Abra o DevTools (F12)
3. Va em Network
4. Preencha o formulario e envie
5. Observe:
   - Requisicao POST com dados do formulario
   - Resposta com redirect ou erro
   - Cookie sendo setado apos login

### 7.3 Inspecionando o Banco de Dados

```bash
# Entrar no container do PostgreSQL
docker exec -it auth-postgres psql -U postgres -d authdb

# Ver usuarios
SELECT id, email, name, failed_attempts, is_blocked FROM users;

# Ver sessoes
SELECT * FROM sessions;

# Ver captchas
SELECT * FROM captcha_tokens;

# Sair
\q
```

---

## Glossario Rapido

| Termo | Significado |
|-------|-------------|
| **Struct** | Estrutura de dados (como uma ficha) |
| **Ponteiro** | Referencia ao endereco de memoria |
| **Hash** | Funcao que transforma texto em codigo irreversivel |
| **Salt** | Valor aleatorio adicionado antes do hash |
| **Cookie** | Arquivo guardado no navegador |
| **Sessao** | Conexao ativa entre usuario e servidor |
| **Token** | Codigo unico que identifica algo |
| **Handler** | Funcao que processa requisicoes HTTP |
| **CAPTCHA** | Teste para provar que e humano |
| **SQL Injection** | Ataque que injeta codigo SQL malicioso |

---

## Proximos Passos

Apos dominar estes conceitos:

1. Adicionar envio de email para reset de senha
2. Implementar autenticacao com Google/Facebook
3. Adicionar autenticacao de dois fatores (2FA)
4. Implementar rate limiting (limitar requisicoes)
5. Adicionar logs de auditoria

Bons estudos!
