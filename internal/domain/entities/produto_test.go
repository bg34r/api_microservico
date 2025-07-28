package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProdutoNew_Success(t *testing.T) {
	tests := []struct {
		name           string
		nomeProduto    string
		categoria      string
		descricao      string
		preco          float32
		expectedResult *Produto
	}{
		{
			name:        "Criar Lanche válido",
			nomeProduto: "Big Mac",
			categoria:   "Lanche",
			descricao:   "Hamburger com dois hambúrgueres",
			preco:       25.90,
			expectedResult: &Produto{
				Nome:      "Big Mac",
				Categoria: Lanche,
				Descricao: "Hamburger com dois hambúrgueres",
				Preco:     25.90,
			},
		},
		{
			name:        "Criar Bebida válida",
			nomeProduto: "Coca Cola",
			categoria:   "Bebida",
			descricao:   "Refrigerante de cola 350ml",
			preco:       5.50,
			expectedResult: &Produto{
				Nome:      "Coca Cola",
				Categoria: Bebida,
				Descricao: "Refrigerante de cola 350ml",
				Preco:     5.50,
			},
		},
		{
			name:        "Criar Acompanhamento válido",
			nomeProduto: "Batata Frita",
			categoria:   "Acompanhamento",
			descricao:   "Batata frita crocante",
			preco:       8.00,
			expectedResult: &Produto{
				Nome:      "Batata Frita",
				Categoria: Acompanhamento,
				Descricao: "Batata frita crocante",
				Preco:     8.00,
			},
		},
		{
			name:        "Criar Sobremesa válida",
			nomeProduto: "Sundae Chocolate",
			categoria:   "Sobremesa",
			descricao:   "Sorvete com calda de chocolate",
			preco:       12.50,
			expectedResult: &Produto{
				Nome:      "Sundae Chocolate",
				Categoria: Sobremesa,
				Descricao: "Sorvete com calda de chocolate",
				Preco:     12.50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			produto, err := ProdutoNew(tt.nomeProduto, tt.categoria, tt.descricao, tt.preco)

			assert.NoError(t, err)
			assert.NotNil(t, produto)
			assert.Equal(t, tt.expectedResult.Nome, produto.Nome)
			assert.Equal(t, tt.expectedResult.Categoria, produto.Categoria)
			assert.Equal(t, tt.expectedResult.Descricao, produto.Descricao)
			assert.Equal(t, tt.expectedResult.Preco, produto.Preco)
			assert.Zero(t, produto.ID) // ID deve ser zero para novos produtos
		})
	}
}

func TestProdutoNew_ErrorNomeVazio(t *testing.T) {
	tests := []struct {
		name        string
		nomeProduto string
	}{
		{"Nome vazio", ""},
		{"Nome apenas espaços", "   "},
		{"Nome com tabs", "\t\t"},
		{"Nome com quebras de linha", "\n\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			produto, err := ProdutoNew(tt.nomeProduto, "Lanche", "Descrição válida", 10.0)

			assert.Error(t, err)
			assert.Nil(t, produto)
			assert.Equal(t, "todos os campos são obrigatórios e o preço maior que zero", err.Error())
		})
	}
}

func TestProdutoNew_ErrorCategoriaVazia(t *testing.T) {
	tests := []struct {
		name      string
		categoria string
	}{
		{"Categoria vazia", ""},
		{"Categoria apenas espaços", "   "},
		{"Categoria com tabs", "\t\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			produto, err := ProdutoNew("Produto Teste", tt.categoria, "Descrição válida", 10.0)

			assert.Error(t, err)
			assert.Nil(t, produto)
			assert.Equal(t, "todos os campos são obrigatórios e o preço maior que zero", err.Error())
		})
	}
}

func TestProdutoNew_ErrorPrecoInvalido(t *testing.T) {
	tests := []struct {
		name  string
		preco float32
	}{
		{"Preço zero", 0.0},
		{"Preço negativo", -1.0},
		{"Preço muito negativo", -100.50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			produto, err := ProdutoNew("Produto Teste", "Lanche", "Descrição válida", tt.preco)

			assert.Error(t, err)
			assert.Nil(t, produto)
			assert.Equal(t, "todos os campos são obrigatórios e o preço maior que zero", err.Error())
		})
	}
}

func TestProdutoNew_ErrorCategoriaInvalida(t *testing.T) {
	tests := []struct {
		name      string
		categoria string
	}{
		{"Categoria inexistente", "Categoria Inexistente"},
		{"Case sensitive - lanche minúsculo", "lanche"},
		{"Case sensitive - LANCHE maiúsculo", "LANCHE"},
		{"Categoria com caracteres especiais", "Lanche@#$"},
		{"Categoria numérica", "123"},
		{"Categoria mista", "Lanche123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			produto, err := ProdutoNew("Produto Teste", tt.categoria, "Descrição válida", 10.0)

			assert.Error(t, err)
			assert.Nil(t, produto)
			assert.Equal(t, "categoria inválida", err.Error())
		})
	}
}

func TestProdutoNew_ErrorMultiplosCampos(t *testing.T) {
	// Teste quando múltiplos campos são inválidos
	produto, err := ProdutoNew("", "", "", 0.0)

	assert.Error(t, err)
	assert.Nil(t, produto)
	assert.Equal(t, "todos os campos são obrigatórios e o preço maior que zero", err.Error())
}

func TestProdutoNew_ValidarConstantesCategorias(t *testing.T) {
	// Teste para garantir que as constantes estão corretas
	assert.Equal(t, CatProduto("Lanche"), Lanche)
	assert.Equal(t, CatProduto("Acompanhamento"), Acompanhamento)
	assert.Equal(t, CatProduto("Bebida"), Bebida)
	assert.Equal(t, CatProduto("Sobremesa"), Sobremesa)
}

func TestProdutoNew_DescricaoVazia(t *testing.T) {
	// Teste se descrição vazia é aceita (baseado no código atual)
	produto, err := ProdutoNew("Produto Teste", "Lanche", "", 10.0)

	assert.NoError(t, err)
	assert.NotNil(t, produto)
	assert.Equal(t, "Produto Teste", produto.Nome)
	assert.Equal(t, Lanche, produto.Categoria)
	assert.Equal(t, "", produto.Descricao)
	assert.Equal(t, float32(10.0), produto.Preco)
}

func TestProdutoNew_PrecoDecimal(t *testing.T) {
	// Teste com valores decimais precisos
	produto, err := ProdutoNew("Produto Teste", "Lanche", "Descrição", 19.99)

	assert.NoError(t, err)
	assert.NotNil(t, produto)
	assert.Equal(t, float32(19.99), produto.Preco)
}

func TestProdutoNew_NomeComEspacos(t *testing.T) {
	// Teste com nome que tem espaços mas não é vazio
	produto, err := ProdutoNew("  Big Mac  ", "Lanche", "Hamburger", 25.90)

	assert.NoError(t, err)
	assert.NotNil(t, produto)
	assert.Equal(t, "  Big Mac  ", produto.Nome) // Mantém os espaços extras
}

// Benchmark para testar performance
func BenchmarkProdutoNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ProdutoNew("Big Mac", "Lanche", "Hamburger com dois hambúrgueres", 25.90)
	}
}
