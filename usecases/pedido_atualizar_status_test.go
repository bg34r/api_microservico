package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"testing"
	"time"
)

// MockPedidoRepositoryAtualizarStatus implements repository.PedidoRepository for testing
type MockPedidoRepositoryAtualizarStatus struct {
	Pedidos []*entities.Pedido
}

type MockEventPublisherAtualizar struct{}

func (m *MockEventPublisherAtualizar) Publish(eventType string, payload interface{}) error {
	// apenas retorna nil, simula sucesso
	return nil
}

func (m *MockPedidoRepositoryAtualizarStatus) CriarPedido(ctx context.Context, pedido *entities.Pedido) error {
	return nil
}

func (m *MockPedidoRepositoryAtualizarStatus) BuscarPedido(ctx context.Context, id int) (*entities.Pedido, error) {
	for _, p := range m.Pedidos {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("pedido não encontrado")
}

func (m *MockPedidoRepositoryAtualizarStatus) AtualizarStatusPedido(ctx context.Context, pedidoID int, status string, ultimaAtualizacao time.Time) error {
	for _, p := range m.Pedidos {
		if p.ID == pedidoID {
			p.Status = entities.StatusPedido(status)
			p.UltimaAtualizacao = ultimaAtualizacao
			return nil
		}
	}
	return errors.New("pedido não encontrado")
}

func (m *MockPedidoRepositoryAtualizarStatus) ListarTodosOsPedidos(ctx context.Context) ([]*entities.Pedido, error) {
	return nil, nil
}

func (m *MockPedidoRepositoryAtualizarStatus) AtualizarStatusPagamento(ctx context.Context, pedidoID int, statusPagamento string, ultimaAtualizacao time.Time) error {
	return nil
}

func TestPedidoAtualizarStatusUseCase_Run_Success(t *testing.T) {
	mockRepo := &MockPedidoRepositoryAtualizarStatus{}
	mockPublisher := &MockEventPublisherAtualizar{}
	useCase := NewPedidoAtualizarStatusUseCase(mockRepo, mockPublisher)

	// Setup pedido no repositório
	pedido := &entities.Pedido{
		ID:          1,
		ClienteNome: "João Silva",
		Status:      entities.Pendente,
	}
	mockRepo.Pedidos = []*entities.Pedido{pedido}

	// Test
	err := useCase.Run(context.Background(), 1, "Recebido")

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mockRepo.Pedidos[0].Status != entities.Recebido {
		t.Errorf("expected status 'Recebido', got %s", mockRepo.Pedidos[0].Status)
	}
}

func TestPedidoAtualizarStatusUseCase_Run_AllValidStatuses(t *testing.T) {
	validStatuses := []string{"Pendente", "Recebido", "Em preparação", "Pronto", "Finalizado"}

	for _, status := range validStatuses {
		mockRepo := &MockPedidoRepositoryAtualizarStatus{}
		mockPublisher := &MockEventPublisherAtualizar{}
		useCase := NewPedidoAtualizarStatusUseCase(mockRepo, mockPublisher)

		// Setup pedido no repositório
		pedido := &entities.Pedido{
			ID:          1,
			ClienteNome: "João Silva",
			Status:      entities.Pendente,
		}
		mockRepo.Pedidos = []*entities.Pedido{pedido}

		// Test
		err := useCase.Run(context.Background(), 1, status)

		// Assertions
		if err != nil {
			t.Errorf("expected no error for status '%s', got %v", status, err)
		}
		if string(mockRepo.Pedidos[0].Status) != status {
			t.Errorf("expected status '%s', got %s", status, mockRepo.Pedidos[0].Status)
		}
	}
}

func TestPedidoAtualizarStatusUseCase_Run_InvalidStatus(t *testing.T) {
	mockRepo := &MockPedidoRepositoryAtualizarStatus{}
	mockPublisher := &MockEventPublisherAtualizar{}
	useCase := NewPedidoAtualizarStatusUseCase(mockRepo, mockPublisher)

	// Setup pedido no repositório
	pedido := &entities.Pedido{
		ID:          1,
		ClienteNome: "João Silva",
		Status:      entities.Pendente,
	}
	mockRepo.Pedidos = []*entities.Pedido{pedido}

	// Test com status inválido
	err := useCase.Run(context.Background(), 1, "StatusInvalido")

	// Assertions
	if err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}
}

func TestPedidoAtualizarStatusUseCase_Run_PedidoNotFound(t *testing.T) {
	mockRepo := &MockPedidoRepositoryAtualizarStatus{}
	mockPublisher := &MockEventPublisherAtualizar{}
	useCase := NewPedidoAtualizarStatusUseCase(mockRepo, mockPublisher)

	// Test sem pedidos no repositório
	err := useCase.Run(context.Background(), 999, "Recebido")

	// Assertions
	if err == nil {
		t.Fatal("expected error for non-existent pedido, got nil")
	}
}

func TestPedidoAtualizarStatusUseCase_Run_StatusProgression(t *testing.T) {
	mockRepo := &MockPedidoRepositoryAtualizarStatus{}
	mockPublisher := &MockEventPublisherAtualizar{}
	useCase := NewPedidoAtualizarStatusUseCase(mockRepo, mockPublisher)

	// Setup pedido no repositório
	pedido := &entities.Pedido{
		ID:          1,
		ClienteNome: "João Silva",
		Status:      entities.Pendente,
	}
	mockRepo.Pedidos = []*entities.Pedido{pedido}

	// Test progression: Pendente -> Recebido -> Em preparação -> Pronto -> Finalizado
	statusProgression := []string{"Recebido", "Em preparação", "Pronto", "Finalizado"}

	for _, status := range statusProgression {
		err := useCase.Run(context.Background(), 1, status)
		if err != nil {
			t.Fatalf("expected no error for status '%s', got %v", status, err)
		}
		if string(mockRepo.Pedidos[0].Status) != status {
			t.Errorf("expected status '%s', got %s", status, mockRepo.Pedidos[0].Status)
		}
	}
}
