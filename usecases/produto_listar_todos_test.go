package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"testing"
)

// MockProdutoRepositoryListarTodos implements repository.ProdutoRepository for testing
type MockProdutoRepositoryListarTodos struct {
	Produtos []*entities.Produto
}

func (m *MockProdutoRepositoryListarTodos) AdicionarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryListarTodos) EditarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryListarTodos) RemoverProduto(ctx context.Context, id int) error {
	return nil
}

func (m *MockProdutoRepositoryListarTodos) BuscarProdutoPorId(ctx context.Context, id int) (*entities.Produto, error) {
	for _, produto := range m.Produtos {
		if produto.ID == id {
			return produto, nil
		}
	}
	return nil, nil
}

func (m *MockProdutoRepositoryListarTodos) ListarTodosOsProdutos(ctx context.Context) ([]*entities.Produto, error) {
	return m.Produtos, nil
}

func (m *MockProdutoRepositoryListarTodos) ListarPorCategoria(ctx context.Context, categoria string) ([]*entities.Produto, error) {
	var produtosFiltrados []*entities.Produto
	for _, produto := range m.Produtos {
		if string(produto.Categoria) == categoria {
			produtosFiltrados = append(produtosFiltrados, produto)
		}
	}
	return produtosFiltrados, nil
}

func TestProdutoListarTodos_Run_Sucesso(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
		{ID: 2, Nome: "Batata Frita", Categoria: entities.Acompanhamento, Descricao: "Acompanhamento crocante", Preco: 10.0},
		{ID: 3, Nome: "Coca-Cola", Categoria: entities.Bebida, Descricao: "Bebida refrescante", Preco: 7.5},
	}

	mockRepo := &MockProdutoRepositoryListarTodos{
		Produtos: produtos,
	}

	useCase := NewProdutoListarTodosUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		t.Fatal("Esperado lista de produtos, recebido nil")
	}

	if len(resultado) != len(produtos) {
		t.Errorf("Esperado %d produtos, recebido %d", len(produtos), len(resultado))
	}

	// Verificar se todos os produtos estão presentes
	for i, produtoEsperado := range produtos {
		if resultado[i].ID != produtoEsperado.ID {
			t.Errorf("Esperado ID %d na posição %d, recebido %d", produtoEsperado.ID, i, resultado[i].ID)
		}
		if resultado[i].Nome != produtoEsperado.Nome {
			t.Errorf("Esperado nome %s na posição %d, recebido %s", produtoEsperado.Nome, i, resultado[i].Nome)
		}
	}
}

func TestProdutoListarTodos_Run_ListaVazia(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryListarTodos{
		Produtos: []*entities.Produto{},
	}

	useCase := NewProdutoListarTodosUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		t.Fatal("Esperado lista vazia, recebido nil")
	}

	if len(resultado) != 0 {
		t.Errorf("Esperado lista vazia, recebido %d produtos", len(resultado))
	}
}

func TestProdutoListarTodos_Run_UmProduto(t *testing.T) {
	// Given
	produto := &entities.Produto{
		ID:        1,
		Nome:      "Produto Único",
		Categoria: entities.Lanche,
		Descricao: "Único produto",
		Preco:     15.0,
	}

	mockRepo := &MockProdutoRepositoryListarTodos{
		Produtos: []*entities.Produto{produto},
	}

	useCase := NewProdutoListarTodosUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		t.Fatal("Esperado lista com um produto, recebido nil")
	}

	if len(resultado) != 1 {
		t.Errorf("Esperado 1 produto, recebido %d", len(resultado))
	}

	if resultado[0].ID != 1 {
		t.Errorf("Esperado ID 1, recebido %d", resultado[0].ID)
	}

	if resultado[0].Nome != "Produto Único" {
		t.Errorf("Esperado nome 'Produto Único', recebido %s", resultado[0].Nome)
	}
}

func TestProdutoListarTodos_Run_TodasCategorias(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
		{ID: 2, Nome: "Batata Frita", Categoria: entities.Acompanhamento, Descricao: "Acompanhamento crocante", Preco: 10.0},
		{ID: 3, Nome: "Coca-Cola", Categoria: entities.Bebida, Descricao: "Bebida refrescante", Preco: 7.5},
		{ID: 4, Nome: "Sorvete", Categoria: entities.Sobremesa, Descricao: "Sobremesa gelada", Preco: 8.0},
	}

	mockRepo := &MockProdutoRepositoryListarTodos{
		Produtos: produtos,
	}

	useCase := NewProdutoListarTodosUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if len(resultado) != 4 {
		t.Errorf("Esperado 4 produtos, recebido %d", len(resultado))
	}

	// Verificar se todas as categorias estão presentes
	categorias := make(map[entities.CatProduto]bool)
	for _, produto := range resultado {
		categorias[produto.Categoria] = true
	}

	categoriasEsperadas := []entities.CatProduto{
		entities.Lanche,
		entities.Acompanhamento,
		entities.Bebida,
		entities.Sobremesa,
	}

	for _, categoria := range categoriasEsperadas {
		if !categorias[categoria] {
			t.Errorf("Categoria %s não encontrada na lista", categoria)
		}
	}
}

func TestProdutoListarTodos_Run_MuitosProdutos(t *testing.T) {
	// Given - criando muitos produtos para testar performance
	var produtos []*entities.Produto
	for i := 1; i <= 100; i++ {
		categoria := entities.Lanche
		if i%4 == 0 {
			categoria = entities.Acompanhamento
		} else if i%4 == 1 {
			categoria = entities.Bebida
		} else if i%4 == 2 {
			categoria = entities.Sobremesa
		}

		produto := &entities.Produto{
			ID:        i,
			Nome:      fmt.Sprintf("Produto %d", i),
			Categoria: categoria,
			Descricao: fmt.Sprintf("Descrição do produto %d", i),
			Preco:     float32(10 + i),
		}
		produtos = append(produtos, produto)
	}

	mockRepo := &MockProdutoRepositoryListarTodos{
		Produtos: produtos,
	}

	useCase := NewProdutoListarTodosUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if len(resultado) != 100 {
		t.Errorf("Esperado 100 produtos, recebido %d", len(resultado))
	}

	// Verificar ordem
	for i, produto := range resultado {
		expectedID := i + 1
		if produto.ID != expectedID {
			t.Errorf("Esperado ID %d na posição %d, recebido %d", expectedID, i, produto.ID)
		}
	}
}
