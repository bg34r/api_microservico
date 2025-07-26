package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"testing"
	"time"
)

// MockPedidoRepositoryBuscar implements repository.PedidoRepository for testing
type MockPedidoRepositoryBuscar struct {
	Pedidos []*entities.Pedido
}

func (m *MockPedidoRepositoryBuscar) CriarPedido(ctx context.Context, pedido *entities.Pedido) error {
	return nil
}

func (m *MockPedidoRepositoryBuscar) BuscarPedido(ctx context.Context, id int) (*entities.Pedido, error) {
	for _, p := range m.Pedidos {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("pedido não encontrado")
}

func (m *MockPedidoRepositoryBuscar) AtualizarStatusPedido(ctx context.Context, pedidoID int, status string, ultimaAtualizacao time.Time) error {
	return nil
}

func (m *MockPedidoRepositoryBuscar) ListarTodosOsPedidos(ctx context.Context) ([]*entities.Pedido, error) {
	return nil, nil
}

func (m *MockPedidoRepositoryBuscar) AtualizarStatusPagamento(ctx context.Context, pedidoID int, statusPagamento string, ultimaAtualizacao time.Time) error {
	return nil
}

func TestPedidoBuscarPorIdUseCase_Run_Success(t *testing.T) {
	mockRepo := &MockPedidoRepositoryBuscar{}
	useCase := NewPedidoBuscarPorIdUseCase(mockRepo)

	// Setup pedido no repositório
	pedido := &entities.Pedido{
		ID:          1,
		ClienteNome: "João Silva",
		Status:      entities.Pendente,
		Produtos: []entities.Produto{
			{ID: 1, Nome: "Hamburguer", Categoria: entities.Lanche, Preco: 25.0},
		},
	}
	mockRepo.Pedidos = []*entities.Pedido{pedido}

	// Test
	result, err := useCase.Run(context.Background(), 1)

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected pedido, got nil")
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.ClienteNome != "João Silva" {
		t.Errorf("expected ClienteNome 'João Silva', got %s", result.ClienteNome)
	}
}

func TestPedidoBuscarPorIdUseCase_Run_NotFound(t *testing.T) {
	mockRepo := &MockPedidoRepositoryBuscar{}
	useCase := NewPedidoBuscarPorIdUseCase(mockRepo)

	// Test
	result, err := useCase.Run(context.Background(), 999)

	// Assertions
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err.Error() != "pedido não encontrado" {
		t.Errorf("expected 'pedido não encontrado', got %s", err.Error())
	}
}

func TestPedidoBuscarPorIdUseCase_Run_InvalidID(t *testing.T) {
	mockRepo := &MockPedidoRepositoryBuscar{}
	useCase := NewPedidoBuscarPorIdUseCase(mockRepo)

	// Test with invalid ID (0 or negative)
	testCases := []int{0, -1, -999}

	for _, id := range testCases {
		result, err := useCase.Run(context.Background(), id)

		if err == nil {
			t.Errorf("expected error for invalid ID %d, got nil", id)
		}
		if result != nil {
			t.Errorf("expected nil result for invalid ID %d, got %+v", id, result)
		}
	}
}
