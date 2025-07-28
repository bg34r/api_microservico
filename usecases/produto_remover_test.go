package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"strings"
	"testing"
)

// MockProdutoRepositoryRemover implements repository.ProdutoRepository for testing
type MockProdutoRepositoryRemover struct {
	Produtos []*entities.Produto
}

type MockEventPublisherRemover struct{}

func (m *MockEventPublisherRemover) Publish(eventType string, payload interface{}) error {
	// apenas retorna nil, simula sucesso
	return nil
}

func (m *MockProdutoRepositoryRemover) AdicionarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryRemover) EditarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryRemover) RemoverProduto(ctx context.Context, id int) error {
	for i, produto := range m.Produtos {
		if produto.ID == id {
			m.Produtos = append(m.Produtos[:i], m.Produtos[i+1:]...)
			return nil
		}
	}
	return errors.New("produto não encontrado")
}

func (m *MockProdutoRepositoryRemover) BuscarProdutoPorId(ctx context.Context, id int) (*entities.Produto, error) {
	for _, produto := range m.Produtos {
		if produto.ID == id {
			return produto, nil
		}
	}
	return nil, errors.New("produto não encontrado")
}

func (m *MockProdutoRepositoryRemover) ListarTodosOsProdutos(ctx context.Context) ([]*entities.Produto, error) {
	return m.Produtos, nil
}

func (m *MockProdutoRepositoryRemover) ListarPorCategoria(ctx context.Context, categoria string) ([]*entities.Produto, error) {
	var produtosFiltrados []*entities.Produto
	for _, produto := range m.Produtos {
		if string(produto.Categoria) == categoria {
			produtosFiltrados = append(produtosFiltrados, produto)
		}
	}
	return produtosFiltrados, nil
}

func TestProdutoRemover_Run_Sucesso(t *testing.T) {
	// Given
	produto := &entities.Produto{
		ID:        1,
		Nome:      "Produto Teste",
		Categoria: entities.Lanche,
		Descricao: "Descrição teste",
		Preco:     10.0,
	}

	mockRepo := &MockProdutoRepositoryRemover{
		Produtos: []*entities.Produto{produto},
	}

	mockPublisher := &MockEventPublisherRemover{}

	useCase := NewProdutoRemoverUseCase(mockRepo, mockPublisher)

	ctx := context.Background()

	// When
	err := useCase.Run(ctx, 1)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	// Verificar se o produto foi removido
	if len(mockRepo.Produtos) != 0 {
		t.Errorf("Esperado 0 produtos após remoção, encontrado %d", len(mockRepo.Produtos))
	}
}

func TestProdutoRemover_Run_ProdutoNaoEncontrado(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryRemover{
		Produtos: []*entities.Produto{},
	}

	mockPublisher := &MockEventPublisherRemover{}

	useCase := NewProdutoRemoverUseCase(mockRepo, mockPublisher)

	ctx := context.Background()

	// When
	err := useCase.Run(ctx, 999)

	// Then
	if err == nil {
		t.Error("Esperado erro, recebido nil")
	}

	if !strings.Contains(err.Error(), "produto não existe no banco de dados") {
		t.Errorf("Esperado erro sobre produto não encontrado, recebido: %v", err)
	}
}

func TestProdutoRemover_Run_IdInvalido(t *testing.T) {
	// Given
	produto := &entities.Produto{
		ID:        1,
		Nome:      "Produto Teste",
		Categoria: entities.Lanche,
		Descricao: "Descrição teste",
		Preco:     10.0,
	}

	mockRepo := &MockProdutoRepositoryRemover{
		Produtos: []*entities.Produto{produto},
	}
	mockPublisher := &MockEventPublisherRemover{}

	useCase := NewProdutoRemoverUseCase(mockRepo, mockPublisher)

	ctx := context.Background()

	// Test com IDs inválidos
	idsInvalidos := []int{0, -1, -999}

	for _, id := range idsInvalidos {
		// When
		err := useCase.Run(ctx, id)

		// Then
		if err == nil {
			t.Errorf("Esperado erro para ID inválido %d, recebido nil", id)
		}

		// Produto original deve permanecer
		if len(mockRepo.Produtos) != 1 {
			t.Errorf("Produto original deve permanecer para ID inválido %d", id)
		}
	}
}

func TestProdutoRemover_Run_MultiplosProdutos(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
		{ID: 2, Nome: "Batata Frita", Categoria: entities.Acompanhamento, Descricao: "Acompanhamento crocante", Preco: 10.0},
		{ID: 3, Nome: "Coca-Cola", Categoria: entities.Bebida, Descricao: "Bebida refrescante", Preco: 7.5},
	}

	mockRepo := &MockProdutoRepositoryRemover{
		Produtos: produtos,
	}

	mockPublisher := &MockEventPublisherRemover{}

	useCase := NewProdutoRemoverUseCase(mockRepo, mockPublisher)

	ctx := context.Background()

	totalInicial := len(mockRepo.Produtos)

	// When - removendo produto do meio
	err := useCase.Run(ctx, 2)

	// Then
	if err != nil {
		t.Errorf("Não esperado erro, recebido %v", err)
	}

	if len(mockRepo.Produtos) != totalInicial-1 {
		t.Errorf("Esperado %d produtos após remoção, encontrado %d", totalInicial-1, len(mockRepo.Produtos))
	}

	// Verificar que os produtos corretos permanecem
	for _, produto := range mockRepo.Produtos {
		if produto.ID == 2 {
			t.Error("Produto com ID 2 deveria ter sido removido")
		}
	}

	// Verificar que produtos 1 e 3 ainda existem
	encontrouProduto1 := false
	encontrouProduto3 := false
	for _, produto := range mockRepo.Produtos {
		if produto.ID == 1 {
			encontrouProduto1 = true
		}
		if produto.ID == 3 {
			encontrouProduto3 = true
		}
	}

	if !encontrouProduto1 {
		t.Error("Produto com ID 1 deveria permanecer")
	}
	if !encontrouProduto3 {
		t.Error("Produto com ID 3 deveria permanecer")
	}
}

func TestProdutoRemover_Run_RemocaoSequencial(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Produto 1", Categoria: entities.Lanche, Descricao: "Descrição 1", Preco: 10.0},
		{ID: 2, Nome: "Produto 2", Categoria: entities.Bebida, Descricao: "Descrição 2", Preco: 15.0},
		{ID: 3, Nome: "Produto 3", Categoria: entities.Sobremesa, Descricao: "Descrição 3", Preco: 20.0},
	}

	mockRepo := &MockProdutoRepositoryRemover{
		Produtos: produtos,
	}

	mockPublisher := &MockEventPublisherRemover{}

	useCase := NewProdutoRemoverUseCase(mockRepo, mockPublisher)

	ctx := context.Background()

	// When - removendo produtos sequencialmente
	idsParaRemover := []int{1, 3, 2}

	for i, id := range idsParaRemover {
		err := useCase.Run(ctx, id)

		// Then
		if err != nil {
			t.Errorf("Não esperado erro na remoção %d, recebido %v", i+1, err)
		}

		expectedCount := len(produtos) - (i + 1)
		if len(mockRepo.Produtos) != expectedCount {
			t.Errorf("Esperado %d produtos após remoção %d, encontrado %d", expectedCount, i+1, len(mockRepo.Produtos))
		}
	}

	// Verificar que todos foram removidos
	if len(mockRepo.Produtos) != 0 {
		t.Errorf("Esperado 0 produtos após remoção completa, encontrado %d", len(mockRepo.Produtos))
	}
}
