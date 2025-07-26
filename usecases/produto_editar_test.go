package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"strings"
	"testing"
)

// MockProdutoRepositoryEditar implements repository.ProdutoRepository for testing
type MockProdutoRepositoryEditar struct {
	Produtos []*entities.Produto
}

func (m *MockProdutoRepositoryEditar) AdicionarProduto(ctx context.Context, produto *entities.Produto) error {
	return nil
}

func (m *MockProdutoRepositoryEditar) EditarProduto(ctx context.Context, produto *entities.Produto) error {
	// Simula sucesso sempre - repositório real faria a edição
	return nil
}

func (m *MockProdutoRepositoryEditar) RemoverProduto(ctx context.Context, id int) error {
	return nil
}

func (m *MockProdutoRepositoryEditar) BuscarProdutoPorId(ctx context.Context, id int) (*entities.Produto, error) {
	for _, produto := range m.Produtos {
		if produto.ID == id {
			return produto, nil
		}
	}
	return nil, errors.New("produto não encontrado")
}

func (m *MockProdutoRepositoryEditar) ListarTodosOsProdutos(ctx context.Context) ([]*entities.Produto, error) {
	return m.Produtos, nil
}

func (m *MockProdutoRepositoryEditar) ListarPorCategoria(ctx context.Context, categoria string) ([]*entities.Produto, error) {
	var produtosFiltrados []*entities.Produto
	for _, produto := range m.Produtos {
		if string(produto.Categoria) == categoria {
			produtosFiltrados = append(produtosFiltrados, produto)
		}
	}
	return produtosFiltrados, nil
}

func TestProdutoEditar_Run_Sucesso(t *testing.T) {
	// Given
	produtoOriginal := &entities.Produto{
		ID:        1,
		Nome:      "Produto Original",
		Categoria: entities.Lanche,
		Descricao: "Descrição original",
		Preco:     15.0,
	}

	mockRepo := &MockProdutoRepositoryEditar{
		Produtos: []*entities.Produto{produtoOriginal},
	}

	useCase := NewProdutoEditarUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx, 1, "Produto Editado", "Bebida", "Nova descrição", 20.0)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		t.Fatal("Esperado produto editado, recebido nil")
	}

	if resultado.Nome != "Produto Editado" {
		t.Errorf("Esperado nome 'Produto Editado', recebido %s", resultado.Nome)
	}

	if resultado.Categoria != entities.Bebida {
		t.Errorf("Esperado categoria 'Bebida', recebido %s", resultado.Categoria)
	}

	if resultado.Preco != 20.0 {
		t.Errorf("Esperado preço 20.0, recebido %f", resultado.Preco)
	}
}

func TestProdutoEditar_Run_ProdutoNaoEncontrado(t *testing.T) {
	// Given
	mockRepo := &MockProdutoRepositoryEditar{
		Produtos: []*entities.Produto{},
	}

	useCase := NewProdutoEditarUseCase(mockRepo)

	ctx := context.Background()

	// When
	resultado, err := useCase.Run(ctx, 999, "Produto Teste", "Lanche", "Descrição", 10.0)

	// Then
	if err == nil {
		t.Error("Esperado erro, recebido nil")
	}

	if resultado != nil {
		t.Error("Esperado nil, recebido produto")
	}

	if !strings.Contains(err.Error(), "produto não cadastrado") {
		t.Errorf("Esperado erro sobre produto não encontrado, recebido: %v", err)
	}
}

func TestProdutoEditar_Run_CamposVazios(t *testing.T) {
	// Given - produto original manterá valores quando campos vazios
	produtoOriginal := &entities.Produto{
		ID:        1,
		Nome:      "Produto Original",
		Categoria: entities.Lanche,
		Descricao: "Descrição original",
		Preco:     15.0,
	}

	mockRepo := &MockProdutoRepositoryEditar{
		Produtos: []*entities.Produto{produtoOriginal},
	}

	useCase := NewProdutoEditarUseCase(mockRepo)

	ctx := context.Background()

	// When - passando campos vazios (devem manter valores originais)
	resultado, err := useCase.Run(ctx, 1, "", "", "", 0)

	// Then
	if err != nil {
		t.Errorf("Esperado nil, recebido %v", err)
	}

	if resultado == nil {
		t.Fatal("Esperado produto editado, recebido nil")
	}

	// Deve manter valores originais
	if resultado.Nome != "Produto Original" {
		t.Errorf("Esperado nome original 'Produto Original', recebido %s", resultado.Nome)
	}

	if resultado.Categoria != entities.Lanche {
		t.Errorf("Esperado categoria original 'Lanche', recebido %s", resultado.Categoria)
	}

	if resultado.Preco != 15.0 {
		t.Errorf("Esperado preço original 15.0, recebido %f", resultado.Preco)
	}
}

func TestProdutoEditar_Run_DadosInvalidos(t *testing.T) {
	// Given
	produtoOriginal := &entities.Produto{
		ID:        1,
		Nome:      "Produto Original",
		Categoria: entities.Lanche,
		Descricao: "Descrição original",
		Preco:     15.0,
	}

	mockRepo := &MockProdutoRepositoryEditar{
		Produtos: []*entities.Produto{produtoOriginal},
	}

	useCase := NewProdutoEditarUseCase(mockRepo)

	ctx := context.Background()

	// When - categoria inválida
	resultado, err := useCase.Run(ctx, 1, "Produto Teste", "CategoriaInvalida", "Descrição", 10.0)

	// Then
	if err == nil {
		t.Error("Esperado erro para categoria inválida, recebido nil")
	}

	if resultado != nil {
		t.Error("Esperado nil para dados inválidos, recebido produto")
	}

	if !strings.Contains(err.Error(), "atualização de produto inválida") {
		t.Errorf("Esperado erro sobre dados inválidos, recebido: %v", err)
	}
}
