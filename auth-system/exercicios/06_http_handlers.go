// =============================================================================
// EXERCICIO 6: HTTP HANDLERS BASICOS
// =============================================================================
// Execute: go run 06_http_handlers.go
// Acesse: http://localhost:8081
//
// OBJETIVO: Entender como funcionam os handlers HTTP
//
// CONCEITO:
// Handler e uma funcao que recebe uma requisicao e envia uma resposta.
// E como um "garcom" que recebe pedidos e entrega respostas.
// =============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// -----------------------------------------------------------------------------
// TEMPLATES HTML (inline para simplificar)
// -----------------------------------------------------------------------------

const paginaInicial = `
<!DOCTYPE html>
<html>
<head>
    <title>Exercicio 6 - HTTP Handlers</title>
    <style>
        body { font-family: Arial; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .card { background: #f5f5f5; padding: 20px; margin: 20px 0; border-radius: 8px; }
        a { color: #007bff; }
        form { margin: 10px 0; }
        input, button { padding: 10px; margin: 5px; }
        button { background: #007bff; color: white; border: none; cursor: pointer; }
    </style>
</head>
<body>
    <h1>Exercicio 6: HTTP Handlers</h1>
    
    <div class="card">
        <h2>1. Handler Simples (GET)</h2>
        <p>Clique para ver um texto simples:</p>
        <a href="/ola">Ver /ola</a>
    </div>
    
    <div class="card">
        <h2>2. Handler com JSON (GET)</h2>
        <p>Retorna dados em formato JSON:</p>
        <a href="/api/hora">Ver /api/hora</a>
    </div>
    
    <div class="card">
        <h2>3. Handler com Template (GET)</h2>
        <p>Pagina HTML com dados dinamicos:</p>
        <a href="/saudacao?nome=Aluno">Ver /saudacao?nome=Aluno</a>
    </div>
    
    <div class="card">
        <h2>4. Handler com Formulario (GET + POST)</h2>
        <p>Formulario que envia dados:</p>
        <a href="/formulario">Ver /formulario</a>
    </div>
    
    <div class="card">
        <h2>5. Handler com Cookie</h2>
        <p>Demonstra como criar e ler cookies:</p>
        <a href="/cookie/criar">Criar cookie</a> |
        <a href="/cookie/ler">Ler cookie</a> |
        <a href="/cookie/apagar">Apagar cookie</a>
    </div>
    
    <div class="card">
        <h2>6. Handler com Redirect</h2>
        <p>Redireciona para outra pagina:</p>
        <a href="/redirecionar">Testar redirect</a>
    </div>
    
    <hr>
    <p><strong>Dica:</strong> Abra o DevTools (F12) -> Network para ver as requisicoes!</p>
</body>
</html>
`

const paginaSaudacao = `
<!DOCTYPE html>
<html>
<head><title>Saudacao</title></head>
<body>
    <h1>Ola, {{.Nome}}!</h1>
    <p>Bem-vindo ao nosso sistema.</p>
    <p>Hora atual: {{.Hora}}</p>
    <a href="/">Voltar</a>
</body>
</html>
`

const paginaFormulario = `
<!DOCTYPE html>
<html>
<head>
    <title>Formulario</title>
    <style>
        body { font-family: Arial; max-width: 500px; margin: 50px auto; }
        input, button { display: block; width: 100%; padding: 10px; margin: 10px 0; box-sizing: border-box; }
        button { background: #28a745; color: white; border: none; cursor: pointer; }
        .resultado { background: #d4edda; padding: 20px; border-radius: 8px; }
    </style>
</head>
<body>
    <h1>Formulario de Contato</h1>
    
    {{if .Enviado}}
    <div class="resultado">
        <h2>Dados Recebidos:</h2>
        <p><strong>Nome:</strong> {{.Nome}}</p>
        <p><strong>Email:</strong> {{.Email}}</p>
        <p><strong>Mensagem:</strong> {{.Mensagem}}</p>
    </div>
    {{else}}
    <form method="POST" action="/formulario">
        <input type="text" name="nome" placeholder="Seu nome" required>
        <input type="email" name="email" placeholder="Seu email" required>
        <textarea name="mensagem" placeholder="Sua mensagem" style="height:100px;width:100%;padding:10px;box-sizing:border-box;"></textarea>
        <button type="submit">Enviar</button>
    </form>
    {{end}}
    
    <p><a href="/">Voltar</a></p>
</body>
</html>
`

func main() {
	fmt.Println("========================================")
	fmt.Println("   EXERCICIO 6: HTTP HANDLERS")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Servidor iniciando...")
	fmt.Println("Acesse: http://localhost:8081")
	fmt.Println()
	fmt.Println("Pressione Ctrl+C para encerrar")
	fmt.Println()

	// -------------------------------------------------------------------------
	// REGISTRANDO OS HANDLERS
	// -------------------------------------------------------------------------
	// http.HandleFunc associa uma URL a uma funcao

	http.HandleFunc("/", handlerInicial)
	http.HandleFunc("/ola", handlerOla)
	http.HandleFunc("/api/hora", handlerHora)
	http.HandleFunc("/saudacao", handlerSaudacao)
	http.HandleFunc("/formulario", handlerFormulario)
	http.HandleFunc("/cookie/criar", handlerCriarCookie)
	http.HandleFunc("/cookie/ler", handlerLerCookie)
	http.HandleFunc("/cookie/apagar", handlerApagarCookie)
	http.HandleFunc("/redirecionar", handlerRedirecionar)
	http.HandleFunc("/destino", handlerDestino)

	// -------------------------------------------------------------------------
	// INICIANDO O SERVIDOR
	// -------------------------------------------------------------------------
	http.ListenAndServe(":8081", nil)
}

// =============================================================================
// HANDLER 1: Pagina Inicial
// =============================================================================
// Tipo: GET
// Retorna: HTML

func handlerInicial(w http.ResponseWriter, r *http.Request) {
	// Logar a requisicao (para aprendizado)
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Escrever HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(paginaInicial))
}

// =============================================================================
// HANDLER 2: Texto Simples
// =============================================================================
// Tipo: GET
// Retorna: Texto puro

func handlerOla(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Resposta simples em texto
	w.Write([]byte("Ola, mundo! Este e um handler simples."))
}

// =============================================================================
// HANDLER 3: JSON
// =============================================================================
// Tipo: GET
// Retorna: JSON

func handlerHora(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Criar dados para retornar
	dados := map[string]interface{}{
		"hora":      time.Now().Format("15:04:05"),
		"data":      time.Now().Format("02/01/2006"),
		"timestamp": time.Now().Unix(),
		"mensagem":  "Este e um exemplo de resposta JSON",
	}

	// Definir Content-Type como JSON
	w.Header().Set("Content-Type", "application/json")

	// Converter para JSON e enviar
	json.NewEncoder(w).Encode(dados)
}

// =============================================================================
// HANDLER 4: Template com Dados
// =============================================================================
// Tipo: GET
// Retorna: HTML com dados dinamicos

func handlerSaudacao(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Obter parametro da URL (?nome=Valor)
	nome := r.URL.Query().Get("nome")
	if nome == "" {
		nome = "Visitante"
	}

	// Dados para o template
	dados := struct {
		Nome string
		Hora string
	}{
		Nome: nome,
		Hora: time.Now().Format("15:04:05"),
	}

	// Parsear e executar template
	tmpl := template.Must(template.New("saudacao").Parse(paginaSaudacao))
	tmpl.Execute(w, dados)
}

// =============================================================================
// HANDLER 5: Formulario (GET e POST)
// =============================================================================
// Tipo: GET (mostrar formulario) + POST (processar dados)

func handlerFormulario(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Dados para o template
	dados := struct {
		Enviado  bool
		Nome     string
		Email    string
		Mensagem string
	}{}

	// Se for POST, processar os dados
	if r.Method == "POST" {
		dados.Enviado = true
		dados.Nome = r.FormValue("nome")
		dados.Email = r.FormValue("email")
		dados.Mensagem = r.FormValue("mensagem")

		fmt.Printf("  -> Dados recebidos: nome=%s, email=%s\n", dados.Nome, dados.Email)
	}

	// Renderizar template
	tmpl := template.Must(template.New("form").Parse(paginaFormulario))
	tmpl.Execute(w, dados)
}

// =============================================================================
// HANDLER 6: Criar Cookie
// =============================================================================

func handlerCriarCookie(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Criar cookie
	cookie := &http.Cookie{
		Name:     "meu_cookie",
		Value:    "valor_do_cookie_" + time.Now().Format("150405"),
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	}

	// Enviar cookie para o navegador
	http.SetCookie(w, cookie)

	w.Write([]byte("Cookie criado! Nome: meu_cookie\nValor: " + cookie.Value + "\n\nVolte e clique em 'Ler cookie'"))
}

// =============================================================================
// HANDLER 7: Ler Cookie
// =============================================================================

func handlerLerCookie(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Tentar ler o cookie
	cookie, err := r.Cookie("meu_cookie")

	if err != nil {
		w.Write([]byte("Cookie nao encontrado!\n\nClique em 'Criar cookie' primeiro."))
		return
	}

	w.Write([]byte("Cookie encontrado!\n\nNome: " + cookie.Name + "\nValor: " + cookie.Value))
}

// =============================================================================
// HANDLER 8: Apagar Cookie
// =============================================================================

func handlerApagarCookie(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Para "apagar" um cookie, criamos outro com mesma nome mas expirado
	cookie := &http.Cookie{
		Name:    "meu_cookie",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour), // Expirado!
	}

	http.SetCookie(w, cookie)

	w.Write([]byte("Cookie apagado!"))
}

// =============================================================================
// HANDLER 9: Redirect
// =============================================================================

func handlerRedirecionar(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s -> Redirecionando para /destino\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	// Redirecionar para outra pagina
	http.Redirect(w, r, "/destino", http.StatusSeeOther) // 303
}

func handlerDestino(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

	w.Write([]byte("Voce foi redirecionado para ca!\n\nO redirect e usado apos formularios (POST -> GET)."))
}
