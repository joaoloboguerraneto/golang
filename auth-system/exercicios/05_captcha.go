// =============================================================================
// EXERCICIO 5: CAPTCHA MATEMATICO
// =============================================================================
// Execute: go run 05_captcha.go
//
// OBJETIVO: Entender como funciona um sistema anti-bot
//
// CONCEITO:
// CAPTCHA e um teste para diferenciar humanos de robos.
// Nosso captcha usa perguntas matematicas simples.
// =============================================================================

package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// STRUCT CAPTCHA
// -----------------------------------------------------------------------------

type Captcha struct {
	Token     string    // Identificador unico
	Pergunta  string    // "Quanto e 5 + 3?"
	Resposta  int       // 8
	ExpiraEm  time.Time // Quando expira
	Usado     bool      // Ja foi usado?
}

// -----------------------------------------------------------------------------
// "BANCO DE DADOS" DE CAPTCHAS
// -----------------------------------------------------------------------------

var captchas = make(map[string]*Captcha)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("   EXERCICIO 5: CAPTCHA MATEMATICO")
	fmt.Println("========================================")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 1: Gerando um captcha
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 1: GERANDO CAPTCHA ---")
	fmt.Println()

	captcha := gerarCaptcha()

	fmt.Println("Captcha gerado!")
	fmt.Println("Token:", captcha.Token)
	fmt.Println("Pergunta:", captcha.Pergunta)
	fmt.Println("Resposta (secreta):", captcha.Resposta)
	fmt.Println("Expira em:", captcha.ExpiraEm.Format("15:04:05"))
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 2: Validando a resposta
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 2: VALIDANDO RESPOSTA ---")
	fmt.Println()

	fmt.Println("PERGUNTA:", captcha.Pergunta)
	fmt.Print("Sua resposta: ")

	entrada, _ := reader.ReadString('\n')
	entrada = strings.TrimSpace(entrada)

	resposta, err := strconv.Atoi(entrada)
	if err != nil {
		fmt.Println("Entrada invalida! Digite apenas numeros.")
	} else {
		resultado := validarCaptcha(captcha.Token, resposta)
		fmt.Println()
		fmt.Println("Resultado:", resultado)
	}
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 3: Tentando reusar o captcha
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 3: TENTANDO REUSAR ---")
	fmt.Println()

	fmt.Println("Tentando usar o mesmo captcha novamente...")
	resultado := validarCaptcha(captcha.Token, captcha.Resposta)
	fmt.Println("Resultado:", resultado)
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 4: Testando expiracao
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 4: TESTANDO EXPIRACAO ---")
	fmt.Println()

	// Criar captcha que expira em 2 segundos
	captchaRapido := &Captcha{
		Token:    gerarToken(),
		Pergunta: "Quanto e 1 + 1?",
		Resposta: 2,
		ExpiraEm: time.Now().Add(2 * time.Second),
		Usado:    false,
	}
	captchas[captchaRapido.Token] = captchaRapido

	fmt.Println("Captcha criado com expiracao de 2 segundos")
	fmt.Println("Esperando 3 segundos...")
	time.Sleep(3 * time.Second)

	resultado = validarCaptcha(captchaRapido.Token, 2)
	fmt.Println("Resultado apos expiracao:", resultado)
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 5: Multiplos captchas
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 5: GERANDO VARIOS CAPTCHAS ---")
	fmt.Println()

	for i := 1; i <= 5; i++ {
		c := gerarCaptcha()
		fmt.Printf("%d. %s (Resposta: %d)\n", i, c.Pergunta, c.Resposta)
	}
	fmt.Println()

	// -------------------------------------------------------------------------
	// EXPLICACAO DO FLUXO
	// -------------------------------------------------------------------------
	fmt.Println("========================================")
	fmt.Println("   COMO FUNCIONA NO SISTEMA REAL")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("1. Usuario acessa /register")
	fmt.Println("   -> Servidor gera captcha e salva no banco")
	fmt.Println("   -> HTML mostra: pergunta + campo resposta + token oculto")
	fmt.Println()
	fmt.Println("2. Usuario preenche e envia formulario")
	fmt.Println("   -> Servidor recebe: token + resposta")
	fmt.Println("   -> Busca captcha pelo token")
	fmt.Println("   -> Verifica: existe? ja usado? expirou? resposta correta?")
	fmt.Println()
	fmt.Println("3. Se correto, continua com o registro")
	fmt.Println("   Se errado, gera NOVO captcha e pede pra tentar de novo")
	fmt.Println()

	// -------------------------------------------------------------------------
	// DESAFIOS
	// -------------------------------------------------------------------------
	fmt.Println("========================================")
	fmt.Println("   DESAFIOS PARA VOCE:")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("1. Adicione operacao de DIVISAO (cuidado com resto!)")
	fmt.Println()
	fmt.Println("2. Crie um captcha de PALAVRAS:")
	fmt.Println("   'Digite GATO ao contrario' -> 'OTAG'")
	fmt.Println()
	fmt.Println("3. Crie um captcha de CORES:")
	fmt.Println("   'Qual cor: vermelho, AZUL, verde?' -> 'AZUL'")
	fmt.Println()
	fmt.Println("4. Adicione um contador de tentativas por IP")
	fmt.Println("   (para limitar ataques de forca bruta)")
}

// -----------------------------------------------------------------------------
// FUNCAO: gerarCaptcha
// -----------------------------------------------------------------------------

func gerarCaptcha() *Captcha {
	// Gerar dois numeros aleatorios (1-20)
	n1 := gerarNumeroAleatorio(1, 20)
	n2 := gerarNumeroAleatorio(1, 20)

	// Escolher operacao
	operacoes := []string{"+", "-", "*"}
	opIndex := gerarNumeroAleatorio(0, 2)
	op := operacoes[opIndex]

	var resposta int
	var pergunta string

	switch op {
	case "+":
		resposta = n1 + n2
		pergunta = fmt.Sprintf("Quanto e %d + %d?", n1, n2)

	case "-":
		// Garantir resultado positivo
		if n1 < n2 {
			n1, n2 = n2, n1
		}
		resposta = n1 - n2
		pergunta = fmt.Sprintf("Quanto e %d - %d?", n1, n2)

	case "*":
		// Usar numeros menores para multiplicacao
		n1 = n1 % 10
		n2 = n2 % 10
		if n1 == 0 {
			n1 = 1
		}
		if n2 == 0 {
			n2 = 1
		}
		resposta = n1 * n2
		pergunta = fmt.Sprintf("Quanto e %d x %d?", n1, n2)
	}

	// Criar token unico
	token := gerarToken()

	// Criar captcha
	captcha := &Captcha{
		Token:    token,
		Pergunta: pergunta,
		Resposta: resposta,
		ExpiraEm: time.Now().Add(5 * time.Minute),
		Usado:    false,
	}

	// "Salvar no banco"
	captchas[token] = captcha

	return captcha
}

// -----------------------------------------------------------------------------
// FUNCAO: validarCaptcha
// -----------------------------------------------------------------------------

func validarCaptcha(token string, resposta int) string {
	// Buscar captcha
	captcha, existe := captchas[token]

	if !existe {
		return "ERRO: Captcha invalido (token nao encontrado)"
	}

	if captcha.Usado {
		return "ERRO: Captcha ja foi utilizado"
	}

	if time.Now().After(captcha.ExpiraEm) {
		return "ERRO: Captcha expirado"
	}

	if resposta != captcha.Resposta {
		return fmt.Sprintf("ERRO: Resposta incorreta (esperado: %d)", captcha.Resposta)
	}

	// Marcar como usado
	captcha.Usado = true

	return "SUCESSO: Captcha validado corretamente!"
}

// -----------------------------------------------------------------------------
// FUNCOES AUXILIARES
// -----------------------------------------------------------------------------

func gerarNumeroAleatorio(min, max int) int {
	diff := big.NewInt(int64(max - min + 1))
	n, _ := rand.Int(rand.Reader, diff)
	return int(n.Int64()) + min
}

func gerarToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
