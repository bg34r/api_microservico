package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lanchonete/internal/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock UseCases ---
type MockProdutoIncluirUseCase struct{ mock.Mock }

func (m *MockProdutoIncluirUseCase) Run(c context.Context, identificacao string, nome, categoria string, preco float32) (*entities.Produto, error) {
	args := m.Called(c, identificacao, nome, categoria, preco)
	return args.Get(0).(*entities.Produto), args.Error(1)
}

type MockProdutoBuscaPorIdUseCase struct{ mock.Mock }

func (m *MockProdutoBuscaPorIdUseCase) Run(c context.Context, id int) (*entities.Produto, error) {
	args := m.Called(c, id)
	return args.Get(0).(*entities.Produto), args.Error(1)
}

type MockProdutoListarTodosUseCase struct{ mock.Mock }

func (m *MockProdutoListarTodosUseCase) Run(c context.Context) ([]*entities.Produto, error) {
	args := m.Called(c)
	return args.Get(0).([]*entities.Produto), args.Error(1)
}

type MockProdutoEditarUseCase struct{ mock.Mock }

func (m *MockProdutoEditarUseCase) Run(c context.Context, id int, nome, categoria, descricao string, preco float32) (*entities.Produto, error) {
	args := m.Called(c, id, nome, categoria, descricao, preco)
	return args.Get(0).(*entities.Produto), args.Error(1)
}

type MockProdutoRemoverUseCase struct{ mock.Mock }

func (m *MockProdutoRemoverUseCase) Run(c context.Context, id int) error {
	args := m.Called(c, id)
	return args.Error(0)
}

type MockProdutoListarPorCategoriaUseCase struct{ mock.Mock }

func (m *MockProdutoListarPorCategoriaUseCase) Run(c context.Context, categoria string) ([]*entities.Produto, error) {
	args := m.Called(c, categoria)
	return args.Get(0).([]*entities.Produto), args.Error(1)
}

// --- Teste do Construtor (IMPORTANTE) ---
func TestNewProdutoHandler(t *testing.T) {
	// Mocks dos use cases
	mockIncluir := new(MockProdutoIncluirUseCase)
	mockBuscar := new(MockProdutoBuscaPorIdUseCase)
	mockListar := new(MockProdutoListarTodosUseCase)
	mockEditar := new(MockProdutoEditarUseCase)
	mockRemover := new(MockProdutoRemoverUseCase)
	mockListarCategoria := new(MockProdutoListarPorCategoriaUseCase)

	// Testar construtor
	handler := NewProdutoHandler(
		mockIncluir,
		mockBuscar,
		mockListar,
		mockEditar,
		mockRemover,
		mockListarCategoria,
	)

	// Verificações
	assert.NotNil(t, handler)
	assert.Equal(t, mockIncluir, handler.ProdutoIncluirUseCase)
	assert.Equal(t, mockBuscar, handler.ProdutoBuscarPorIdUseCase)
	assert.Equal(t, mockListar, handler.ProdutoListarTodosUseCase)
	assert.Equal(t, mockEditar, handler.ProdutoEditarUseCase)
	assert.Equal(t, mockRemover, handler.ProdutoRemoverUseCase)
	assert.Equal(t, mockListarCategoria, handler.ProdutoListarPorCategoriaUseCase)
}

// --- Testes dos Métodos ---
func TestProdutoHandler_ProdutoIncluir(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(MockProdutoIncluirUseCase)
	handler := &ProdutoHandler{
		ProdutoIncluirUseCase: mockUC,
	}

	prod := entities.Produto{
		Nome:      "Coca-Cola",
		Categoria: "Bebida",
		Descricao: "Refrigerante",
		Preco:     5.0,
	}
	mockUC.On("Run", mock.Anything, prod.Nome, string(prod.Categoria), prod.Descricao, prod.Preco).
		Return(&prod, nil)

	body, _ := json.Marshal(prod)
	req, _ := http.NewRequest(http.MethodPost, "/produto", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ProdutoIncluir(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Produto incluído com sucesso")
}

func TestProdutoHandler_ProdutoBuscarPorId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(MockProdutoBuscaPorIdUseCase)
	handler := &ProdutoHandler{
		ProdutoBuscarPorIdUseCase: mockUC,
	}

	prod := &entities.Produto{Nome: "Coca-Cola"}
	mockUC.On("Run", mock.Anything, 1).Return(prod, nil)

	req, _ := http.NewRequest(http.MethodGet, "/produto/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = req

	handler.ProdutoBuscarPorId(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Coca-Cola")
}

func TestProdutoHandler_ProdutoListarTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(MockProdutoListarTodosUseCase)
	handler := &ProdutoHandler{
		ProdutoListarTodosUseCase: mockUC,
	}

	prods := []*entities.Produto{{Nome: "Coca-Cola"}}
	mockUC.On("Run", mock.Anything).Return(prods, nil)

	req, _ := http.NewRequest(http.MethodGet, "/produtos", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ProdutoListarTodos(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Coca-Cola")
}

func TestProdutoHandler_ProdutoEditar(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(MockProdutoEditarUseCase)
	handler := &ProdutoHandler{
		ProdutoEditarUseCase: mockUC,
	}

	prod := entities.Produto{
		ID:        1,
		Nome:      "Coca-Cola",
		Categoria: "Bebida",
		Descricao: "Refrigerante",
		Preco:     5.0,
	}
	mockUC.On("Run", mock.Anything, 1, prod.Nome, string(prod.Categoria), prod.Descricao, prod.Preco).
		Return(&prod, nil)

	body, _ := json.Marshal(prod)
	req, _ := http.NewRequest(http.MethodPost, "/produto/editar", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ProdutoEditar(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Produto editado com sucesso")
}

func TestProdutoHandler_ProdutoRemover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(MockProdutoRemoverUseCase)
	handler := &ProdutoHandler{
		ProdutoRemoverUseCase: mockUC,
	}

	mockUC.On("Run", mock.Anything, 1).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/produto/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = req

	handler.ProdutoRemover(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Produto removido com sucesso")
}

func TestProdutoHandler_ProdutoListarPorCategoria(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(MockProdutoListarPorCategoriaUseCase)
	handler := &ProdutoHandler{
		ProdutoListarPorCategoriaUseCase: mockUC,
	}

	prods := []*entities.Produto{{Nome: "Coca-Cola", Categoria: "Bebida"}}
	mockUC.On("Run", mock.Anything, "Bebida").Return(prods, nil)

	req, _ := http.NewRequest(http.MethodGet, "/produtos/Bebida", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "categoria", Value: "Bebida"}}
	c.Request = req

	handler.ProdutoListarPorCategoria(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Coca-Cola")
}
