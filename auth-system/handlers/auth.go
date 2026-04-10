// =============================================================================
// HANDLERS/AUTH.GO - CONTROLADORES DE AUTENTICACAO
// =============================================================================
// Este arquivo contem os "handlers" - funcoes que processam requisicoes HTTP.
//
// O QUE E UM HANDLER?
// E uma funcao que recebe uma requisicao HTTP e envia uma resposta.
// Pense como um "atendente" que recebe pedidos e responde.
//
// FLUXO DE UMA REQUISICAO WEB:
// 1. Usuario acessa URL (ex: /login)
// 2. Servidor identifica qual handler processar
// 3. Handler executa logica (verifica dados, acessa banco, etc)
// 4. Handler envia resposta (HTML, JSON, redirect)
//
// METODOS HTTP:
// - GET: buscar/exibir dados (ex: mostrar formulario)
// - POST: enviar/criar dados (ex: submeter formulario)
// - PUT: atualizar dados
// - DELETE: remover dados
//
// CONCEITOS IMPORTANTES:
// - http.ResponseWriter (w): onde escrevemos a resposta
// - *http.Request (r): contem dados da requisicao (metodo, URL, formulario)
// - Template: arquivo HTML com variaveis dinamicas
// - Cookie: pequeno arquivo salvo no navegador do usuario
// - Redirect: redirecionar usuario para outra pagina
// =============================================================================

package handlers

import (
	"encoding/json" // Para converter structs em JSON
	"html/template" // Para renderizar templates HTML
	"net/http"      // Pacote HTTP do Go
	"time"          // Para trabalhar com datas

	"auth-system/models" // Nossos modelos (User, Captcha)
)

// =============================================================================
// VARIAVEL GLOBAL: templates
// =============================================================================
// Carregamos todos os templates HTML uma vez na inicializacao.
// Isso e mais eficiente do que carregar a cada requisicao.
//
// template.ParseGlob("templates/*.html"):
// - Carrega todos arquivos .html da pasta templates
// - Compila os templates para uso posterior
// =============================================================================

var templates *template.Template

// init() e uma funcao especial que executa automaticamente
// quando o pacote e importado (antes do main)
func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

// =============================================================================
// STRUCT: PageData
// =============================================================================
// Estrutura que contem dados para passar aos templates HTML.
// Os templates podem acessar esses campos usando {{.NomeCampo}}
//
// Exemplo no HTML:
// <p>Erro: {{.Error}}</p>
// <input value="{{.Email}}">
// =============================================================================

type PageData struct {
	// Mensagem de erro a exibir (ex: "Senha incorreta")
	Error string

	// Mensagem de sucesso (ex: "Conta criada!")
	Success string

	// Email do usuario (para preencher campo apos erro)
	Email string

	// Nome do usuario (para exibir no dashboard)
	Name string

	// Numero de tentativas de login falhas
	Attempts int

	// Se a conta esta bloqueada
	IsBlocked bool

	// Token do captcha atual
	CaptchaToken string

	// Pergunta do captcha (ex: "Quanto e 5 + 3?")
	CaptchaQuestion string
}

// =============================================================================
// HANDLER: HomeHandler
// =============================================================================
// ROTA: GET /
// OBJETIVO: Pagina inicial - redireciona conforme estado do usuario
//
// LOGICA:
// - Se usuario logado: vai para dashboard
// - Se nao logado: vai para login
// =============================================================================

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// -------------------------------------------------------------------------
	// Verificar se usuario ja esta logado
	// -------------------------------------------------------------------------
	// Tentamos ler o cookie "session_token"
	// Se existir e for valido, usuario esta logado

	cookie, err := r.Cookie("session_token")

	// Se conseguiu ler o cookie (err == nil)
	if err == nil {
		// Validar se a sessao ainda e valida
		user, err := models.ValidateSession(cookie.Value)

		// Se sessao valida e usuario existe
		if err == nil && user != nil {
			// Redirecionar para dashboard
			// StatusSeeOther (303) = "va para esta outra pagina"
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return // Importante: encerrar funcao apos redirect
		}
	}

	// Se chegou aqui, nao esta logado - ir para login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// =============================================================================
// HANDLER: LoginHandler
// =============================================================================
// ROTA: GET /login  - Exibe formulario de login
// ROTA: POST /login - Processa tentativa de login
//
// FLUXO DO LOGIN:
// 1. Usuario acessa /login (GET) -> mostra formulario
// 2. Usuario preenche e envia (POST)
// 3. Verificamos se usuario existe
// 4. Se nao existe -> redireciona para registro
// 5. Se existe, verificamos senha
// 6. Se senha errada -> incrementa tentativas
// 7. Se 3 erros -> bloqueia conta
// 8. Se senha certa -> cria sessao e vai para dashboard
// =============================================================================

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// -------------------------------------------------------------------------
	// Verificar se ja esta logado
	// -------------------------------------------------------------------------
	cookie, err := r.Cookie("session_token")
	if err == nil {
		user, err := models.ValidateSession(cookie.Value)
		if err == nil && user != nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	// -------------------------------------------------------------------------
	// METODO GET: Exibir formulario de login
	// -------------------------------------------------------------------------
	if r.Method == "GET" {
		// Renderizar template com dados vazios
		templates.ExecuteTemplate(w, "login.html", PageData{})
		return
	}

	// -------------------------------------------------------------------------
	// METODO POST: Processar login
	// -------------------------------------------------------------------------
	// r.FormValue() obtem valores enviados pelo formulario
	// O nome deve corresponder ao atributo "name" do input HTML
	// <input name="email"> -> r.FormValue("email")

	email := r.FormValue("email")
	password := r.FormValue("password")

	// -------------------------------------------------------------------------
	// PASSO 1: Verificar se usuario existe
	// -------------------------------------------------------------------------
	exists, _ := models.UserExists(email)

	if !exists {
		// Usuario nao existe -> redirecionar para registro
		// Passamos o email na URL para pre-preencher o campo
		// &new=1 indica que veio do login (para mostrar mensagem)
		http.Redirect(w, r, "/register?email="+email+"&new=1", http.StatusSeeOther)
		return
	}

	// -------------------------------------------------------------------------
	// PASSO 2: Tentar autenticar
	// -------------------------------------------------------------------------
	user, err := models.Authenticate(email, password)

	if err != nil {
		// Autenticacao falhou - verificar qual erro
		attempts, _ := models.GetFailedAttempts(email)

		// switch: verifica qual tipo de erro ocorreu
		switch err {

		case models.ErrUserBlocked:
			// Conta bloqueada
			templates.ExecuteTemplate(w, "login.html", PageData{
				Error:     "Conta bloqueada. Por favor, redefina sua senha.",
				Email:     email,
				IsBlocked: true, // Template pode mostrar link para reset
			})

		case models.ErrInvalidPassword:
			// Senha incorreta
			remaining := 3 - attempts // Tentativas restantes
			templates.ExecuteTemplate(w, "login.html", PageData{
				Error:    "Senha incorreta. Tentativas restantes: " + string(rune('0'+remaining)),
				Email:    email,
				Attempts: attempts,
			})

		default:
			// Outro erro qualquer
			templates.ExecuteTemplate(w, "login.html", PageData{
				Error: "Erro ao fazer login. Tente novamente.",
				Email: email,
			})
		}
		return
	}

	// -------------------------------------------------------------------------
	// PASSO 3: Login bem sucedido - criar sessao
	// -------------------------------------------------------------------------
	token, err := models.CreateSession(user.ID)

	if err != nil {
		templates.ExecuteTemplate(w, "login.html", PageData{
			Error: "Erro ao criar sessao.",
		})
		return
	}

	// -------------------------------------------------------------------------
	// PASSO 4: Definir cookie no navegador
	// -------------------------------------------------------------------------
	// http.Cookie define um cookie com varias opcoes:
	// - Name: nome do cookie
	// - Value: valor (nosso token de sessao)
	// - Path: em quais URLs o cookie e enviado ("/" = todas)
	// - Expires: quando expira
	// - HttpOnly: JavaScript NAO pode acessar (seguranca contra XSS)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // Expira em 24h
		HttpOnly: true,                            // Protege contra XSS
	})

	// -------------------------------------------------------------------------
	// PASSO 5: Redirecionar para dashboard
	// -------------------------------------------------------------------------
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// =============================================================================
// HANDLER: RegisterHandler
// =============================================================================
// ROTA: GET /register  - Exibe formulario de registro
// ROTA: POST /register - Processa criacao de conta
//
// FLUXO DO REGISTRO:
// 1. Gerar captcha
// 2. Usuario preenche formulario (nome, email, senha, captcha)
// 3. Validar captcha (anti-bot)
// 4. Validar senhas coincidem
// 5. Validar tamanho da senha
// 6. Criar usuario no banco
// 7. Redirecionar para login
// =============================================================================

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// -------------------------------------------------------------------------
	// METODO GET: Exibir formulario
	// -------------------------------------------------------------------------
	if r.Method == "GET" {
		// Gerar novo captcha para o formulario
		captcha, err := models.GenerateCaptcha()

		if err != nil {
			templates.ExecuteTemplate(w, "register.html", PageData{
				Error: "Erro ao gerar captcha.",
			})
			return
		}

		// Preparar dados para o template
		data := PageData{
			Email:           r.URL.Query().Get("email"), // Pre-preencher se veio do login
			CaptchaToken:    captcha.Token,
			CaptchaQuestion: captcha.Question,
		}

		// Se veio do login com usuario nao encontrado
		if r.URL.Query().Get("new") == "1" {
			data.Error = "Usuario nao encontrado. Crie uma conta."
		}

		templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	// -------------------------------------------------------------------------
	// METODO POST: Processar registro
	// -------------------------------------------------------------------------
	// Obter todos os campos do formulario
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	name := r.FormValue("name")
	birthDateStr := r.FormValue("birth_date")     // Vem como string "2000-01-15"
	captchaToken := r.FormValue("captcha_token")  // Token oculto
	captchaAnswer := r.FormValue("captcha_answer") // Resposta do usuario

	// -------------------------------------------------------------------------
	// VALIDACAO 1: Captcha (anti-bot)
	// -------------------------------------------------------------------------
	// Verificamos PRIMEIRO o captcha para evitar que bots
	// tentem criar contas repetidamente

	if err := models.ValidateCaptcha(captchaToken, captchaAnswer); err != nil {
		// Captcha errado - gerar novo e mostrar erro
		captcha, _ := models.GenerateCaptcha()
		templates.ExecuteTemplate(w, "register.html", PageData{
			Error:           "Resposta do captcha incorreta. Tente novamente.",
			Email:           email,
			Name:            name,
			CaptchaToken:    captcha.Token,
			CaptchaQuestion: captcha.Question,
		})
		return
	}

	// -------------------------------------------------------------------------
	// VALIDACAO 2: Senhas coincidem
	// -------------------------------------------------------------------------
	if password != confirmPassword {
		captcha, _ := models.GenerateCaptcha()
		templates.ExecuteTemplate(w, "register.html", PageData{
			Error:           "As senhas nao coincidem.",
			Email:           email,
			Name:            name,
			CaptchaToken:    captcha.Token,
			CaptchaQuestion: captcha.Question,
		})
		return
	}

	// -------------------------------------------------------------------------
	// VALIDACAO 3: Tamanho minimo da senha
	// -------------------------------------------------------------------------
	if len(password) < 6 {
		captcha, _ := models.GenerateCaptcha()
		templates.ExecuteTemplate(w, "register.html", PageData{
			Error:           "A senha deve ter pelo menos 6 caracteres.",
			Email:           email,
			Name:            name,
			CaptchaToken:    captcha.Token,
			CaptchaQuestion: captcha.Question,
		})
		return
	}

	// -------------------------------------------------------------------------
	// VALIDACAO 4: Data de nascimento
	// -------------------------------------------------------------------------
	// time.Parse converte string para time.Time
	// "2006-01-02" e o formato padrao de data em Go (sim, essa data especifica!)
	// E uma convencao do Go: usar essa data como "template" do formato

	birthDate, err := time.Parse("2006-01-02", birthDateStr)

	if err != nil {
		captcha, _ := models.GenerateCaptcha()
		templates.ExecuteTemplate(w, "register.html", PageData{
			Error:           "Data de nascimento invalida.",
			Email:           email,
			Name:            name,
			CaptchaToken:    captcha.Token,
			CaptchaQuestion: captcha.Question,
		})
		return
	}

	// -------------------------------------------------------------------------
	// CRIAR USUARIO
	// -------------------------------------------------------------------------
	_, err = models.CreateUser(email, password, name, birthDate)

	if err != nil {
		captcha, _ := models.GenerateCaptcha()
		errorMsg := "Erro ao criar conta."

		// Verificar se e erro de usuario existente
		if err == models.ErrUserExists {
			errorMsg = "Este email ja esta cadastrado."
		}

		templates.ExecuteTemplate(w, "register.html", PageData{
			Error:           errorMsg,
			Email:           email,
			Name:            name,
			CaptchaToken:    captcha.Token,
			CaptchaQuestion: captcha.Question,
		})
		return
	}

	// -------------------------------------------------------------------------
	// SUCESSO: Redirecionar para login
	// -------------------------------------------------------------------------
	// ?registered=1 e usado pelo JavaScript para mostrar mensagem de sucesso
	http.Redirect(w, r, "/login?registered=1", http.StatusSeeOther)
}

// =============================================================================
// HANDLER: ResetPasswordHandler
// =============================================================================
// ROTA: GET /reset-password  - Exibe formulario de reset
// ROTA: POST /reset-password - Processa reset de senha
//
// FLUXO DO RESET:
// 1. Usuario informa email e data de nascimento
// 2. Validamos se a data de nascimento confere
// 3. Se conferir, permitimos definir nova senha
// 4. Conta e desbloqueada automaticamente
// =============================================================================

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// -------------------------------------------------------------------------
	// METODO GET: Exibir formulario
	// -------------------------------------------------------------------------
	if r.Method == "GET" {
		// Pre-preencher email se veio da pagina de login
		email := r.URL.Query().Get("email")
		templates.ExecuteTemplate(w, "reset.html", PageData{Email: email})
		return
	}

	// -------------------------------------------------------------------------
	// METODO POST: Processar reset
	// -------------------------------------------------------------------------
	email := r.FormValue("email")
	birthDateStr := r.FormValue("birth_date")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	// Validar senhas coincidem
	if newPassword != confirmPassword {
		templates.ExecuteTemplate(w, "reset.html", PageData{
			Error: "As senhas nao coincidem.",
			Email: email,
		})
		return
	}

	// Validar tamanho
	if len(newPassword) < 6 {
		templates.ExecuteTemplate(w, "reset.html", PageData{
			Error: "A senha deve ter pelo menos 6 caracteres.",
			Email: email,
		})
		return
	}

	// Converter data
	birthDate, err := time.Parse("2006-01-02", birthDateStr)

	if err != nil {
		templates.ExecuteTemplate(w, "reset.html", PageData{
			Error: "Data de nascimento invalida.",
			Email: email,
		})
		return
	}

	// Tentar resetar senha (valida data de nascimento internamente)
	err = models.ResetPassword(email, birthDate, newPassword)

	if err != nil {
		errorMsg := "Erro ao redefinir senha."

		switch err {
		case models.ErrUserNotFound:
			errorMsg = "Usuario nao encontrado."
		case models.ErrInvalidBirthDate:
			errorMsg = "Data de nascimento incorreta."
		}

		templates.ExecuteTemplate(w, "reset.html", PageData{
			Error: errorMsg,
			Email: email,
		})
		return
	}

	// Sucesso - redirecionar para login
	http.Redirect(w, r, "/login?reset=1", http.StatusSeeOther)
}

// =============================================================================
// HANDLER: DashboardHandler
// =============================================================================
// ROTA: GET /dashboard
// OBJETIVO: Exibir painel do usuario logado
//
// SEGURANCA: Verifica se usuario esta logado
// Se nao estiver, redireciona para login
// =============================================================================

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// -------------------------------------------------------------------------
	// Verificar se esta logado (TEM cookie valido?)
	// -------------------------------------------------------------------------
	cookie, err := r.Cookie("session_token")

	if err != nil {
		// Sem cookie -> nao logado -> vai para login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validar se sessao e valida
	user, err := models.ValidateSession(cookie.Value)

	if err != nil {
		// Sessao invalida ou expirada
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// -------------------------------------------------------------------------
	// Usuario logado - mostrar dashboard
	// -------------------------------------------------------------------------
	templates.ExecuteTemplate(w, "dashboard.html", PageData{
		Name:  user.Name,
		Email: user.Email,
	})
}

// =============================================================================
// HANDLER: LogoutHandler
// =============================================================================
// ROTA: GET /logout
// OBJETIVO: Encerrar sessao do usuario
//
// FLUXO:
// 1. Obter token da sessao do cookie
// 2. Deletar sessao do banco de dados
// 3. Remover cookie do navegador
// 4. Redirecionar para login
// =============================================================================

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Tentar obter cookie
	cookie, err := r.Cookie("session_token")

	if err == nil {
		// Se existe, deletar sessao do banco
		models.DeleteSession(cookie.Value)
	}

	// -------------------------------------------------------------------------
	// Remover cookie do navegador
	// -------------------------------------------------------------------------
	// Para "remover" um cookie, definimos ele com:
	// - Valor vazio
	// - Data de expiracao no passado

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",                              // Valor vazio
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Expirado (1 hora atras)
		HttpOnly: true,
	})

	// Redirecionar para login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// =============================================================================
// HANDLER: CaptchaHandler
// =============================================================================
// ROTA: GET /captcha
// OBJETIVO: Gerar novo captcha via AJAX (requisicao JavaScript)
//
// RESPOSTA: JSON com token e pergunta
// { "token": "abc123...", "question": "Quanto e 5 + 3?" }
//
// USO: Quando usuario clica no botao "Novo captcha" no formulario
// =============================================================================

func CaptchaHandler(w http.ResponseWriter, r *http.Request) {
	// Gerar novo captcha
	captcha, err := models.GenerateCaptcha()

	if err != nil {
		// http.Error envia resposta de erro simples
		http.Error(w, "Erro ao gerar captcha", http.StatusInternalServerError)
		return
	}

	// -------------------------------------------------------------------------
	// Responder com JSON
	// -------------------------------------------------------------------------
	// w.Header().Set define cabecalhos da resposta
	// Content-Type indica o formato dos dados

	w.Header().Set("Content-Type", "application/json")

	// json.NewEncoder(w).Encode() converte struct/map para JSON
	// e escreve diretamente na resposta
	json.NewEncoder(w).Encode(map[string]string{
		"token":    captcha.Token,
		"question": captcha.Question,
	})
}

// =============================================================================
// HANDLER: VerifyCaptchaHandler
// =============================================================================
// ROTA: POST /verify-captcha
// OBJETIVO: Verificar resposta do captcha via AJAX
//
// REQUISICAO: Form data com token e answer
// RESPOSTA: JSON com resultado
// { "valid": true } ou { "valid": false, "error": "mensagem" }
// =============================================================================

func VerifyCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	answer := r.FormValue("answer")

	// Validar captcha
	err := models.ValidateCaptcha(token, answer)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		// Captcha invalido
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	// Captcha valido
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid": true,
	})
}

// =============================================================================
// RESUMO DOS STATUS HTTP USADOS
// =============================================================================
//
// 200 OK              - Requisicao bem sucedida (padrao)
// 303 See Other       - Redirecionar para outra pagina (apos POST)
// 500 Internal Error  - Erro no servidor
//
// =============================================================================

// =============================================================================
// EXERCICIOS PRATICOS
// =============================================================================
//
// 1. TESTANDO ROTAS MANUALMENTE:
//    Acesse no navegador:
//    - http://localhost:8080/login
//    - http://localhost:8080/register
//    - http://localhost:8080/dashboard (deve redirecionar se nao logado)
//
// 2. INSPECIONANDO COOKIES:
//    No navegador (Chrome/Firefox):
//    - F12 -> Application -> Cookies
//    - Observe o cookie "session_token" apos login
//
// 3. TESTANDO BLOQUEIO:
//    - Crie uma conta
//    - Faca logout
//    - Tente logar com senha errada 3 vezes
//    - Observe a mensagem de bloqueio
//    - Use a data de nascimento para resetar
//
// 4. DESAFIO - ADICIONAR RATE LIMITING:
//    Limite o numero de requisicoes por IP.
//    Dica: crie um map[string]int para contar requisicoes por IP
//
// 5. DESAFIO - ADICIONAR "LEMBRAR-ME":
//    Adicione checkbox "Lembrar-me" que aumenta duracao do cookie
//    Dica: se marcado, Expires = 30 dias; senao, Expires = 24 horas
//
// =============================================================================
