// =============================================================================
// EXERCICIO 4: SISTEMA DE BLOQUEIO DE CONTA
// =============================================================================
// Execute: go run 04_bloqueio.go
//
// OBJETIVO: Entender a logica de bloqueio de conta apos tentativas falhas
//
// REGRA DE NEGOCIO:
// - Usuario tem 3 tentativas para acertar a senha
// - Na 3a tentativa errada, conta e bloqueada
// - Para desbloquear, precisa informar data de nascimento
// =============================================================================

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// -----------------------------------------------------------------------------
// STRUCT USUARIO (simplificada para o exercicio)
// -----------------------------------------------------------------------------

type Usuario struct {
	Email          string
	Senha          string // Em producao seria hash!
	Nome           string
	DataNascimento string
	Tentativas     int
	Bloqueado      bool
}

// -----------------------------------------------------------------------------
// "BANCO DE DADOS" (em memoria)
// -----------------------------------------------------------------------------

var usuarios = map[string]*Usuario{
	"joao@email.com": {
		Email:          "joao@email.com",
		Senha:          "senha123",
		Nome:           "Joao Silva",
		DataNascimento: "1990-05-15",
		Tentativas:     0,
		Bloqueado:      false,
	},
	"maria@email.com": {
		Email:          "maria@email.com",
		Senha:          "abc456",
		Nome:           "Maria Santos",
		DataNascimento: "1985-03-20",
		Tentativas:     0,
		Bloqueado:      false,
	},
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("   EXERCICIO 4: SISTEMA DE BLOQUEIO")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Usuarios disponiveis para teste:")
	fmt.Println("- joao@email.com (senha: senha123)")
	fmt.Println("- maria@email.com (senha: abc456)")
	fmt.Println()
	fmt.Println("Comandos: 'sair' para encerrar, 'reset' para resetar usuarios")
	fmt.Println()

	for {
		fmt.Println("----------------------------------------")
		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		if email == "sair" {
			fmt.Println("Ate mais!")
			break
		}

		if email == "reset" {
			resetarUsuarios()
			fmt.Println("Usuarios resetados!")
			continue
		}

		// Verificar se usuario existe
		usuario, existe := usuarios[email]
		if !existe {
			fmt.Println()
			fmt.Println(">>> USUARIO NAO ENCONTRADO <<<")
			fmt.Println("    Em um sistema real, redirecionariamos para registro.")
			continue
		}

		// Verificar se esta bloqueado
		if usuario.Bloqueado {
			fmt.Println()
			fmt.Println(">>> CONTA BLOQUEADA <<<")
			fmt.Println("    Deseja resetar a senha? (s/n)")
			resposta, _ := reader.ReadString('\n')
			resposta = strings.TrimSpace(resposta)

			if resposta == "s" || resposta == "S" {
				resetarSenha(reader, usuario)
			}
			continue
		}

		// Pedir senha
		fmt.Print("Senha: ")
		senha, _ := reader.ReadString('\n')
		senha = strings.TrimSpace(senha)

		// Tentar autenticar
		resultado := autenticar(usuario, senha)
		fmt.Println()
		fmt.Println(resultado)
	}
}

// -----------------------------------------------------------------------------
// FUNCAO: autenticar
// -----------------------------------------------------------------------------
// Esta e a logica principal que voces vao encontrar no user.go
// Aqui esta simplificada para entendimento

func autenticar(usuario *Usuario, senha string) string {
	// Verificar se esta bloqueado (checagem dupla por seguranca)
	if usuario.Bloqueado {
		return ">>> ERRO: Conta bloqueada! <<<"
	}

	// Verificar senha
	if senha != usuario.Senha {
		// SENHA INCORRETA!

		// Incrementar contador de tentativas
		usuario.Tentativas++

		// Calcular tentativas restantes
		restantes := 3 - usuario.Tentativas

		// Verificar se atingiu limite
		if usuario.Tentativas >= 3 {
			// BLOQUEAR!
			usuario.Bloqueado = true
			return fmt.Sprintf(
				">>> CONTA BLOQUEADA! <<<\n"+
					"    Voce errou 3 vezes.\n"+
					"    Para desbloquear, informe sua data de nascimento.",
			)
		}

		return fmt.Sprintf(
			">>> SENHA INCORRETA <<<\n"+
				"    Tentativas restantes: %d",
			restantes,
		)
	}

	// SENHA CORRETA!

	// Resetar contador de tentativas
	usuario.Tentativas = 0

	return fmt.Sprintf(
		">>> LOGIN REALIZADO COM SUCESSO! <<<\n"+
			"    Bem-vindo(a), %s!",
		usuario.Nome,
	)
}

// -----------------------------------------------------------------------------
// FUNCAO: resetarSenha
// -----------------------------------------------------------------------------
// Simula o fluxo de reset de senha com validacao de data de nascimento

func resetarSenha(reader *bufio.Reader, usuario *Usuario) {
	fmt.Println()
	fmt.Println("=== RESET DE SENHA ===")
	fmt.Println("Para confirmar sua identidade, informe sua data de nascimento.")
	fmt.Println("(Formato: AAAA-MM-DD)")
	fmt.Println()

	fmt.Print("Data de nascimento: ")
	data, _ := reader.ReadString('\n')
	data = strings.TrimSpace(data)

	// Verificar data de nascimento
	if data != usuario.DataNascimento {
		fmt.Println()
		fmt.Println(">>> DATA INCORRETA <<<")
		fmt.Println("    Nao foi possivel verificar sua identidade.")
		return
	}

	// Data correta - pedir nova senha
	fmt.Print("Nova senha: ")
	novaSenha, _ := reader.ReadString('\n')
	novaSenha = strings.TrimSpace(novaSenha)

	if len(novaSenha) < 6 {
		fmt.Println()
		fmt.Println(">>> SENHA MUITO CURTA <<<")
		fmt.Println("    A senha deve ter pelo menos 6 caracteres.")
		return
	}

	// Atualizar senha e desbloquear
	usuario.Senha = novaSenha
	usuario.Bloqueado = false
	usuario.Tentativas = 0

	fmt.Println()
	fmt.Println(">>> SENHA ALTERADA COM SUCESSO! <<<")
	fmt.Println("    Sua conta foi desbloqueada.")
	fmt.Println("    Faca login com a nova senha.")
}

// -----------------------------------------------------------------------------
// FUNCAO: resetarUsuarios
// -----------------------------------------------------------------------------
// Para facilitar os testes, volta tudo ao estado inicial

func resetarUsuarios() {
	for _, u := range usuarios {
		u.Tentativas = 0
		u.Bloqueado = false
	}
	usuarios["joao@email.com"].Senha = "senha123"
	usuarios["maria@email.com"].Senha = "abc456"
}
