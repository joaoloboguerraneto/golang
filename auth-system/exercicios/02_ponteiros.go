// =============================================================================
// EXERCICIO 2: ENTENDENDO PONTEIROS
// =============================================================================
// Execute: go run 02_ponteiros.go
//
// OBJETIVO: Entender o que sao ponteiros e quando usar
//
// CONCEITO:
// Ponteiro e uma variavel que guarda o ENDERECO de outra variavel.
// Em vez de copiar o valor, voce "aponta" para onde ele esta.
//
// ANALOGIA:
// Imagine que voce tem um livro. Voce pode:
// - Fazer uma COPIA do livro (caro, demorado)
// - Anotar ONDE o livro esta na biblioteca (rapido, economico)
// O ponteiro e como anotar a localizacao.
// =============================================================================

package main

import "fmt"

func main() {
	fmt.Println("========================================")
	fmt.Println("   EXERCICIO 2: ENTENDENDO PONTEIROS")
	fmt.Println("========================================")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 1: Criando um ponteiro
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 1: Criando um ponteiro ---")

	numero := 42 // Variavel normal

	// & = "me da o ENDERECO de"
	ponteiro := &numero

	fmt.Println("Valor de 'numero':", numero)
	fmt.Println("Endereco de 'numero' (ponteiro):", ponteiro)
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 2: Acessando valor atraves do ponteiro
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 2: Acessando valor pelo ponteiro ---")

	// * = "va no endereco e pegue o valor"
	valorNoPonteiro := *ponteiro

	fmt.Println("Valor acessado pelo ponteiro:", valorNoPonteiro)
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 3: Modificando valor atraves do ponteiro
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 3: Modificando pelo ponteiro ---")

	fmt.Println("Valor original:", numero)

	*ponteiro = 100 // Modifica o valor ORIGINAL atraves do ponteiro

	fmt.Println("Valor apos modificar pelo ponteiro:", numero)
	fmt.Println("(Perceba que 'numero' mudou, nao criamos uma copia!)")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 4: Por que isso e util? Funcoes!
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 4: Ponteiros em funcoes ---")

	idade := 25
	fmt.Println("Idade antes de chamar funcao:", idade)

	// SEM ponteiro: a funcao recebe uma COPIA
	tentarMudarSemPonteiro(idade)
	fmt.Println("Idade depois de tentarMudarSemPonteiro:", idade, "(nao mudou!)")

	// COM ponteiro: a funcao recebe o ENDERECO
	mudarComPonteiro(&idade)
	fmt.Println("Idade depois de mudarComPonteiro:", idade, "(mudou!)")
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 5: Ponteiros com structs
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 5: Ponteiros com structs ---")

	type Usuario struct {
		Nome string
		Idade int
	}

	usuario := Usuario{Nome: "Joao", Idade: 30}
	fmt.Println("Usuario antes:", usuario)

	// Passando ponteiro para funcao
	fazerAniversario(&usuario)
	fmt.Println("Usuario depois:", usuario)
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 6: nil - O ponteiro vazio
	// -------------------------------------------------------------------------
	fmt.Println("--- PARTE 6: Ponteiro nil (vazio) ---")

	var ponteiroVazio *int // Ponteiro que nao aponta para nada

	fmt.Println("Ponteiro vazio:", ponteiroVazio)
	fmt.Println("E nil?", ponteiroVazio == nil)

	// CUIDADO: Acessar um ponteiro nil causa ERRO!
	// fmt.Println(*ponteiroVazio) // PANIC!

	// Sempre verifique antes de usar:
	if ponteiroVazio != nil {
		fmt.Println("Valor:", *ponteiroVazio)
	} else {
		fmt.Println("Ponteiro esta vazio, nao posso acessar")
	}
	fmt.Println()

	// -------------------------------------------------------------------------
	// PARTE 7: Quando usar ponteiros?
	// -------------------------------------------------------------------------
	fmt.Println("========================================")
	fmt.Println("   QUANDO USAR PONTEIROS?")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("1. Quando a funcao precisa MODIFICAR o valor original")
	fmt.Println("2. Quando a struct e GRANDE (evita copiar muitos dados)")
	fmt.Println("3. Quando o valor pode ser NULO (nil)")
	fmt.Println()
	fmt.Println("NO NOSSO PROJETO:")
	fmt.Println("- *User: funcoes retornam ponteiro de usuario")
	fmt.Println("- &cookie: passamos endereco do cookie")
	fmt.Println("- *time.Time: data que pode estar vazia (BlockedAt)")
	fmt.Println()

	// -------------------------------------------------------------------------
	// DESAFIOS
	// -------------------------------------------------------------------------
	fmt.Println("========================================")
	fmt.Println("   DESAFIOS PARA VOCE:")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("1. Crie uma funcao 'dobrar' que recebe *int e dobra o valor")
	fmt.Println()
	fmt.Println("2. Crie uma struct Conta com saldo, e uma funcao 'depositar'")
	fmt.Println("   que recebe *Conta e adiciona valor ao saldo")
	fmt.Println()
	fmt.Println("3. Por que no projeto usamos *User em vez de User?")
}

// Funcao SEM ponteiro - recebe COPIA
func tentarMudarSemPonteiro(valor int) {
	valor = 999 // Muda apenas a copia local
}

// Funcao COM ponteiro - recebe ENDERECO
func mudarComPonteiro(valor *int) {
	*valor = 999 // Muda o valor original
}

// Funcao que modifica struct via ponteiro
func fazerAniversario(u *Usuario) {
	u.Idade++ // Incrementa idade no usuario ORIGINAL
}

// Struct auxiliar para a funcao acima
type Usuario struct {
	Nome  string
	Idade int
}
