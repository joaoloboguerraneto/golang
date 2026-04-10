// =============================================================================
// MODELS/USER.GO - MODELO DE USUARIO
// =============================================================================
// Este arquivo contem toda a logica relacionada aos usuarios:
// - Criar usuario
// - Buscar usuario
// - Autenticar (login)
// - Bloquear/Desbloquear
// - Resetar senha
// - Gerenciar sessoes
//
// CONCEITOS IMPORTANTES:
// - Struct: estrutura de dados que agrupa informacoes relacionadas
// - Funcao: bloco de codigo reutilizavel que executa uma tarefa
// - Ponteiro (*): referencia ao endereco de memoria de uma variavel
// - Error: tipo especial em Go para tratar erros
// =============================================================================

package models

import (
	"crypto/rand"   // Para gerar tokens aleatorios seguros
	"encoding/hex"  // Para converter bytes em texto hexadecimal
	"errors"        // Para criar erros personalizados
	"time"          // Para trabalhar com datas e horarios

	"auth-system/database" // Nosso pacote de conexao com banco

	"golang.org/x/crypto/bcrypt" // Biblioteca para hash de senhas
)

// =============================================================================
// STRUCT USER - ESTRUTURA DO USUARIO
// =============================================================================
// Uma struct e como uma "ficha" que define quais informacoes um usuario tem.
// Cada campo tem:
// - Nome (ex: Email)
// - Tipo (ex: string, int, bool)
//
// ANALOGIA: Pense em um formulario de cadastro. Cada campo do formulario
// e um campo da struct.
// =============================================================================

type User struct {
	// ID e o identificador unico do usuario no banco de dados
	// Tipo: int (numero inteiro)
	// O banco gera automaticamente (SERIAL no PostgreSQL)
	ID int

	// Email do usuario - usado para login
	// Tipo: string (texto)
	Email string

	// Senha do usuario - NUNCA armazenamos a senha real!
	// Armazenamos um "hash" (versao criptografada)
	// Tipo: string
	Password string

	// Nome completo do usuario
	// Tipo: string
	Name string

	// Data de nascimento - usada para resetar senha
	// Tipo: time.Time (tipo especial para datas)
	BirthDate time.Time

	// Contador de tentativas de login falhas
	// Tipo: int
	// Quando chega a 3, bloqueamos a conta
	FailedAttempts int

	// Indica se a conta esta bloqueada
	// Tipo: bool (true ou false)
	IsBlocked bool

	// Data/hora em que a conta foi bloqueada
	// Tipo: *time.Time (ponteiro - pode ser nulo)
	// Usamos ponteiro porque esse campo pode estar vazio (nil)
	BlockedAt *time.Time

	// Data de criacao do registro
	CreatedAt time.Time

	// Data da ultima atualizacao
	UpdatedAt time.Time
}

// =============================================================================
// ERROS PERSONALIZADOS
// =============================================================================
// Em Go, criamos erros especificos para cada situacao.
// Isso facilita identificar o que deu errado.
//
// errors.New("mensagem") - cria um novo erro com a mensagem especificada
//
// VANTAGEM: Podemos verificar qual erro ocorreu e tomar acoes diferentes
// Ex: se for ErrUserNotFound, redirecionamos para registro
//     se for ErrInvalidPassword, mostramos "senha incorreta"
// =============================================================================

var (
	// Erro quando o usuario nao existe no banco
	ErrUserNotFound = errors.New("usuario nao encontrado")

	// Erro quando a senha esta incorreta
	ErrInvalidPassword = errors.New("senha incorreta")

	// Erro quando a conta esta bloqueada
	ErrUserBlocked = errors.New("usuario bloqueado")

	// Erro quando tentam criar usuario que ja existe
	ErrUserExists = errors.New("usuario ja existe")

	// Erro quando a data de nascimento nao confere (no reset de senha)
	ErrInvalidBirthDate = errors.New("data de nascimento incorreta")
)

// =============================================================================
// FUNCAO: CreateUser
// =============================================================================
// OBJETIVO: Criar um novo usuario no banco de dados
//
// PARAMETROS:
// - email: string com o email do usuario
// - password: string com a senha (sera criptografada)
// - name: string com o nome completo
// - birthDate: time.Time com a data de nascimento
//
// RETORNO:
// - *User: ponteiro para o usuario criado (ou nil se der erro)
// - error: erro se algo der errado (ou nil se tudo ok)
//
// FLUXO:
// 1. Verificar se usuario ja existe
// 2. Criptografar a senha
// 3. Inserir no banco de dados
// 4. Retornar o usuario criado
// =============================================================================

func CreateUser(email, password, name string, birthDate time.Time) (*User, error) {
	// -------------------------------------------------------------------------
	// PASSO 1: Verificar se usuario ja existe
	// -------------------------------------------------------------------------
	// Antes de criar, verificamos se o email ja esta cadastrado.
	// Isso evita duplicatas no banco de dados.
	exists, err := UserExists(email)

	// Se houve erro na verificacao, retornamos o erro
	if err != nil {
		return nil, err
	}

	// Se o usuario ja existe, retornamos erro especifico
	if exists {
		return nil, ErrUserExists
	}

	// -------------------------------------------------------------------------
	// PASSO 2: Criptografar a senha (HASH)
	// -------------------------------------------------------------------------
	// NUNCA armazenamos senhas em texto puro!
	// Usamos bcrypt para criar um "hash" da senha.
	//
	// O QUE E HASH?
	// E uma funcao que transforma qualquer texto em uma sequencia fixa de caracteres.
	// Exemplo: "senha123" -> "$2a$10$N9qo8uLOickgx2ZMRZoMy..."
	//
	// CARACTERISTICAS DO HASH:
	// - Unidirecional: nao da para "voltar" do hash para a senha original
	// - Deterministico: mesma entrada sempre gera mesma saida
	// - Pequena mudanca na entrada = grande mudanca na saida
	//
	// bcrypt.DefaultCost = nivel de complexidade (maior = mais seguro, mais lento)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),      // Converte string para bytes
		bcrypt.DefaultCost,    // Custo padrao (10)
	)

	// Se houve erro ao criar hash, retornamos
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------------------------
	// PASSO 3: Inserir no banco de dados
	// -------------------------------------------------------------------------
	// Query SQL com parametros ($1, $2, etc) para prevenir SQL Injection
	//
	// O QUE E SQL INJECTION?
	// E um ataque onde o usuario mal-intencionado coloca codigo SQL no input.
	// Exemplo: se o email fosse: "'; DROP TABLE users; --"
	// Sem parametros, isso poderia apagar a tabela!
	//
	// RETURNING: faz o PostgreSQL retornar os valores gerados (id, created_at, etc)
	query := `
		INSERT INTO users (email, password, name, birth_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	// Criar struct do usuario para preencher com os dados
	user := &User{
		Email:     email,
		Password:  string(hashedPassword), // Convertemos bytes para string
		Name:      name,
		BirthDate: birthDate,
	}

	// Executar a query e preencher os campos retornados
	// QueryRow: executa query que retorna uma unica linha
	// Scan: preenche as variaveis com os valores retornados
	err = database.DB.QueryRow(query, email, string(hashedPassword), name, birthDate).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Retornar o usuario criado (sucesso!)
	return user, nil
}

// =============================================================================
// FUNCAO: UserExists
// =============================================================================
// OBJETIVO: Verificar se um email ja esta cadastrado
//
// PARAMETROS:
// - email: string com o email a verificar
//
// RETORNO:
// - bool: true se existe, false se nao existe
// - error: erro se algo der errado
//
// COMO FUNCIONA:
// Contamos quantos registros existem com esse email.
// Se count > 0, o usuario existe.
// =============================================================================

func UserExists(email string) (bool, error) {
	var count int

	// COUNT(*) retorna a quantidade de registros que atendem a condicao
	query := `SELECT COUNT(*) FROM users WHERE email = $1`

	// QueryRow + Scan para obter o resultado
	err := database.DB.QueryRow(query, email).Scan(&count)

	if err != nil {
		return false, err
	}

	// Se count > 0, retorna true (usuario existe)
	return count > 0, nil
}

// =============================================================================
// FUNCAO: GetUserByEmail
// =============================================================================
// OBJETIVO: Buscar todos os dados de um usuario pelo email
//
// PARAMETROS:
// - email: string com o email a buscar
//
// RETORNO:
// - *User: ponteiro para o usuario encontrado
// - error: ErrUserNotFound se nao encontrar
//
// USO: Usado no login para verificar se o usuario existe e obter seus dados
// =============================================================================

func GetUserByEmail(email string) (*User, error) {
	// Query para buscar todos os campos do usuario
	query := `
		SELECT id, email, password, name, birth_date, 
		       failed_attempts, is_blocked, blocked_at, 
		       created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	// Criar struct vazia para preencher
	user := &User{}

	// Executar query e preencher cada campo do struct
	// A ordem dos campos no Scan DEVE corresponder a ordem no SELECT
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.BirthDate,
		&user.FailedAttempts,
		&user.IsBlocked,
		&user.BlockedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Se nao encontrou, retornamos nosso erro personalizado
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// =============================================================================
// FUNCAO: Authenticate
// =============================================================================
// OBJETIVO: Verificar se email e senha estao corretos (fazer login)
//
// PARAMETROS:
// - email: email digitado pelo usuario
// - password: senha digitada pelo usuario
//
// RETORNO:
// - *User: dados do usuario se login ok
// - error: tipo de erro (nao encontrado, senha errada, bloqueado)
//
// FLUXO:
// 1. Buscar usuario pelo email
// 2. Verificar se esta bloqueado
// 3. Comparar senha com hash
// 4. Se errou, incrementar contador
// 5. Se acertou, resetar contador
// =============================================================================

func Authenticate(email, password string) (*User, error) {
	// -------------------------------------------------------------------------
	// PASSO 1: Buscar usuario pelo email
	// -------------------------------------------------------------------------
	user, err := GetUserByEmail(email)
	if err != nil {
		// Usuario nao existe
		return nil, ErrUserNotFound
	}

	// -------------------------------------------------------------------------
	// PASSO 2: Verificar se conta esta bloqueada
	// -------------------------------------------------------------------------
	// Se estiver bloqueada, nem verificamos a senha
	if user.IsBlocked {
		return nil, ErrUserBlocked
	}

	// -------------------------------------------------------------------------
	// PASSO 3: Comparar senha digitada com hash armazenado
	// -------------------------------------------------------------------------
	// bcrypt.CompareHashAndPassword compara de forma segura
	// Ele aplica o mesmo hash na senha digitada e compara com o armazenado
	//
	// IMPORTANTE: Nao comparamos strings diretamente!
	// O bcrypt usa "salt" (valor aleatorio) que torna cada hash unico
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password), // Hash armazenado no banco
		[]byte(password),      // Senha digitada pelo usuario
	)

	// Se err != nil, a senha esta incorreta
	if err != nil {
		// ---------------------------------------------------------------------
		// PASSO 4: Senha incorreta - incrementar tentativas
		// ---------------------------------------------------------------------
		user.FailedAttempts++

		// Verificar se atingiu o limite de 3 tentativas
		if user.FailedAttempts >= 3 {
			// BLOQUEAR A CONTA!
			BlockUser(user.ID)
			return nil, ErrUserBlocked
		}

		// Salvar o novo numero de tentativas no banco
		UpdateFailedAttempts(user.ID, user.FailedAttempts)
		return nil, ErrInvalidPassword
	}

	// -------------------------------------------------------------------------
	// PASSO 5: Login bem sucedido - resetar contador de erros
	// -------------------------------------------------------------------------
	ResetFailedAttempts(user.ID)

	return user, nil
}

// =============================================================================
// FUNCAO: BlockUser
// =============================================================================
// OBJETIVO: Bloquear a conta de um usuario
//
// PARAMETROS:
// - userID: ID do usuario a bloquear
//
// RETORNO:
// - error: erro se falhar
//
// O QUE FAZ:
// - Define is_blocked = TRUE
// - Registra a data/hora do bloqueio
// =============================================================================

func BlockUser(userID int) error {
	query := `UPDATE users SET is_blocked = TRUE, blocked_at = NOW() WHERE id = $1`

	// Exec: usado para queries que nao retornam dados (UPDATE, DELETE, INSERT sem RETURNING)
	_, err := database.DB.Exec(query, userID)
	return err
}

// =============================================================================
// FUNCAO: UnblockUser
// =============================================================================
// OBJETIVO: Desbloquear a conta de um usuario
//
// PARAMETROS:
// - userID: ID do usuario a desbloquear
//
// O QUE FAZ:
// - Define is_blocked = FALSE
// - Limpa blocked_at
// - Reseta failed_attempts para 0
// =============================================================================

func UnblockUser(userID int) error {
	query := `UPDATE users SET is_blocked = FALSE, blocked_at = NULL, failed_attempts = 0 WHERE id = $1`
	_, err := database.DB.Exec(query, userID)
	return err
}

// =============================================================================
// FUNCAO: UpdateFailedAttempts
// =============================================================================
// OBJETIVO: Atualizar o contador de tentativas falhas
//
// PARAMETROS:
// - userID: ID do usuario
// - attempts: novo numero de tentativas
// =============================================================================

func UpdateFailedAttempts(userID int, attempts int) error {
	query := `UPDATE users SET failed_attempts = $1, updated_at = NOW() WHERE id = $2`
	_, err := database.DB.Exec(query, attempts, userID)
	return err
}

// =============================================================================
// FUNCAO: ResetFailedAttempts
// =============================================================================
// OBJETIVO: Zerar o contador de tentativas (apos login bem sucedido)
// =============================================================================

func ResetFailedAttempts(userID int) error {
	query := `UPDATE users SET failed_attempts = 0, updated_at = NOW() WHERE id = $1`
	_, err := database.DB.Exec(query, userID)
	return err
}

// =============================================================================
// FUNCAO: ResetPassword
// =============================================================================
// OBJETIVO: Permitir que usuario bloqueado redefina sua senha
//
// PARAMETROS:
// - email: email do usuario
// - birthDate: data de nascimento (para validacao)
// - newPassword: nova senha desejada
//
// FLUXO:
// 1. Buscar usuario
// 2. Validar data de nascimento
// 3. Criar hash da nova senha
// 4. Atualizar no banco e desbloquear
//
// POR QUE VALIDAR DATA DE NASCIMENTO?
// E uma forma simples de verificar identidade.
// Em sistemas reais, usariamos email de confirmacao ou SMS.
// =============================================================================

func ResetPassword(email string, birthDate time.Time, newPassword string) error {
	// Buscar usuario
	user, err := GetUserByEmail(email)
	if err != nil {
		return ErrUserNotFound
	}

	// -------------------------------------------------------------------------
	// Validar data de nascimento
	// -------------------------------------------------------------------------
	// Comparamos apenas ano, mes e dia (ignorando hora/minuto/segundo)
	if !sameDate(user.BirthDate, birthDate) {
		return ErrInvalidBirthDate
	}

	// -------------------------------------------------------------------------
	// Criar hash da nova senha
	// -------------------------------------------------------------------------
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// -------------------------------------------------------------------------
	// Atualizar senha e desbloquear conta
	// -------------------------------------------------------------------------
	query := `
		UPDATE users 
		SET password = $1, 
		    is_blocked = FALSE, 
		    blocked_at = NULL, 
		    failed_attempts = 0, 
		    updated_at = NOW() 
		WHERE id = $2
	`
	_, err = database.DB.Exec(query, string(hashedPassword), user.ID)
	return err
}

// =============================================================================
// FUNCAO: GetFailedAttempts
// =============================================================================
// OBJETIVO: Obter o numero de tentativas falhas de um usuario
// USO: Para mostrar "Tentativas restantes: X" na tela de login
// =============================================================================

func GetFailedAttempts(email string) (int, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return 0, err
	}
	return user.FailedAttempts, nil
}

// =============================================================================
// FUNCAO: CreateSession
// =============================================================================
// OBJETIVO: Criar uma sessao para usuario logado
//
// O QUE E SESSAO?
// Quando o usuario faz login, criamos um "token" (codigo aleatorio).
// Esse token e salvo no cookie do navegador.
// A cada requisicao, verificamos se o token e valido.
//
// PARAMETROS:
// - userID: ID do usuario logado
//
// RETORNO:
// - string: token da sessao
// - error: erro se falhar
// =============================================================================

func CreateSession(userID int) (string, error) {
	// Gerar token aleatorio de 64 caracteres hexadecimais
	token := generateToken()

	// Sessao expira em 24 horas
	expiresAt := time.Now().Add(24 * time.Hour)

	// Salvar no banco
	query := `INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, userID, token, expiresAt)

	if err != nil {
		return "", err
	}

	return token, nil
}

// =============================================================================
// FUNCAO: ValidateSession
// =============================================================================
// OBJETIVO: Verificar se um token de sessao e valido
//
// PARAMETROS:
// - token: token recebido do cookie
//
// RETORNO:
// - *User: dados do usuario se sessao valida
// - error: erro se sessao invalida ou expirada
//
// COMO FUNCIONA:
// 1. Busca a sessao pelo token
// 2. Verifica se nao expirou (expires_at > NOW())
// 3. Retorna os dados do usuario associado
// =============================================================================

func ValidateSession(token string) (*User, error) {
	// JOIN: combina dados de duas tabelas (sessions e users)
	query := `
		SELECT u.id, u.email, u.name, u.birth_date
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.token = $1 AND s.expires_at > NOW()
	`

	user := &User{}
	err := database.DB.QueryRow(query, token).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.BirthDate,
	)

	if err != nil {
		return nil, errors.New("sessao invalida")
	}

	return user, nil
}

// =============================================================================
// FUNCAO: DeleteSession
// =============================================================================
// OBJETIVO: Remover uma sessao (logout)
// =============================================================================

func DeleteSession(token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	_, err := database.DB.Exec(query, token)
	return err
}

// =============================================================================
// FUNCOES AUXILIARES (HELPERS)
// =============================================================================

// generateToken cria um token aleatorio seguro
// Usa crypto/rand que e criptograficamente seguro
// (diferente de math/rand que e previsivel)
func generateToken() string {
	// Criar slice de 32 bytes
	bytes := make([]byte, 32)

	// Preencher com valores aleatorios
	rand.Read(bytes)

	// Converter para hexadecimal (64 caracteres)
	return hex.EncodeToString(bytes)
}

// sameDate compara duas datas ignorando hora/minuto/segundo
func sameDate(d1, d2 time.Time) bool {
	// Date() retorna ano, mes e dia separados
	y1, m1, day1 := d1.Date()
	y2, m2, day2 := d2.Date()

	// Compara cada parte
	return y1 == y2 && m1 == m2 && day1 == day2
}

// =============================================================================
// EXERCICIOS PRATICOS
// =============================================================================
//
// 1. ENTENDENDO HASH:
//    Execute este codigo para ver como bcrypt funciona:
//
//    senha := "minhasenha123"
//    hash1, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
//    hash2, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
//    fmt.Println(string(hash1))
//    fmt.Println(string(hash2))
//    // Observe: mesmo senha gera hashes DIFERENTES (por causa do salt)
//
// 2. TESTANDO COMPARACAO:
//    hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
//    err1 := bcrypt.CompareHashAndPassword(hash, []byte("senha123"))
//    err2 := bcrypt.CompareHashAndPassword(hash, []byte("senhaerrada"))
//    fmt.Println("Senha correta:", err1 == nil)  // true
//    fmt.Println("Senha errada:", err2 == nil)   // false
//
// 3. TESTANDO TOKENS:
//    token1 := generateToken()
//    token2 := generateToken()
//    fmt.Println(token1)
//    fmt.Println(token2)
//    // Observe: cada token e completamente diferente
//
// 4. DESAFIO:
//    Adicione uma funcao que liste todos os usuarios bloqueados.
//    Dica: SELECT * FROM users WHERE is_blocked = TRUE
//
// =============================================================================
