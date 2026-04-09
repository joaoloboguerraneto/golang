# Curso de Programação em Go (Golang)

> De zero a criar aplicações web — Um guia completo para iniciantes

---

## Índice

1. [Sobre o Curso](#sobre-o-curso)
2. [Mapa da Jornada](#mapa-da-jornada)
3. [Aula 1: Fundamentos da Programação](#aula-1-fundamentos-da-programação)
4. [Aula 2: Lógica de Programação](#aula-2-lógica-de-programação)
5. [Aula 3: Go na Prática](#aula-3-go-na-prática)
6. [Aula 4: Criando APIs Web](#aula-4-criando-apis-web)
7. [Recursos e Próximos Passos](#recursos-e-próximos-passos)

---

## Sobre o Curso

Este curso foi criado para quem **nunca programou** e quer aprender Go (Golang) para criar sites e aplicações web.

**Go** é uma linguagem moderna criada pelo Google, conhecida por ser:
- Simples de aprender
- Rápida e eficiente
- Perfeita para aplicações web
- Usada por empresas como Google, Uber, Docker e Netflix

---

## Mapa da Jornada

```
+------------------+
| 1. Fundamentos   |  <- O que é programação
+--------+---------+
         |
         v
+------------------+
| 2. Lógica        |  <- Pensar como programador
+--------+---------+
         |
         v
+------------------+
| 3. Go Básico     |  <- Escrever código real
+--------+---------+
         |
         v
+------------------+
| 4. APIs Web      |  <- Criar aplicações
+--------+---------+
         |
         v
   Desenvolvedor Go!
```

---

## Aula 1: Fundamentos da Programação

### O que é Programação?

Programação é simplesmente **dar instruções para o computador**. O computador é muito rápido, mas muito "burro" — ele faz exatamente o que você mandar, nem mais, nem menos.

Quando você programa, você escreve um "roteiro" que o computador segue passo a passo.

### O que é uma Linguagem de Programação?

Assim como usamos Português para conversar com pessoas, usamos **linguagens de programação** para "conversar" com computadores. Go (ou Golang) é uma dessas linguagens.

Cada linguagem tem suas regras (sintaxe), assim como o Português tem gramática.

### O que é Lógica de Programação?

É a **forma de pensar** para resolver problemas. Antes de escrever código, você precisa saber *o quê* quer fazer. A lógica é universal — se você aprende a pensar logicamente, pode programar em qualquer linguagem.

**Os 3 pilares da lógica:**

| Pilar | Descrição | Exemplo |
|-------|-----------|---------|
| **Sequência** | Fazer uma coisa de cada vez, em ordem | Passo 1, Passo 2, Passo 3 |
| **Decisão** | Se acontecer X, faça Y. Senão, faça Z | Se chover, leve guarda-chuva |
| **Repetição** | Repita isso até conseguir o resultado | Enquanto não terminar, continue |

---

## Aula 2: Lógica de Programação

### 1. Variáveis: Caixas para guardar informações

Uma **variável** é como uma caixa com um nome, onde você guarda um valor.

```
+-------------+    +-------------+    +-------------+
|    idade    |    |    nome     |    |    ativo    |
|      25     |    |   "João"    |    |    true     |
|   (número)  |    |   (texto)   |    |  (booleano) |
+-------------+    +-------------+    +-------------+
```

**Os 4 tipos básicos de dados:**

| Tipo | Em Go | Exemplo | Para que serve |
|------|-------|---------|----------------|
| Número inteiro | `int` | `25`, `100`, `-5` | Idades, quantidades, contagens |
| Número decimal | `float64` | `3.14`, `99.90` | Preços, medidas, porcentagens |
| Texto | `string` | `"Olá"`, `"João"` | Nomes, mensagens, endereços |
| Verdadeiro/Falso | `bool` | `true`, `false` | Ligado/desligado, sim/não |

### 2. Condicionais: Tomando decisões

A vida é cheia de decisões: **SE** está chovendo, **ENTÃO** leve guarda-chuva, **SENÃO** vá sem.

```
           +------------------+
           |   idade >= 18?   |
           +--------+---------+
                    |
         +----------+----------+
         |                     |
    SIM  v                     v  NÃO
+-------------+         +-------------+
| Pode votar! |         | Não pode    |
|             |         | votar       |
+-------------+         +-------------+
```

**Em pseudocódigo:**

```
SE idade >= 18 ENTÃO
    mostrar "Pode votar!"
SENÃO
    mostrar "Não pode votar"
FIM
```

### 3. Loops: Repetindo ações

Em vez de escrever o mesmo código várias vezes, usamos um **loop** (laço de repetição).

```
PARA i = 1 ATÉ 5 FAÇA
    mostrar(i)
FIM

Saída: 1, 2, 3, 4, 5
```

### 4. Funções: Blocos reutilizáveis

Uma **função** é um bloco de código com um nome que você pode usar várias vezes.

```
+---------+      +------------------+      +---------+
| ENTRADA |  ->  |    dobrar()      |  ->  |  SAÍDA  |
|    5    |      | retorna n * 2    |      |   10    |
+---------+      +------------------+      +---------+
```

### Exemplo completo em pseudocódigo

```
// Função que verifica se pode dirigir
FUNÇÃO podeDirigir(idade)
    SE idade >= 18 ENTÃO
        retornar VERDADEIRO
    SENÃO
        retornar FALSO
    FIM
FIM

// Programa principal
nome = "Maria"
idade = 17

SE podeDirigir(idade) ENTÃO
    mostrar nome + " pode dirigir!"
SENÃO
    mostrar nome + " ainda não pode dirigir."
    faltam = 18 - idade
    mostrar "Faltam " + faltam + " ano(s)."
FIM

// Saída: 
// "Maria ainda não pode dirigir."
// "Faltam 1 ano(s)."
```

---

## Aula 3: Go na Prática

### Preparando o ambiente

| Opção | Quando usar |
|-------|-------------|
| **Go Playground** (online) | Para aprender e testar rápido — sem instalar nada |
| **Instalar Go** no computador | Para projetos reais — recomendado depois |

Link: https://go.dev/play/

### 1. Seu primeiro programa: Hello World

```go
package main

import "fmt"

func main() {
    fmt.Println("Olá, mundo!")
}
```

**Explicação:**

| Parte | O que faz |
|-------|-----------|
| `package main` | Todo programa Go começa com um pacote. O `main` é o principal. |
| `import "fmt"` | Importa o pacote `fmt` (format) para escrever na tela. |
| `func main()` | Função principal — onde o programa começa a executar. |
| `fmt.Println(...)` | Imprime texto na tela e pula uma linha. |

### 2. Variáveis em Go

**Forma 1: Explícita**
```go
var nome string = "João"
var idade int = 25
var ativo bool = true
```

**Forma 2: Curta (recomendada)**
```go
nome := "João"
idade := 25
ativo := true
```

O operador `:=` cria a variável e atribui valor ao mesmo tempo. Go descobre o tipo automaticamente.

**Exemplo completo:**

```go
package main

import "fmt"

func main() {
    nome := "Maria"
    idade := 22
    altura := 1.65

    fmt.Println("Nome:", nome)
    fmt.Println("Idade:", idade)
    fmt.Println("Altura:", altura)
}
```

### 3. Condicionais em Go (if/else)

| Pseudocódigo | Go |
|--------------|-----|
| `SE` | `if` |
| `ENTÃO` | `{` |
| `SENÃO` | `else` |
| `FIM` | `}` |

```go
package main

import "fmt"

func main() {
    idade := 17

    if idade >= 18 {
        fmt.Println("Maior de idade")
    } else {
        fmt.Println("Menor de idade")
    }
}
```

### 4. Loops em Go (for)

Go usa **apenas `for`** para todos os tipos de repetição:

```go
for i := 1; i <= 5; i++ {
    fmt.Println(i)
}
```

| Parte | Significado |
|-------|-------------|
| `i := 1` | Início: começa em 1 |
| `i <= 5` | Condição: enquanto <= 5 |
| `i++` | Incremento: soma 1 |

### 5. Funções em Go

```go
// Função que recebe um número e retorna o dobro
func dobrar(numero int) int {
    return numero * 2
}

func main() {
    resultado := dobrar(5)
    fmt.Println("O dobro de 5 é:", resultado)
}
```

**Anatomia de uma função:**

```
func   dobrar   (numero int)   int        { return numero * 2 }
 |       |           |          |                   |
 |       |           |          |                   +-- corpo
 |       |           |          +-- tipo de retorno
 |       |           +-- parâmetro (tipo)
 |       +-- nome
 +-- palavra-chave
```

### Projeto: Calculadora de idade

```go
package main

import "fmt"

// Função que calcula a idade baseado no ano de nascimento
func calcularIdade(anoNascimento int) int {
    anoAtual := 2024
    return anoAtual - anoNascimento
}

// Função que verifica se pode votar
func podeVotar(idade int) bool {
    return idade >= 16
}

// Função que verifica se pode dirigir
func podeDirigir(idade int) bool {
    return idade >= 18
}

func main() {
    nome := "Carlos"
    anoNascimento := 2006

    idade := calcularIdade(anoNascimento)

    fmt.Println("=== Verificador de Idade ===")
    fmt.Println("Nome:", nome)
    fmt.Println("Idade:", idade, "anos")
    fmt.Println("")

    // Verificar se pode votar
    if podeVotar(idade) {
        fmt.Println("[OK] Pode votar")
    } else {
        fmt.Println("[X] Não pode votar ainda")
    }

    // Verificar se pode dirigir
    if podeDirigir(idade) {
        fmt.Println("[OK] Pode dirigir")
    } else {
        fmt.Println("[X] Não pode dirigir ainda")
        faltam := 18 - idade
        fmt.Println("    Faltam", faltam, "ano(s)")
    }
}
```

---

## Aula 4: Criando APIs Web

### O que é uma API?

API (Interface de Programação de Aplicações) é como um **garçom** entre o cliente (navegador/app) e o servidor (seu programa Go).

```
+----------+         +--------------+         +--------------+
| Cliente  | ------> |   Sua API    | <-----> |  Banco de    |
| (Browser)| <------ | (servidor Go)|         |    Dados     |
+----------+         +--------------+         +--------------+
  Requisição           Processa               Armazena
  Resposta             Retorna
```

### 1. Seu primeiro servidor web

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Define o que acontece quando acessam "/"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Olá! Bem-vindo à minha API!")
    })

    // Inicia o servidor na porta 8080
    fmt.Println("Servidor rodando em http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

**Explicação:**

| Parte | O que faz |
|-------|-----------|
| `import "net/http"` | Pacote nativo do Go para criar servidores web |
| `http.HandleFunc("/", ...)` | Registra uma rota. Quando acessam "/", executa a função |
| `w http.ResponseWriter` | Onde você escreve a resposta |
| `r *http.Request` | Dados da requisição |
| `http.ListenAndServe(":8080", nil)` | Inicia o servidor na porta 8080 |

### 2. Criando múltiplas rotas

```go
package main

import (
    "fmt"
    "net/http"
)

// Handler para a página inicial
func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Página inicial")
}

// Handler para sobre
func sobreHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Sobre nós: Esta é uma API feita em Go!")
}

// Handler para contato
func contatoHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Email: contato@exemplo.com")
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/sobre", sobreHandler)
    http.HandleFunc("/contato", contatoHandler)

    fmt.Println("Servidor rodando em http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

### 3. Retornando JSON

APIs modernas retornam dados em formato **JSON**:

**Struct em Go -> JSON resultante:**

```go
type Usuario struct {           {
    Nome  string                    "nome": "Maria",
    Idade int         ->            "idade": 25,
    Email string                    "email": "maria@email.com"
}                               }
```

**Código completo:**

```go
package main

import (
    "encoding/json"
    "net/http"
)

// Struct que define a estrutura do usuário
type Usuario struct {
    Nome  string `json:"nome"`
    Idade int    `json:"idade"`
    Email string `json:"email"`
}

func usuariosHandler(w http.ResponseWriter, r *http.Request) {
    // Criar um usuário
    usuario := Usuario{
        Nome:  "Maria",
        Idade: 25,
        Email: "maria@email.com",
    }

    // Definir o tipo de resposta como JSON
    w.Header().Set("Content-Type", "application/json")

    // Converter struct para JSON e enviar
    json.NewEncoder(w).Encode(usuario)
}

func main() {
    http.HandleFunc("/usuario", usuariosHandler)

    println("Servidor rodando em http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

### Projeto Final: API de Lista de Tarefas

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// Struct da Tarefa
type Tarefa struct {
    ID        int    `json:"id"`
    Titulo    string `json:"titulo"`
    Concluida bool   `json:"concluida"`
}

// "Banco de dados" em memória
var tarefas = []Tarefa{
    {ID: 1, Titulo: "Aprender Go", Concluida: true},
    {ID: 2, Titulo: "Criar uma API", Concluida: true},
    {ID: 3, Titulo: "Fazer deploy", Concluida: false},
}

// GET /tarefas - Lista todas as tarefas
func listarTarefas(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tarefas)
}

// GET / - Página inicial
func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "API de Tarefas v1.0")
    fmt.Fprintf(w, "\n\nRotas disponíveis:")
    fmt.Fprintf(w, "\n  GET /tarefas - Lista todas as tarefas")
}

func main() {
    http.HandleFunc("/", home)
    http.HandleFunc("/tarefas", listarTarefas)

    fmt.Println("=================================")
    fmt.Println("  API de Tarefas rodando!")
    fmt.Println("  http://localhost:8080")
    fmt.Println("=================================")

    http.ListenAndServe(":8080", nil)
}
```

**Resposta da API (GET /tarefas):**

```json
[
    {
        "id": 1,
        "titulo": "Aprender Go",
        "concluida": true
    },
    {
        "id": 2,
        "titulo": "Criar uma API",
        "concluida": true
    },
    {
        "id": 3,
        "titulo": "Fazer deploy",
        "concluida": false
    }
]
```

---

## Bônus: Strings em Go

### O que é uma String?

Uma **string** é uma sequência de caracteres — letras, números, símbolos, espaços.

```
nome := "Olá Mundo"

+---+---+---+---+---+---+---+---+---+
| O | l | á |   | M | u | n | d | o |
+---+---+---+---+---+---+---+---+---+
  0   1   2   3   4   5   6   7   8    <- índices
```

### Operações com Strings

| Operação | Código | Resultado |
|----------|--------|-----------|
| Concatenar | `"Olá" + " " + "Mundo"` | `"Olá Mundo"` |
| Tamanho | `len("Golang")` | `6` |
| Acessar caractere | `string("Golang"[0])` | `"G"` |
| Fatiar | `"Golang"[0:2]` | `"Go"` |

### Funções do pacote strings

```go
import "strings"

strings.ToUpper("golang")           // "GOLANG"
strings.ToLower("GOLANG")           // "golang"
strings.TrimSpace("  oi  ")         // "oi"
strings.Contains("Golang", "Go")    // true
strings.Replace("Olá", "á", "a", 1) // "Ola"
strings.Split("a,b,c", ",")         // ["a", "b", "c"]
strings.Count("banana", "a")        // 3
```

### Caracteres especiais

| Código | Significado | Exemplo |
|--------|-------------|---------|
| `\n` | Nova linha | `"Linha1\nLinha2"` |
| `\t` | Tab | `"Nome:\tMaria"` |
| `\"` | Aspas | `"Ele disse \"Oi\""` |
| `\\` | Barra invertida | `"C:\\pasta"` |

---

## Resumo Geral

### Tabela de Conversão: Pseudocódigo para Go

| Conceito | Pseudocódigo | Go |
|----------|--------------|-----|
| Variável | `idade = 25` | `idade := 25` |
| Condicional | `SE ... ENTÃO ... SENÃO` | `if ... { } else { }` |
| Loop | `PARA i = 1 ATÉ 5` | `for i := 1; i <= 5; i++` |
| Função | `FUNÇÃO nome()` | `func nome() { }` |
| Imprimir | `mostrar(...)` | `fmt.Println(...)` |

### Conceitos de API Web

| Conceito | O que faz | Código |
|----------|-----------|--------|
| Servidor | Escuta requisições | `http.ListenAndServe(":8080", nil)` |
| Rota | Define um endpoint | `http.HandleFunc("/rota", handler)` |
| Handler | Processa a requisição | `func(w, r) { ... }` |
| JSON | Formato de resposta | `json.NewEncoder(w).Encode(dados)` |
| Struct | Modelo de dados | `type Usuario struct { ... }` |

---

## Recursos e Próximos Passos

### Próximos passos

| Nível | O que aprender | Por quê |
|-------|----------------|---------|
| Agora | Instalar Go no computador | Desenvolver localmente |
| Próximo | Banco de dados (PostgreSQL/MySQL) | Persistir dados de verdade |
| Depois | Framework Gin ou Echo | APIs mais robustas |
| Avançado | Docker + Deploy | Colocar no ar para o mundo |

### Links úteis

- Documentação oficial: https://go.dev/doc/
- Tour interativo: https://go.dev/tour/
- Playground online: https://go.dev/play/

---
## Conclusão