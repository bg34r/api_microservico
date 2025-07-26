package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"testing"
	"time"
)

// MockPedidoRepositoryListar implements repository.PedidoRepository for testing
type MockPedidoRepositoryListar struct {
	Pedidos []*entities.Pedido
}

func (m *MockPedidoRepositoryListar) CriarPedido(ctx context.Context, pedido *entities.Pedido) error {
	return nil
}

func (m *MockPedidoRepositoryListar) BuscarPedido(ctx context.Context, id int) (*entities.Pedido, error) {
	return nil, nil
}

func (m *MockPedidoRepositoryListar) AtualizarStatusPedido(ctx context.Context, pedidoID int, status string, ultimaAtualizacao time.Time) error {
	return nil
}

func (m *MockPedidoRepositoryListar) ListarTodosOsPedidos(ctx context.Context) ([]*entities.Pedido, error) {
	return m.Pedidos, nil
}

func (m *MockPedidoRepositoryListar) AtualizarStatusPagamento(ctx context.Context, pedidoID int, statusPagamento string, ultimaAtualizacao time.Time) error {
	return nil
}

func TestPedidoListarTodosUseCase_Run_Success(t *testing.T) {
	mockRepo := &MockPedidoRepositoryListar{}
	useCase := NewPedidoListarTodosUseCase(mockRepo)

	// Setup pedidos no repositório
	pedidos := []*entities.Pedido{
		{ID: 1, ClienteNome: "João", Status: entities.Pendente},
		{ID: 2, ClienteNome: "Maria", Status: entities.Recebido},
		{ID: 3, ClienteNome: "Pedro", Status: entities.Pronto},
	}
	mockRepo.Pedidos = pedidos

	// Test
	result, err := useCase.Run(context.Background())

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 pedidos, got %d", len(result))
	}
	if result[0].ClienteNome != "João" {
		t.Errorf("expected first pedido ClienteNome 'João', got %s", result[0].ClienteNome)
	}
}

func TestPedidoListarTodosUseCase_Run_Empty(t *testing.T) {
	mockRepo := &MockPedidoRepositoryListar{}
	useCase := NewPedidoListarTodosUseCase(mockRepo)

	// Test (repositório vazio)
	result, err := useCase.Run(context.Background())

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 pedidos, got %d", len(result))
	}
}

func TestPedidoListarTodosUseCase_Run_OrderedByStatus(t *testing.T) {
	mockRepo := &MockPedidoRepositoryListar{}
	useCase := NewPedidoListarTodosUseCase(mockRepo)

	// Setup pedidos com diferentes status
	pedidos := []*entities.Pedido{
		{ID: 1, ClienteNome: "João", Status: entities.Finalizado},
		{ID: 2, ClienteNome: "Maria", Status: entities.Pendente},
		{ID: 3, ClienteNome: "Pedro", Status: entities.EmPreparacao},
		{ID: 4, ClienteNome: "Ana", Status: entities.Pronto},
	}
	mockRepo.Pedidos = pedidos

	// Test
	result, err := useCase.Run(context.Background())

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 4 {
		t.Errorf("expected 4 pedidos, got %d", len(result))
	}

	// Verify all pedidos are returned
	statusCount := make(map[entities.StatusPedido]int)
	for _, pedido := range result {
		statusCount[pedido.Status]++
	}

	if statusCount[entities.Pendente] != 1 {
		t.Errorf("expected 1 Pendente pedido, got %d", statusCount[entities.Pendente])
	}
	if statusCount[entities.EmPreparacao] != 1 {
		t.Errorf("expected 1 EmPreparacao pedido, got %d", statusCount[entities.EmPreparacao])
	}
}
