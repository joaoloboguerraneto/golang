// =============================================================================
// MODELS/CAPTCHA.GO - SISTEMA ANTI-BOT
// =============================================================================
// Este arquivo implementa um CAPTCHA matematico simples para evitar que
// robos (bots) criem contas automaticamente.
//
// O QUE E CAPTCHA?
// CAPTCHA = "Completely Automated Public Turing test to tell Computers and Humans Apart"
// Em portugues: Teste publico completamente automatizado para diferenciar
// computadores de humanos.
//
// COMO NOSSO CAPTCHA FUNCIONA:
// 1. Geramos uma pergunta matematica simples (ex: "Quanto e 5 + 3?")
// 2. Salvamos a resposta correta no banco de dados
// 3. O usuario responde no formulario
// 4. Comparamos a resposta do usuario com a resposta salva
//
// POR QUE CAPTCHA MATEMATICO?
// - Simples de implementar
// - Nao depende de servicos externos
// - Funciona sem JavaScript
// - Acessivel (funciona com leitores de tela)
//
// ALTERNATIVAS MAIS ROBUSTAS (para sistemas reais):
// - Google reCAPTCHA
// - hCaptcha
// - Cloudflare Turnstile
// =============================================================================

package models

import (
	"crypto/rand"  // Gerador de numeros aleatorios seguro
	"encoding/hex" // Converter bytes para hexadecimal
	"errors"       // Criar erros personalizados
	"fmt"          // Formatacao de strings
	"math/big"     // Numeros grandes (para rand.Int)
	"time"         // Trabalhar com tempo

	"auth-system/database" // Conexao com banco
)

// =============================================================================
// STRUCT CAPTCHA
// =============================================================================
// Representa um desafio captcha
// =============================================================================

type Captcha struct {
	// Token unico que identifica este captcha
	// E enviado como campo oculto no formulario
	Token string

	// Pergunta que o usuario vera (ex: "Quanto e 5 + 3?")
	Question string

	// Resposta correta (ex: "8")
	// NUNCA enviamos isso para o navegador!
	Answer string

	// Quando este captcha expira
	// Captchas antigos nao sao aceitos (evita reuso)
	ExpiresAt time.Time
}

// =============================================================================
// FUNCAO: GenerateCaptcha
// =============================================================================
// OBJETIVO: Criar um novo desafio captcha
//
// RETORNO:
// - *Captcha: o captcha gerado com pergunta e token
// - error: erro se falhar
//
// FLUXO:
// 1. Gerar dois numeros aleatorios
// 2. Escolher uma operacao aleatoria (+, -, *)
// 3. Calcular a resposta
// 4. Criar token unico
// 5. Salvar no banco de dados
// 6. Retornar captcha (sem a resposta exposta)
// =============================================================================

func GenerateCaptcha() (*Captcha, error) {
	// -------------------------------------------------------------------------
	// PASSO 1: Gerar dois numeros aleatorios
	// -------------------------------------------------------------------------
	// rand.Int gera um numero aleatorio entre 0 e o valor maximo especificado
	//
	// crypto/rand vs math/rand:
	// - crypto/rand: seguro, imprevisivel (para criptografia)
	// - math/rand: rapido, previsivel (para jogos, simulacoes)
	//
	// big.NewInt(20) = numero maximo 20
	// Resultado: numero entre 0 e 19

	num1, _ := rand.Int(rand.Reader, big.NewInt(20))
	num2, _ := rand.Int(rand.Reader, big.NewInt(20))

	// Converter de *big.Int para int normal
	// Somamos 1 para evitar zero (0 + 0 = 0 e muito facil)
	n1 := int(num1.Int64()) + 1 // Resultado: 1 a 20
	n2 := int(num2.Int64()) + 1 // Resultado: 1 a 20

	// -------------------------------------------------------------------------
	// PASSO 2: Escolher operacao aleatoria
	// -------------------------------------------------------------------------
	// Temos 3 operacoes: adicao, subtracao, multiplicacao
	// Nao usamos divisao para evitar resultados decimais

	operations := []string{"+", "-", "*"}

	// Escolher indice aleatorio (0, 1 ou 2)
	opIndex, _ := rand.Int(rand.Reader, big.NewInt(3))
	op := operations[opIndex.Int64()]

	// -------------------------------------------------------------------------
	// PASSO 3: Calcular resposta e montar pergunta
	// -------------------------------------------------------------------------

	var answer int
	var question string

	// switch: estrutura de decisao para multiplos casos
	// Similar a varios if/else, mas mais limpo
	switch op {
	case "+":
		// Adicao: resultado sempre positivo
		answer = n1 + n2
		question = fmt.Sprintf("Quanto e %d + %d?", n1, n2)

	case "-":
		// Subtracao: garantir resultado positivo
		// Se n1 < n2, invertemos para evitar numeros negativos
		if n1 < n2 {
			n1, n2 = n2, n1 // Troca de valores (Go permite isso!)
		}
		answer = n1 - n2
		question = fmt.Sprintf("Quanto e %d - %d?", n1, n2)

	case "*":
		// Multiplicacao: usar numeros menores para nao complicar
		// 20 * 20 = 400, pode ser dificil de cabeca
		// Limitamos a 1-9 para facilitar
		n1 = n1 % 10 // Resto da divisao por 10 (0-9)
		n2 = n2 % 10
		if n1 == 0 {
			n1 = 1 // Evita multiplicar por zero
		}
		if n2 == 0 {
			n2 = 1
		}
		answer = n1 * n2
		question = fmt.Sprintf("Quanto e %d x %d?", n1, n2)
	}

	// -------------------------------------------------------------------------
	// PASSO 4: Gerar token unico
	// -------------------------------------------------------------------------
	// O token identifica este captcha especifico
	// E enviado como campo oculto (<input type="hidden">)
	// Quando o usuario responde, enviamos o token junto
	// Assim sabemos qual resposta correta comparar

	tokenBytes := make([]byte, 16) // 16 bytes = 32 caracteres hex
	rand.Read(tokenBytes)          // Preenche com bytes aleatorios
	token := hex.EncodeToString(tokenBytes)

	// -------------------------------------------------------------------------
	// PASSO 5: Definir expiracao
	// -------------------------------------------------------------------------
	// Captcha expira em 5 minutos
	// Isso evita que alguem guarde captchas para usar depois

	expiresAt := time.Now().Add(5 * time.Minute)

	// -------------------------------------------------------------------------
	// PASSO 6: Salvar no banco de dados
	// -------------------------------------------------------------------------
	// Salvamos token, resposta correta e data de expiracao
	// A resposta NUNCA vai para o navegador!

	query := `INSERT INTO captcha_tokens (token, answer, expires_at) VALUES ($1, $2, $3)`

	_, err := database.DB.Exec(query, token, fmt.Sprintf("%d", answer), expiresAt)
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------------------------
	// PASSO 7: Retornar captcha
	// -------------------------------------------------------------------------
	// Retornamos a pergunta e o token
	// A resposta esta no struct mas NAO sera enviada ao usuario

	return &Captcha{
		Token:     token,
		Question:  question,
		Answer:    fmt.Sprintf("%d", answer), // Guardamos como string
		ExpiresAt: expiresAt,
	}, nil
}

// =============================================================================
// FUNCAO: ValidateCaptcha
// =============================================================================
// OBJETIVO: Verificar se a resposta do usuario esta correta
//
// PARAMETROS:
// - token: token do captcha (veio do campo oculto do formulario)
// - answer: resposta digitada pelo usuario
//
// RETORNO:
// - error: nil se correto, erro especifico se incorreto
//
// VERIFICACOES:
// 1. Token existe no banco?
// 2. Captcha ja foi usado?
// 3. Captcha expirou?
// 4. Resposta esta correta?
// =============================================================================

func ValidateCaptcha(token, answer string) error {
	// -------------------------------------------------------------------------
	// PASSO 1: Buscar captcha no banco
	// -------------------------------------------------------------------------
	query := `
		SELECT answer, expires_at, used 
		FROM captcha_tokens 
		WHERE token = $1
	`

	var storedAnswer string // Resposta correta salva no banco
	var expiresAt time.Time // Quando expira
	var used bool           // Se ja foi usado

	err := database.DB.QueryRow(query, token).Scan(&storedAnswer, &expiresAt, &used)

	// Se nao encontrou o token, e invalido
	if err != nil {
		return errors.New("captcha invalido")
	}

	// -------------------------------------------------------------------------
	// PASSO 2: Verificar se ja foi usado
	// -------------------------------------------------------------------------
	// Cada captcha so pode ser usado UMA VEZ
	// Isso evita que alguem descubra a resposta e reuse o captcha
	if used {
		return errors.New("captcha ja utilizado")
	}

	// -------------------------------------------------------------------------
	// PASSO 3: Verificar se expirou
	// -------------------------------------------------------------------------
	// time.Now().After(expiresAt) = "agora e depois de expiresAt?"
	if time.Now().After(expiresAt) {
		return errors.New("captcha expirado")
	}

	// -------------------------------------------------------------------------
	// PASSO 4: Verificar resposta
	// -------------------------------------------------------------------------
	// Comparacao simples de strings
	if storedAnswer != answer {
		return errors.New("resposta incorreta")
	}

	// -------------------------------------------------------------------------
	// PASSO 5: Marcar como usado
	// -------------------------------------------------------------------------
	// Mesmo que a resposta esteja correta, marcamos como usado
	// para evitar reutilizacao
	updateQuery := `UPDATE captcha_tokens SET used = TRUE WHERE token = $1`
	database.DB.Exec(updateQuery, token)

	// Tudo certo!
	return nil
}

// =============================================================================
// FUNCAO: CleanExpiredCaptchas
// =============================================================================
// OBJETIVO: Limpar captchas antigos do banco de dados
//
// POR QUE LIMPAR?
// - Captchas expirados ocupam espaco no banco
// - Captchas usados nao tem mais utilidade
// - Manter o banco limpo melhora performance
//
// QUANDO USAR:
// Pode ser chamada periodicamente (ex: a cada hora)
// Em producao, usariamos um "cron job" ou scheduler
// =============================================================================

func CleanExpiredCaptchas() error {
	// Deleta captchas que:
	// - Expiraram (expires_at < agora)
	// - OU ja foram usados
	query := `DELETE FROM captcha_tokens WHERE expires_at < NOW() OR used = TRUE`

	_, err := database.DB.Exec(query)
	return err
}

// =============================================================================
// COMO O CAPTCHA FUNCIONA NO FLUXO COMPLETO
// =============================================================================
//
// 1. USUARIO ACESSA PAGINA DE REGISTRO:
//    - Handler chama GenerateCaptcha()
//    - Captcha salvo no banco com token + resposta
//    - HTML renderizado com pergunta e token (campo oculto)
//
//    <form>
//      <p>Quanto e 5 + 3?</p>                     <!-- Pergunta visivel -->
//      <input type="hidden" name="captcha_token" value="abc123...">  <!-- Token oculto -->
//      <input type="text" name="captcha_answer">  <!-- Resposta do usuario -->
//    </form>
//
// 2. USUARIO PREENCHE E ENVIA:
//    - Formulario envia: token="abc123", answer="8"
//
// 3. SERVIDOR VALIDA:
//    - Handler chama ValidateCaptcha("abc123", "8")
//    - Busca no banco: token="abc123" tem resposta="8"? SIM!
//    - Marca captcha como usado
//    - Continua com registro do usuario
//
// 4. SE USUARIO ERRAR:
//    - ValidateCaptcha retorna erro
//    - Handler gera NOVO captcha
//    - Renderiza formulario novamente com nova pergunta
//
// =============================================================================

// =============================================================================
// EXERCICIOS PRATICOS
// =============================================================================
//
// 1. TESTANDO GERACAO:
//    Execute varias vezes para ver perguntas diferentes:
//
//    for i := 0; i < 5; i++ {
//        captcha, _ := GenerateCaptcha()
//        fmt.Printf("Pergunta: %s (Resposta: %s)\n", captcha.Question, captcha.Answer)
//    }
//
// 2. TESTANDO VALIDACAO:
//    captcha, _ := GenerateCaptcha()
//    fmt.Println("Pergunta:", captcha.Question)
//    fmt.Println("Token:", captcha.Token)
//
//    // Resposta correta
//    err := ValidateCaptcha(captcha.Token, captcha.Answer)
//    fmt.Println("Correto:", err == nil)
//
//    // Tentar usar de novo (deve falhar)
//    err = ValidateCaptcha(captcha.Token, captcha.Answer)
//    fmt.Println("Reuso:", err) // "captcha ja utilizado"
//
// 3. DESAFIO - ADICIONAR DIVISAO:
//    Modifique GenerateCaptcha para incluir divisao.
//    Dicas:
//    - Use numeros que dividem exatamente (ex: 8 / 2 = 4)
//    - Verifique se o divisor nao e zero
//    - Verifique se a divisao e exata (sem resto)
//
// 4. DESAFIO - CAPTCHA DE PALAVRAS:
//    Crie um captcha que mostre uma palavra e peca para
//    o usuario digitar ao contrario.
//    Ex: "Digite CASA ao contrario" -> "ASAC"
//
// =============================================================================

// =============================================================================
// SEGURANCA DO CAPTCHA
// =============================================================================
//
// PONTOS FORTES:
// - Resposta nunca vai para o navegador
// - Token unico para cada captcha
// - Expiracao evita ataques de replay
// - Uso unico evita reutilizacao
//
// LIMITACOES:
// - Bots avancados podem resolver matematica simples
// - OCR pode ler texto de imagens
// - Usuarios podem achar chato
//
// MELHORIAS POSSIVEIS:
// - Aumentar complexidade das operacoes
// - Adicionar imagens (mais dificil para bots)
// - Rate limiting (limitar tentativas por IP)
// - Integrar com servicos como reCAPTCHA
//
// =============================================================================
