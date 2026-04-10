// =============================================================================
// EXERCICIO 1: ENTENDENDO STRUCTS
// =============================================================================
// Execute: go run 01_structs.go
//
// OBJETIVO: Entender como criar e usar structs em Go
//
// CONCEITO:
// Struct e como uma FICHA que define quais informacoes guardar juntas.
// E similar a uma classe em outras linguagens, mas mais simples.
// =============================================================================

package main

import "fmt"

// -----------------------------------------------------------------------------
// DEFININDO UMA STRUCT
// -----------------------------------------------------------------------------
// Aqui estamos criando um "modelo" de Usuario
// Ainda nao e um usuario real, e so a DEFINICAO de como um usuario deve ser

type Usuario struct {
	Nome  string // Campo do tipo texto
	Email string // Campo do tipo texto
	Idade int    // Campo do tipo numero inteiro
	Ativo bool   // Campo do tipo verdadeiro/falso
}

// -----------------------------------------------------------------------------
// OUTRA STRUCT: Produto
// -----------------------------------------------------------------------------

type Produto struct {
	ID        int
	Nome      string
	Preco     float64 // Numero com casas decimais
	Estoque   int
	Descricao string
}

func main() {
	fmt.Println("========================================")
	fmt.Println("   EXERCICIO 1: ENTENDENDO STRUCTS")
	fmt.Println("========================================")
	fmt.Println()

	// -------------------------------------------------------------------------
	// CRIANDO UM USUARIO (preenchendo a ficha)
	// -------------------------------------------------------------------------
	fmt.Println("--- Criando um usuario ---")

	// Forma 1: Especificando os campos por nome
	usuario1 := Usuario{
		Nome:  "Joao Silva",
		Email: "joao@email.com",
		Idade: 25,
		Ativo: true,
	}

	// Forma 2: Na ordem dos campos (menos recomendado)
	usuario2 := Usuario{"Maria Santos", "maria@email.com", 30, false}

	// Forma 3: Criando vazio e preenchendo depois
	var usuario3 Usuario
	usuario3.Nome = "Pedro Costa"
	usuario3.Email = "pedro@email.com"
	usuario3.Idade = 28
	usuario3.Ativo = true

	// -------------------------------------------------------------------------
	// ACESSANDO OS CAMPOS
	// -------------------------------------------------------------------------
	fmt.Println()
	fmt.Println("--- Acessando campos ---")

	fmt.Println("Usuario 1:")
	fmt.Println("  Nome:", usuario1.Nome)
	fmt.Println("  Email:", usuario1.Email)
	fmt.Println("  Idade:", usuario1.Idade)
	fmt.Println("  Ativo:", usuario1.Ativo)

	fmt.Println()
	fmt.Println("Usuario 2:")
	fmt.Println("  Nome:", usuario2.Nome)
	fmt.Println("  Ativo:", usuario2.Ativo)

	// -------------------------------------------------------------------------
	// MODIFICANDO CAMPOS
	// -------------------------------------------------------------------------
	fmt.Println()
	fmt.Println("--- Modificando campos ---")

	fmt.Println("Idade do usuario3 antes:", usuario3.Idade)
	usuario3.Idade = 29 // Aniversario!
	fmt.Println("Idade do usuario3 depois:", usuario3.Idade)

	// -------------------------------------------------------------------------
	// STRUCT PRODUTO
	// -------------------------------------------------------------------------
	fmt.Println()
	fmt.Println("--- Criando um produto ---")

	produto := Produto{
		ID:        1,
		Nome:      "Notebook Dell",
		Preco:     3500.99,
		Estoque:   15,
		Descricao: "Notebook Dell Inspiron 15",
	}

	fmt.Printf("Produto: %s\n", produto.Nome)
	fmt.Printf("Preco: R$ %.2f\n", produto.Preco)
	fmt.Printf("Estoque: %d unidades\n", produto.Estoque)

	// -------------------------------------------------------------------------
	// COMPARANDO STRUCTS
	// -------------------------------------------------------------------------
	fmt.Println()
	fmt.Println("--- Comparando structs ---")

	a := Usuario{Nome: "Ana", Email: "ana@email.com", Idade: 20, Ativo: true}
	b := Usuario{Nome: "Ana", Email: "ana@email.com", Idade: 20, Ativo: true}
	c := Usuario{Nome: "Ana", Email: "ana@email.com", Idade: 21, Ativo: true}

	fmt.Println("a == b?", a == b) // true (todos campos iguais)
	fmt.Println("a == c?", a == c) // false (idade diferente)

	// -------------------------------------------------------------------------
	// DESAFIOS
	// -------------------------------------------------------------------------
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("   DESAFIOS PARA VOCE:")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("1. Crie uma struct 'Endereco' com campos:")
	fmt.Println("   - Rua (string)")
	fmt.Println("   - Numero (int)")
	fmt.Println("   - Cidade (string)")
	fmt.Println("   - CEP (string)")
	fmt.Println()
	fmt.Println("2. Crie uma variavel do tipo Endereco e preencha")
	fmt.Println()
	fmt.Println("3. Adicione um campo 'Endereco' na struct Usuario")
	fmt.Println("   (sim, uma struct dentro de outra!)")
}
