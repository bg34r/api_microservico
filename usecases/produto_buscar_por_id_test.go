package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"strings"
	"testing"
)

// MockProdutoRepositoryBuscar implements repository.ProdutoRepository for testing
type MockProdutoRepositoryBuscar struct {
	Produtos []*entities.Produto
}

func (m *MockProdutoRepositoryBuscar) AdicionarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryBuscar) EditarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryBuscar) RemoverProduto(ctx context.Context, id int) error {
	return nil
}

func (m *MockProdutoRepositoryBuscar) BuscarProdutoPorId(ctx context.Context, id int) (*entities.Produto, error) {
	for _, produto := range m.Produtos {
		if produto.ID == id {
			return produto, nil
		}
	}
	return nil, errors.New("produto não encontrado")
}

func (m *MockProdutoRepositoryBuscar) ListarTodosOsProdutos(ctx context.Context) ([]*entities.Produto, error) {
	return m.Produtos, nil
}

func (m *MockProdutoRepositoryBuscar) ListarPorCategoria(ctx context.Context, categoria string) ([]*entities.Produto, error) {
	var produtosFiltrados []*entities.Produto
	for _, produto := range m.Produtos {
		if string(produto.Categoria) == categoria {
			produtosFiltrados = append(produtosFiltrados, produto)
		}
	}
	return produtosFiltrados, nil
}

func TestProdutoBuscarPorId_Run_Sucesso(t *testing.T) {
	// Given
	produto := &entities.Produto{
		ID:        1,
		Nome:      "Produto Teste",
		Categoria: entities.Bebida,
		Descricao: "Descrição teste",
		Preco:     10.0,
	}

	mockRepo := &MockProdutoRepositoryBuscar{
		Produtos: []*entities.Produto{produto},
	}

	useCase := NewProdutoBuscaPorIdUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx, 1)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		t.Fatal("Esperado produto, recebido nil")
	}

	if resultado.ID != 1 {
		t.Errorf("Esperado ID 1, recebido %d", resultado.ID)
	}

	if resultado.Nome != "Produto Teste" {
		t.Errorf("Esperado nome 'Produto Teste', recebido %s", resultado.Nome)
	}
}

func TestProdutoBuscarPorId_Run_ProdutoNaoEncontrado(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryBuscar{
		Produtos: []*entities.Produto{},
	}

	useCase := NewProdutoBuscaPorIdUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx, 1)

	// Then
	if err == nil {
		t.Error("Esperado erro, recebido nil")
	}

	if resultado != nil {
		t.Error("Esperado nil, recebido produto")
	}

	if !strings.Contains(err.Error(), "não foi possível buscar produto") {
		t.Errorf("Esperado erro sobre produto não encontrado, recebido: %v", err)
	}
}

func TestProdutoBuscarPorId_Run_IdInvalido(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryBuscar{}

	useCase := NewProdutoBuscaPorIdUseCase(mockRepo)

	ctx := context.Background()

	// Test com IDs inválidos
	idsInvalidos := []int{0, -1, -999}

	for _, id := range idsInvalidos {
		// When
		resultado, err := useCase.Run(ctx, id)

		// Then
		if err == nil {
			t.Errorf("Esperado erro para ID inválido %d, recebido nil", id)
		}
		if resultado != nil {
			t.Errorf("Esperado nil para ID inválido %d, recebido %+v", id, resultado)
		}
	}
}

func TestProdutoBuscarPorId_Run_MultiplosProdutos(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
		{ID: 2, Nome: "Batata Frita", Categoria: entities.Acompanhamento, Descricao: "Acompanhamento crocante", Preco: 10.0},
		{ID: 3, Nome: "Coca-Cola", Categoria: entities.Bebida, Descricao: "Bebida refrescante", Preco: 7.5},
	}

	mockRepo := &MockProdutoRepositoryBuscar{
		Produtos: produtos,
	}

	useCase := NewProdutoBuscaPorIdUseCase(mockRepo)

	ctx := context.Background()

	// Test buscar cada produto
	for _, produtoEsperado := range produtos {
		// When
		resultado, err := useCase.Run(ctx, produtoEsperado.ID)

		// Then
		if err != nil {
			t.Fatalf("Não esperado erro para ID %d, recebido %v", produtoEsperado.ID, err)
		}
		if resultado.ID != produtoEsperado.ID {
			t.Errorf("Esperado ID %d, recebido %d", produtoEsperado.ID, resultado.ID)
		}
		if resultado.Nome != produtoEsperado.Nome {
			t.Errorf("Esperado nome %s, recebido %s", produtoEsperado.Nome, resultado.Nome)
		}
	}
}
