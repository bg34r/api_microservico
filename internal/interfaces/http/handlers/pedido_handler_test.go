package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"lanchonete/internal/domain/entities"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock UseCases ---

type MockPedidoIncluirUseCase struct{ mock.Mock }

func (m *MockPedidoIncluirUseCase) Run(ctx context.Context, clienteNome string, produtos []entities.Produto, personalizacao *string) (*entities.Pedido, error) {
	args := m.Called(ctx, clienteNome, produtos, personalizacao)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Pedido), args.Error(1)
}

type MockPedidoBuscarPorIdUseCase struct{ mock.Mock }

func (m *MockPedidoBuscarPorIdUseCase) Run(ctx context.Context, identificacao int) (*entities.Pedido, error) {
	args := m.Called(ctx, identificacao)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Pedido), args.Error(1)
}

type MockPedidoAtualizarStatusUseCase struct{ mock.Mock }

func (m *MockPedidoAtualizarStatusUseCase) Run(ctx context.Context, identificacao int, status string) error {
	args := m.Called(ctx, identificacao, status)
	return args.Error(0)
}

type MockPedidoAtualizarStatusPagamentoUseCase struct{ mock.Mock }

func (m *MockPedidoAtualizarStatusPagamentoUseCase) Run(ctx context.Context, identificacao int, statusPagamento string) error {
	args := m.Called(ctx, identificacao, statusPagamento)
	return args.Error(0)
}

type MockPedidoListarTodosUseCase struct{ mock.Mock }

func (m *MockPedidoListarTodosUseCase) Run(ctx context.Context) ([]*entities.Pedido, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Pedido), args.Error(1)
}

type MockProdutoBuscarUseCase struct{ mock.Mock }

func (m *MockProdutoBuscarUseCase) Run(c context.Context, id int) (*entities.Produto, error) {
	args := m.Called(c, id)
	return args.Get(0).(*entities.Produto), args.Error(1)
}

// --- Setup Handler ---

func setupPedidoHandlerWithMocks() (*PedidoHandler,
	*MockPedidoIncluirUseCase,
	*MockPedidoBuscarPorIdUseCase,
	*MockPedidoAtualizarStatusUseCase,
	*MockPedidoAtualizarStatusPagamentoUseCase,
	*MockProdutoBuscarUseCase,
	*MockPedidoListarTodosUseCase,
) {
	mockIncluir := new(MockPedidoIncluirUseCase)
	mockBuscar := new(MockPedidoBuscarPorIdUseCase)
	mockAtualizar := new(MockPedidoAtualizarStatusUseCase)
	mockAtualizarPagamento := new(MockPedidoAtualizarStatusPagamentoUseCase)
	mockProdutoBuscar := new(MockProdutoBuscarUseCase)
	mockListar := new(MockPedidoListarTodosUseCase)

	handler := &PedidoHandler{
		PedidoIncluirUseCase:                  mockIncluir,
		PedidoBuscarPorIdUseCase:              mockBuscar,
		PedidoAtualizarStatusUseCase:          mockAtualizar,
		PedidoAtualizarStatusPagamentoUseCase: mockAtualizarPagamento,
		ProdutoBuscarPorIdUseCase:             mockProdutoBuscar,
		PedidoListarTodosUseCase:              mockListar,
	}
	return handler, mockIncluir, mockBuscar, mockAtualizar, mockAtualizarPagamento, mockProdutoBuscar, mockListar
}

// --- Tests ---

func TestPedidoHandler_CriarPedido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, mockIncluir, _, _, _, mockProdutoBuscar, _ := setupPedidoHandlerWithMocks()

	clienteNome := "João Silva"
	produto := entities.Produto{ID: 1, Nome: "Produto Teste", Categoria: entities.Lanche, Preco: 10}
	pedido := entities.Pedido{ID: 1, ClienteNome: clienteNome, Produtos: []entities.Produto{produto}}

	mockProdutoBuscar.On("Run", mock.Anything, produto.ID).Return(&produto, nil)
	mockIncluir.On("Run", mock.Anything, clienteNome, []entities.Produto{produto}, (*string)(nil)).
		Return(&entities.Pedido{ID: 1}, nil)

	body, _ := json.Marshal(pedido)
	req, _ := http.NewRequest(http.MethodPost, "/pedidos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CriarPedido(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pedido criado com sucesso")
}

func TestPedidoHandler_BuscarPedido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, mockBuscar, _, _, _, _ := setupPedidoHandlerWithMocks()

	pedido := &entities.Pedido{ID: 1}
	mockBuscar.On("Run", mock.Anything, 1).Return(pedido, nil)

	req, _ := http.NewRequest(http.MethodGet, "/pedidos/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "nroPedido", Value: "1"}}
	c.Request = req

	handler.BuscarPedido(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPedidoHandler_AtualizarStatusPedido(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _, mockAtualizar, _, _, _ := setupPedidoHandlerWithMocks()

	mockAtualizar.On("Run", mock.Anything, 1, "Finalizado").Return(nil)

	req, _ := http.NewRequest(http.MethodPut, "/pedidos/1/status/Finalizado", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "nroPedido", Value: "1"},
		{Key: "status", Value: "Finalizado"},
	}
	c.Request = req

	handler.AtualizarStatusPedido(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Status do pedido atualizado com sucesso")
}

func TestPedidoHandler_ListarTodosOsPedidos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _, _, _, _, mockListar := setupPedidoHandlerWithMocks()

	pedidos := []*entities.Pedido{
		{ID: 1},
		{ID: 2},
	}
	mockListar.On("Run", mock.Anything).Return(pedidos, nil)

	req, _ := http.NewRequest(http.MethodGet, "/pedidos/listartodos", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListarTodosOsPedidos(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPedidoHandler_AtualizarStatusPagamento(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _, _, mockAtualizarPagamento, _, _ := setupPedidoHandlerWithMocks()

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
	mockAtualizarPagamento.AssertExpectations(t)
}
