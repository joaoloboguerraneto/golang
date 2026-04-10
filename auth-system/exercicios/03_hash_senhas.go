// =============================================================================
// EXERCICIO 3: HASH DE SENHAS
// =============================================================================
// Execute: go run 03_hash_senhas.go
//
// ANTES DE EXECUTAR, instale a biblioteca:
// go get golang.org/x/crypto/bcrypt
//
// OBJETIVO: Entender por que e como criptografamos senhas
//
// CONCEITO:
// NUNCA armazenamos senhas em texto puro!
// Usamos uma funcao HASH que transforma a senha em um codigo irreversivel.
// =============================================================================

package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("   EXERCICIO 3: HASH DE SENHAS")
	fmt.Println("========================================")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 1: Por que nao salvar senha em texto puro?
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 1: O PROBLEMA ---")
	fmt.Println()
	fmt.Println("Imagine um banco de dados assim:")
	fmt.Println()
	fmt.Println("| email              | senha        |")
	fmt.Println("|--------------------|--------------|")
	fmt.Println("| joao@email.com     | minhasenha   |")
	fmt.Println("| maria@email.com    | 123456       |")
	fmt.Println()
	fmt.Println("Se um HACKER invadir o banco, ele ve TODAS as senhas!")
	fmt.Println("E muita gente usa a mesma senha em varios sites...")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 2: O que e Hash?
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 2: O QUE E HASH? ---")
	fmt.Println()
	fmt.Println("Hash e uma funcao que:")
	fmt.Println("1. Transforma QUALQUER texto em um codigo de tamanho fixo")
	fmt.Println("2. E UNIDIRECIONAL (nao da pra voltar)")
	fmt.Println("3. Pequena mudanca na entrada = grande mudanca na saida")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 3: Gerando hash com bcrypt
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 3: GERANDO HASH ---")
	fmt.Println()

	senha := "minhasenha123"

	// Gerar hash
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	fmt.Println("Senha original:", senha)
	fmt.Println("Hash gerado:", string(hash))
	fmt.Println()
	fmt.Println("Tamanho do hash:", len(hash), "caracteres")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 4: Salt - Por que hashes sao diferentes?
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 4: SALT (SAL) ---")
	fmt.Println()
	fmt.Println("Gerando 3 hashes da MESMA senha:")
	fmt.Println()

	for i := 1; i <= 3; i++ {
		h, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
		fmt.Printf("Hash %d: %s...\n", i, string(h)[:40])
	}

	fmt.Println()
	fmt.Println("Por que sao DIFERENTES?")
	fmt.Println("O bcrypt adiciona um 'salt' (valor aleatorio) antes de gerar o hash.")
	fmt.Println("Isso impede ataques de 'rainbow table' (tabelas pre-calculadas).")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 5: Verificando senha
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 5: VERIFICANDO SENHA ---")
	fmt.Println()

	// Gerar um hash para testar
	hashParaTeste, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

	fmt.Println("Hash armazenado:", string(hashParaTeste)[:40]+"...")
	fmt.Println()

	// Testar senha CORRETA
	err = bcrypt.CompareHashAndPassword(hashParaTeste, []byte("senha123"))
	if err == nil {
		fmt.Println("Teste 1 - 'senha123': CORRETA!")
	} else {
		fmt.Println("Teste 1 - 'senha123': incorreta")
	}

	// Testar senha ERRADA
	err = bcrypt.CompareHashAndPassword(hashParaTeste, []byte("senhaerrada"))
	if err == nil {
		fmt.Println("Teste 2 - 'senhaerrada': CORRETA!")
	} else {
		fmt.Println("Teste 2 - 'senhaerrada': incorreta (como esperado)")
	}

	// Testar senha PARECIDA
	err = bcrypt.CompareHashAndPassword(hashParaTeste, []byte("Senha123"))
	if err == nil {
		fmt.Println("Teste 3 - 'Senha123' (maiuscula): CORRETA!")
	} else {
		fmt.Println("Teste 3 - 'Senha123' (maiuscula): incorreta (case sensitive!)")
	}
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 6: Custo do bcrypt
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 6: CUSTO DO BCRYPT ---")
	fmt.Println()
	fmt.Println("bcrypt.DefaultCost = 10")
	fmt.Println()
	fmt.Println("O 'custo' define quantas vezes o algoritmo e repetido.")
	fmt.Println("Maior custo = mais seguro, mas mais lento.")
	fmt.Println()
	fmt.Println("| Custo | Tempo aproximado |")
	fmt.Println("|-------|------------------|")
	fmt.Println("|   10  |   ~100ms         |")
	fmt.Println("|   12  |   ~300ms         |")
	fmt.Println("|   14  |   ~1 segundo     |")
	fmt.Println()
	fmt.Println("Por que isso e BOM?")
	fmt.Println("Um hacker tentando 1 milhao de senhas demoraria MUITO tempo!")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 7: Como fica no banco de dados
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 7: BANCO DE DADOS SEGURO ---")
	fmt.Println()
	fmt.Println("Agora nosso banco fica assim:")
	fmt.Println()
	fmt.Println("| email              | senha                            |")
	fmt.Println("|--------------------|----------------------------------|")
	fmt.Println("| joao@email.com     | $2a$10$xK8fQ9L2nM...             |")
	fmt.Println("| maria@email.com    | $2a$10$pR7sT4uW6x...             |")
	fmt.Println()
	fmt.Println("Mesmo se o hacker invadir, ele NAO consegue descobrir as senhas!")
	fmt.Println()

	// -------------------------------------------------------------------------
	// DESAFIOS
	// -------------------------------------------------------------------------
	fmt.Println("========================================")
	fmt.Println("   DESAFIOS PARA VOCE:")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("1. Crie uma funcao 'criarUsuario(email, senha string)' que:")
	fmt.Println("   - Recebe email e senha")
	fmt.Println("   - Gera hash da senha")
	fmt.Println("   - Retorna uma struct Usuario com email e hash")
	fmt.Println()
	fmt.Println("2. Crie uma funcao 'verificarSenha(usuario Usuario, senha string) bool'")
	fmt.Println("   que retorna true se a senha estiver correta")
	fmt.Println()
	fmt.Println("3. Teste com custo 4 (minimo) e custo 14 para ver a diferenca de tempo")
}
