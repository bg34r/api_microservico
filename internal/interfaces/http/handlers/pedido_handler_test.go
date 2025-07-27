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
type MockPedidoIncluirUseCase struct{ mock.Mock }

func (m *MockPedidoIncluirUseCase) Run(ctx context.Context, clienteNome string, produtos []entities.Produto, personalizacao *string) (*entities.Pedido, error) {
	args := m.Called(ctx, clienteNome, produtos, personalizacao)
	return args.Get(0).(*entities.Pedido), args.Error(1)
}

type MockPedidoBuscarPorIdUseCase struct{ mock.Mock }

func (m *MockPedidoBuscarPorIdUseCase) Run(ctx context.Context, pedidoID int) (*entities.Pedido, error) {
	args := m.Called(ctx, pedidoID)
	return args.Get(0).(*entities.Pedido), args.Error(1)
}

type MockPedidoAtualizarStatusUseCase struct{ mock.Mock }

func (m *MockPedidoAtualizarStatusUseCase) Run(ctx context.Context, pedidoID int, novoStatus string) error {
	args := m.Called(ctx, pedidoID, novoStatus)
	return args.Error(0)
}

type MockPedidoAtualizarStatusPagamentoUseCase struct{ mock.Mock }

func (m *MockPedidoAtualizarStatusPagamentoUseCase) Run(ctx context.Context, pedidoID int, statusPagamento string) error {
	args := m.Called(ctx, pedidoID, statusPagamento)
	return args.Error(0)
}

type MockProdutoBuscarPorIdUseCase struct{ mock.Mock }

func (m *MockProdutoBuscarPorIdUseCase) Run(ctx context.Context, id int) (*entities.Produto, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Produto), args.Error(1)
}

type MockPedidoListarTodosUseCase struct{ mock.Mock }

func (m *MockPedidoListarTodosUseCase) Run(ctx context.Context) ([]*entities.Pedido, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Pedido), args.Error(1)
}

// --- Teste do Construtor ---
func TestNewPedidoHandler(t *testing.T) {
	// Mocks dos use cases
	mockIncluir := new(MockPedidoIncluirUseCase)
	mockBuscar := new(MockPedidoBuscarPorIdUseCase)
	mockAtualizarStatus := new(MockPedidoAtualizarStatusUseCase)
	mockAtualizarStatusPagamento := new(MockPedidoAtualizarStatusPagamentoUseCase)
	mockProdutoBuscar := new(MockProdutoBuscarPorIdUseCase)
	mockListarTodos := new(MockPedidoListarTodosUseCase)

	// Testar construtor
	handler := NewPedidoHandler(
		mockIncluir,
		mockBuscar,
		mockAtualizarStatus,
		mockAtualizarStatusPagamento,
		mockProdutoBuscar,
		mockListarTodos,
	)

	// Verificações
	assert.NotNil(t, handler)
	assert.Equal(t, mockIncluir, handler.PedidoIncluirUseCase)
	assert.Equal(t, mockBuscar, handler.PedidoBuscarPorIdUseCase)
	assert.Equal(t, mockAtualizarStatus, handler.PedidoAtualizarStatusUseCase)
	assert.Equal(t, mockAtualizarStatusPagamento, handler.PedidoAtualizarStatusPagamentoUseCase)
	assert.Equal(t, mockProdutoBuscar, handler.ProdutoBuscarPorIdUseCase)
	assert.Equal(t, mockListarTodos, handler.PedidoListarTodosUseCase)
}

// --- Testes dos Métodos ---
func TestPedidoHandler_CriarPedido(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPedidoIncluir := new(MockPedidoIncluirUseCase)
	mockProdutoBuscar := new(MockProdutoBuscarPorIdUseCase)

	handler := &PedidoHandler{
		PedidoIncluirUseCase:      mockPedidoIncluir,
		ProdutoBuscarPorIdUseCase: mockProdutoBuscar,
	}

	// Produto mockado que será retornado pela busca
	produtoCompleto := &entities.Produto{
		ID:        1,
		Nome:      "Hamburger",
		Categoria: "Lanche",
		Descricao: "Hamburger clássico",
		Preco:     15.0,
	}

	// Pedido de entrada com apenas o ID do produto
	personalizacao := "Sem cebola"
	pedidoRequest := entities.Pedido{
		ClienteNome: "João Silva",
		Produtos: []entities.Produto{
			{ID: 1}, // Apenas ID será enviado
		},
		Personalizacao: &personalizacao,
	}

	// Pedido que será retornado pelo use case
	pedidoRetorno := &entities.Pedido{
		ID:             123,
		ClienteNome:    "João Silva",
		Produtos:       []entities.Produto{*produtoCompleto},
		Status:         entities.Pendente,
		Personalizacao: &personalizacao,
	}

	// Mocks
	mockProdutoBuscar.On("Run", mock.Anything, 1).Return(produtoCompleto, nil)
	mockPedidoIncluir.On("Run", mock.Anything, "João Silva", []entities.Produto{*produtoCompleto}, &personalizacao).
		Return(pedidoRetorno, nil)

	// Preparar request
	body, _ := json.Marshal(pedidoRequest)
	req, _ := http.NewRequest(http.MethodPost, "/pedidos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Executar
	handler.CriarPedido(c)

	// Verificações
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pedido criado com sucesso")
	mockProdutoBuscar.AssertExpectations(t)
	mockPedidoIncluir.AssertExpectations(t)
}

func TestPedidoHandler_BuscarPedido(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBuscar := new(MockPedidoBuscarPorIdUseCase)
	handler := &PedidoHandler{
		PedidoBuscarPorIdUseCase: mockBuscar,
	}

	pedido := &entities.Pedido{
		ID:          1,
		ClienteNome: "João Silva",
		Status:      entities.Pendente,
	}

	mockBuscar.On("Run", mock.Anything, 1).Return(pedido, nil)

	req, _ := http.NewRequest(http.MethodGet, "/pedidos/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "nroPedido", Value: "1"}}
	c.Request = req

	handler.BuscarPedido(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "João Silva")
	mockBuscar.AssertExpectations(t)
}

func TestPedidoHandler_AtualizarStatusPedido(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAtualizar := new(MockPedidoAtualizarStatusUseCase)
	handler := &PedidoHandler{
		PedidoAtualizarStatusUseCase: mockAtualizar,
	}

	mockAtualizar.On("Run", mock.Anything, 1, "EmPreparacao").Return(nil)

	req, _ := http.NewRequest(http.MethodPut, "/pedidos/1/status/EmPreparacao", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "nroPedido", Value: "1"},
		{Key: "status", Value: "EmPreparacao"},
	}
	c.Request = req

	handler.AtualizarStatusPedido(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Status do pedido atualizado com sucesso")
	mockAtualizar.AssertExpectations(t)
}

func TestPedidoHandler_AtualizarStatusPagamento(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAtualizarPagamento := new(MockPedidoAtualizarStatusPagamentoUseCase)
	handler := &PedidoHandler{
		PedidoAtualizarStatusPagamentoUseCase: mockAtualizarPagamento,
	}

	mockAtualizarPagamento.On("Run", mock.Anything, 1, "Pago").Return(nil)

	req, _ := http.NewRequest(http.MethodPut, "/pedidos/1/pagamento/Pago", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "nroPedido", Value: "1"},
		{Key: "statusPagamento", Value: "Pago"},
	}
	c.Request = req

	handler.AtualizarStatusPagamento(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Status de pagamento atualizado com sucesso")
	mockAtualizarPagamento.AssertExpectations(t)
}

func TestPedidoHandler_ListarTodosOsPedidos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockListar := new(MockPedidoListarTodosUseCase)
	handler := &PedidoHandler{
		PedidoListarTodosUseCase: mockListar,
	}

	pedidos := []*entities.Pedido{
		{
			ID:          1,
			ClienteNome: "João Silva",
			Status:      entities.Pendente,
		},
		{
			ID:          2,
			ClienteNome: "Maria Santos",
			Status:      entities.Pendente,
		},
	}

	mockListar.On("Run", mock.Anything).Return(pedidos, nil)

	req, _ := http.NewRequest(http.MethodGet, "/pedidos/listartodos", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListarTodosOsPedidos(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "João Silva")
	assert.Contains(t, w.Body.String(), "Maria Santos")
	mockListar.AssertExpectations(t)
}

// --- Testes de Casos de Erro ---
func TestPedidoHandler_BuscarPedido_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &PedidoHandler{}

	req, _ := http.NewRequest(http.MethodGet, "/pedidos/abc", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "nroPedido", Value: "abc"}}
	c.Request = req

	handler.BuscarPedido(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Número do pedido inválido")
}

func TestPedidoHandler_AtualizarStatusPedido_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &PedidoHandler{}

	req, _ := http.NewRequest(http.MethodPut, "/pedidos/abc/status/EmPreparacao", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "nroPedido", Value: "abc"},
		{Key: "status", Value: "EmPreparacao"},
	}
	c.Request = req

	handler.AtualizarStatusPedido(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Número do pedido inválido")
}

func TestPedidoHandler_AtualizarStatusPagamento_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &PedidoHandler{}

	req, _ := http.NewRequest(http.MethodPut, "/pedidos/abc/pagamento/Pago", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "nroPedido", Value: "abc"},
		{Key: "statusPagamento", Value: "Pago"},
	}
	c.Request = req

	handler.AtualizarStatusPagamento(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Número do pedido inválido")
}
