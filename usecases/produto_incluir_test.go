package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"testing"
)

// MockProdutoRepositoryIncluir implements repository.ProdutoRepository for testing
type MockProdutoRepositoryIncluir struct {
	Produtos []*entities.Produto
}

func (m *MockProdutoRepositoryIncluir) AdicionarProduto(ctx context.Context, produto *entities.Produto) error {
	produto.ID = len(m.Produtos) + 1
	m.Produtos = append(m.Produtos, produto)
	return nil
}

func (m *MockProdutoRepositoryIncluir) EditarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryIncluir) RemoverProduto(ctx context.Context, id int) error {
	return nil
}

func (m *MockProdutoRepositoryIncluir) BuscarProdutoPorId(ctx context.Context, id int) (*entities.Produto, error) {
	for _, produto := range m.Produtos {
		if produto.ID == id {
			return produto, nil
		}
	}
	return nil, nil
}

func (m *MockProdutoRepositoryIncluir) ListarTodosOsProdutos(ctx context.Context) ([]*entities.Produto, error) {
	return m.Produtos, nil
}

func (m *MockProdutoRepositoryIncluir) ListarPorCategoria(ctx context.Context, categoria string) ([]*entities.Produto, error) {
	var produtosFiltrados []*entities.Produto
	for _, produto := range m.Produtos {
		if string(produto.Categoria) == categoria {
			produtosFiltrados = append(produtosFiltrados, produto)
		}
	}
	return produtosFiltrados, nil
}

func TestProdutoIncluir_Run_Sucesso(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryIncluir{
		Produtos: []*entities.Produto{},
	}

	useCase := NewProdutoIncluirUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx, "Hamburguer", "Lanche", "Delicioso hamburguer", 25.0)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		t.Fatal("Esperado produto criado, recebido nil")
	}

	if resultado.Nome != "Hamburguer" {
		t.Errorf("Esperado nome 'Hamburguer', recebido %s", resultado.Nome)
	}

	if resultado.Categoria != entities.Lanche {
		t.Errorf("Esperado categoria 'Lanche', recebido %s", resultado.Categoria)
	}

	if resultado.Preco != 25.0 {
		t.Errorf("Esperado preço 25.0, recebido %f", resultado.Preco)
	}

	if len(mockRepo.Produtos) != 1 {
		t.Errorf("Esperado 1 produto no repositório, encontrado %d", len(mockRepo.Produtos))
	}
}

func TestProdutoIncluir_Run_DadosInvalidos(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryIncluir{
		Produtos: []*entities.Produto{},
	}

	useCase := NewProdutoIncluirUseCase(mockRepo)

	ctx := context.Background()

	// Test cases com dados inválidos
	testCases := []struct {
		nome      string
		categoria string
		descricao string
		preco     float32
		descTest  string
	}{
		{"", "Lanche", "Descrição", 10.0, "nome vazio"},
		{"Produto", "", "Descrição", 10.0, "categoria vazia"},
		{"Produto", "Lanche", "Descrição", 0, "preço zero"},
		{"Produto", "Lanche", "Descrição", -5.0, "preço negativo"},
		{"Produto", "CategoriaInvalida", "Descrição", 10.0, "categoria inválida"},
	}

	for _, tc := range testCases {
		// When
		resultado, err := useCase.Run(ctx, tc.nome, tc.categoria, tc.descricao, tc.preco)

		// Then
		if err == nil {
			t.Errorf("Esperado erro para %s, recebido nil", tc.descTest)
		}

		if resultado != nil {
			t.Errorf("Esperado nil para %s, recebido produto", tc.descTest)
		}
	}
}

func TestProdutoIncluir_Run_MultiplosProdutos(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryIncluir{
		Produtos: []*entities.Produto{},
	}

	useCase := NewProdutoIncluirUseCase(mockRepo)

	ctx := context.Background()

	produtos := []struct {
		nome      string
		categoria string
		descricao string
		preco     float32
	}{
		{"Hamburguer", "Lanche", "Delicioso hamburguer", 25.0},
		{"Batata Frita", "Acompanhamento", "Batata crocante", 10.0},
		{"Coca-Cola", "Bebida", "Refrigerante gelado", 7.5},
		{"Sorvete", "Sobremesa", "Sorvete cremoso", 8.0},
	}

	// When - criando múltiplos produtos
	for i, p := range produtos {
		resultado, err := useCase.Run(ctx, p.nome, p.categoria, p.descricao, p.preco)

		// Then
		if err != nil {
			t.Errorf("Não esperado erro para produto %d, recebido %v", i+1, err)
		}

		if resultado == nil {
			t.Errorf("Esperado produto criado para produto %d, recebido nil", i+1)
		}

		if resultado.Nome != p.nome {
			t.Errorf("Esperado nome %s, recebido %s", p.nome, resultado.Nome)
		}
	}

	// Verificar se todos foram adicionados
	if len(mockRepo.Produtos) != len(produtos) {
		t.Errorf("Esperado %d produtos no repositório, encontrado %d", len(produtos), len(mockRepo.Produtos))
	}
}

func TestProdutoIncluir_Run_TodasCategorias(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryIncluir{
		Produtos: []*entities.Produto{},
	}

	useCase := NewProdutoIncluirUseCase(mockRepo)

	ctx := context.Background()

	categorias := []string{"Lanche", "Acompanhamento", "Bebida", "Sobremesa"}

	// When - testando cada categoria válida
	for _, categoria := range categorias {
		resultado, err := useCase.Run(ctx, "Produto "+categoria, categoria, "Descrição", 10.0)

		// Then
		if err != nil {
			t.Errorf("Não esperado erro para categoria %s, recebido %v", categoria, err)
		}

		if resultado == nil {
			t.Errorf("Esperado produto criado para categoria %s, recebido nil", categoria)
		}

		if string(resultado.Categoria) != categoria {
			t.Errorf("Esperado categoria %s, recebido %s", categoria, resultado.Categoria)
		}
	}

	if len(mockRepo.Produtos) != len(categorias) {
		t.Errorf("Esperado %d produtos criados, encontrado %d", len(categorias), len(mockRepo.Produtos))
	}
}
