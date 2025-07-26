package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"testing"
)

// MockProdutoRepositoryListarPorCategoria implements repository.ProdutoRepository for testing
type MockProdutoRepositoryListarPorCategoria struct {
	Produtos []*entities.Produto
}

func (m *MockProdutoRepositoryListarPorCategoria) AdicionarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryListarPorCategoria) EditarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryListarPorCategoria) RemoverProduto(ctx context.Context, id int) error {
	return nil
}

func (m *MockProdutoRepositoryListarPorCategoria) BuscarProdutoPorId(ctx context.Context, id int) (*entities.Produto, error) {
	for _, produto := range m.Produtos {
		if produto.ID == id {
			return produto, nil
		}
	}
	return nil, nil
}

func (m *MockProdutoRepositoryListarPorCategoria) ListarTodosOsProdutos(ctx context.Context) ([]*entities.Produto, error) {
	return m.Produtos, nil
}

func (m *MockProdutoRepositoryListarPorCategoria) ListarPorCategoria(ctx context.Context, categoria string) ([]*entities.Produto, error) {
	var produtosFiltrados []*entities.Produto
	for _, produto := range m.Produtos {
		if string(produto.Categoria) == categoria {
			produtosFiltrados = append(produtosFiltrados, produto)
		}
	}
	return produtosFiltrados, nil
}

func TestProdutoListarPorCategoria_Run_Sucesso(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
		{ID: 2, Nome: "Cheeseburger", Categoria: entities.Lanche, Descricao: "Lanche com queijo", Preco: 28.0},
		{ID: 3, Nome: "Batata Frita", Categoria: entities.Acompanhamento, Descricao: "Acompanhamento crocante", Preco: 10.0},
		{ID: 4, Nome: "Coca-Cola", Categoria: entities.Bebida, Descricao: "Bebida refrescante", Preco: 7.5},
	}

	mockRepo := &MockProdutoRepositoryListarPorCategoria{
		Produtos: produtos,
	}

	useCase := NewProdutoListarPorCategoriaUseCase(mockRepo)

	ctx := context.Background()

	// When - buscando produtos da categoria Lanche
	resultado, err := useCase.Run(ctx, "Lanche")

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		t.Fatal("Esperado lista de produtos, recebido nil")
	}

	if len(resultado) != 2 {
		t.Errorf("Esperado 2 produtos da categoria Lanche, recebido %d", len(resultado))
	}

	// Verificar se todos os produtos são da categoria correta
	for _, produto := range resultado {
		if produto.Categoria != entities.Lanche {
			t.Errorf("Esperado categoria Lanche, recebido %s", produto.Categoria)
		}
	}

	// Verificar produtos específicos
	nomes := []string{resultado[0].Nome, resultado[1].Nome}
	if !contains(nomes, "Hamburguer") || !contains(nomes, "Cheeseburger") {
		t.Error("Produtos esperados (Hamburguer, Cheeseburger) não encontrados")
	}
}

func TestProdutoListarPorCategoria_Run_CategoriaVazia(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
		{ID: 2, Nome: "Coca-Cola", Categoria: entities.Bebida, Descricao: "Bebida refrescante", Preco: 7.5},
	}

	mockRepo := &MockProdutoRepositoryListarPorCategoria{
		Produtos: produtos,
	}

	useCase := NewProdutoListarPorCategoriaUseCase(mockRepo)

	ctx := context.Background()

	// When - buscando categoria que não tem produtos
	resultado, err := useCase.Run(ctx, "Sobremesa")

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		resultado = []*entities.Produto{} // Tratar nil como lista vazia
	}

	if len(resultado) != 0 {
		t.Errorf("Esperado 0 produtos da categoria Sobremesa, recebido %d", len(resultado))
	}
}

func TestProdutoListarPorCategoria_Run_TodasCategorias(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
		{ID: 2, Nome: "X-Bacon", Categoria: entities.Lanche, Descricao: "Lanche com bacon", Preco: 30.0},
		{ID: 3, Nome: "Batata Frita", Categoria: entities.Acompanhamento, Descricao: "Acompanhamento crocante", Preco: 10.0},
		{ID: 4, Nome: "Onion Rings", Categoria: entities.Acompanhamento, Descricao: "Anéis de cebola", Preco: 12.0},
		{ID: 5, Nome: "Coca-Cola", Categoria: entities.Bebida, Descricao: "Bebida refrescante", Preco: 7.5},
		{ID: 6, Nome: "Suco", Categoria: entities.Bebida, Descricao: "Suco natural", Preco: 6.0},
		{ID: 7, Nome: "Sorvete", Categoria: entities.Sobremesa, Descricao: "Sobremesa gelada", Preco: 8.0},
		{ID: 8, Nome: "Pudim", Categoria: entities.Sobremesa, Descricao: "Sobremesa doce", Preco: 9.0},
	}

	mockRepo := &MockProdutoRepositoryListarPorCategoria{
		Produtos: produtos,
	}

	useCase := NewProdutoListarPorCategoriaUseCase(mockRepo)

	ctx := context.Background()

	testCases := []struct {
		categoria          string
		quantidadeEsperada int
	}{
		{"Lanche", 2},
		{"Acompanhamento", 2},
		{"Bebida", 2},
		{"Sobremesa", 2},
	}

	for _, tc := range testCases {
		// When
		resultado, err := useCase.Run(ctx, tc.categoria)

		// Then
		if err != nil {
			t.Errorf("Não esperado erro para categoria %s, recebido %v", tc.categoria, err)
		}

		if len(resultado) != tc.quantidadeEsperada {
			t.Errorf("Esperado %d produtos da categoria %s, recebido %d", tc.quantidadeEsperada, tc.categoria, len(resultado))
		}

		// Verificar se todos os produtos são da categoria correta
		for _, produto := range resultado {
			if string(produto.Categoria) != tc.categoria {
				t.Errorf("Esperado categoria %s, recebido %s", tc.categoria, produto.Categoria)
			}
		}
	}
}

func TestProdutoListarPorCategoria_Run_CategoriaInvalida(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
	}

	mockRepo := &MockProdutoRepositoryListarPorCategoria{
		Produtos: produtos,
	}

	useCase := NewProdutoListarPorCategoriaUseCase(mockRepo)

	ctx := context.Background()

	// When - buscando categoria inválida
	resultado, err := useCase.Run(ctx, "CategoriaInexistente")

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		resultado = []*entities.Produto{} // Tratar nil como lista vazia
	}

	if len(resultado) != 0 {
		t.Errorf("Esperado 0 produtos para categoria inválida, recebido %d", len(resultado))
	}
}

func TestProdutoListarPorCategoria_Run_ListaVaziaGeral(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryListarPorCategoria{
		Produtos: []*entities.Produto{},
	}

	useCase := NewProdutoListarPorCategoriaUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx, "Lanche")

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		resultado = []*entities.Produto{} // Tratar nil como lista vazia
	}

	if len(resultado) != 0 {
		t.Errorf("Esperado 0 produtos quando repositório vazio, recebido %d", len(resultado))
	}
}

func TestProdutoListarPorCategoria_Run_CaseSensitive(t *testing.T) {
	// Given
	produtos := []*entities.Produto{
		{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Lanche delicioso", Preco: 25.0},
	}

	mockRepo := &MockProdutoRepositoryListarPorCategoria{
		Produtos: produtos,
	}

	useCase := NewProdutoListarPorCategoriaUseCase(mockRepo)

	ctx := context.Background()

	// Test cases com diferentes variações de case
	testCases := []struct {
		categoria          string
		quantidadeEsperada int
	}{
		{"Lanche", 1}, // correto
		{"lanche", 0}, // minúscula
		{"LANCHE", 0}, // maiúscula
		{"LaNcHe", 0}, // misto
	}

	for _, tc := range testCases {
		// When
		resultado, err := useCase.Run(ctx, tc.categoria)

		// Then
		if err != nil {
			t.Errorf("Não esperado erro para categoria %s, recebido %v", tc.categoria, err)
		}

		if len(resultado) != tc.quantidadeEsperada {
			t.Errorf("Para categoria '%s' esperado %d produtos, recebido %d", tc.categoria, tc.quantidadeEsperada, len(resultado))
		}
	}
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
