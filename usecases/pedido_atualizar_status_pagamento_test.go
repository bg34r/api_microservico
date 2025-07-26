package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"testing"
	"time"
)

// MockPedidoRepositoryAtualizarPagamento implements repository.PedidoRepository for testing
type MockPedidoRepositoryAtualizarPagamento struct {
	Pedidos []*entities.Pedido
}

func (m *MockPedidoRepositoryAtualizarPagamento) CriarPedido(ctx context.Context, pedido *entities.Pedido) error {
	return nil
}

func (m *MockPedidoRepositoryAtualizarPagamento) BuscarPedido(ctx context.Context, id int) (*entities.Pedido, error) {
	for _, p := range m.Pedidos {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("pedido não encontrado")
}

func (m *MockPedidoRepositoryAtualizarPagamento) AtualizarStatusPedido(ctx context.Context, pedidoID int, status string, ultimaAtualizacao time.Time) error {
	return nil
}

func (m *MockPedidoRepositoryAtualizarPagamento) ListarTodosOsPedidos(ctx context.Context) ([]*entities.Pedido, error) {
	return nil, nil
}

func (m *MockPedidoRepositoryAtualizarPagamento) AtualizarStatusPagamento(ctx context.Context, pedidoID int, statusPagamento string, ultimaAtualizacao time.Time) error {
	for _, p := range m.Pedidos {
		if p.ID == pedidoID {
			p.StatusPagamento = statusPagamento
			p.UltimaAtualizacao = ultimaAtualizacao
			return nil
		}
	}
	return errors.New("pedido não encontrado")
}

func TestPedidoAtualizarStatusPagamentoUseCase_Run_Success(t *testing.T) {
	mockRepo := &MockPedidoRepositoryAtualizarPagamento{}
	useCase := NewPedidoAtualizarStatusPagamentoUseCase(mockRepo)

	// Setup pedido no repositório
	pedido := &entities.Pedido{
		ID:              1,
		ClienteNome:     "João Silva",
		StatusPagamento: "Pendente",
	}
	mockRepo.Pedidos = []*entities.Pedido{pedido}

	// Test
	err := useCase.Run(context.Background(), 1, "Pago")

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mockRepo.Pedidos[0].StatusPagamento != "Pago" {
		t.Errorf("expected StatusPagamento 'Pago', got %s", mockRepo.Pedidos[0].StatusPagamento)
	}
}

func TestPedidoAtualizarStatusPagamentoUseCase_Run_AllValidStatuses(t *testing.T) {
	validStatuses := []string{"Pendente", "Pago", "Recusado", "Cancelado"}

	for _, status := range validStatuses {
		mockRepo := &MockPedidoRepositoryAtualizarPagamento{}
		useCase := NewPedidoAtualizarStatusPagamentoUseCase(mockRepo)

		// Setup pedido no repositório
		pedido := &entities.Pedido{
			ID:              1,
			ClienteNome:     "João Silva",
			StatusPagamento: "Pendente",
		}
		mockRepo.Pedidos = []*entities.Pedido{pedido}

		// Test
		err := useCase.Run(context.Background(), 1, status)

		// Assertions
		if err != nil {
			t.Errorf("expected no error for status '%s', got %v", status, err)
		}
		if mockRepo.Pedidos[0].StatusPagamento != status {
			t.Errorf("expected StatusPagamento '%s', got %s", status, mockRepo.Pedidos[0].StatusPagamento)
		}
	}
}

func TestPedidoAtualizarStatusPagamentoUseCase_Run_InvalidStatus(t *testing.T) {
	mockRepo := &MockPedidoRepositoryAtualizarPagamento{}
	useCase := NewPedidoAtualizarStatusPagamentoUseCase(mockRepo)

	// Setup pedido no repositório
	pedido := &entities.Pedido{
		ID:              1,
		ClienteNome:     "João Silva",
		StatusPagamento: "Pendente",
	}
	mockRepo.Pedidos = []*entities.Pedido{pedido}

	// Test com status inválido
	err := useCase.Run(context.Background(), 1, "StatusInvalido")

	// Assertions
	if err == nil {
		t.Fatal("expected error for invalid payment status, got nil")
	}
}

func TestPedidoAtualizarStatusPagamentoUseCase_Run_PedidoNotFound(t *testing.T) {
	mockRepo := &MockPedidoRepositoryAtualizarPagamento{}
	useCase := NewPedidoAtualizarStatusPagamentoUseCase(mockRepo)

	// Test sem pedidos no repositório
	err := useCase.Run(context.Background(), 999, "Pago")

	// Assertions
	if err == nil {
		t.Fatal("expected error for non-existent pedido, got nil")
	}
}

func TestPedidoAtualizarStatusPagamentoUseCase_Run_PaymentFlow(t *testing.T) {
	mockRepo := &MockPedidoRepositoryAtualizarPagamento{}
	useCase := NewPedidoAtualizarStatusPagamentoUseCase(mockRepo)

	// Setup pedido no repositório
	pedido := &entities.Pedido{
		ID:              1,
		ClienteNome:     "João Silva",
		StatusPagamento: "Pendente",
	}
	mockRepo.Pedidos = []*entities.Pedido{pedido}

	// Test typical payment flow: Pendente -> Pago
	err := useCase.Run(context.Background(), 1, "Pago")
	if err != nil {
		t.Fatalf("expected no error for 'Pago', got %v", err)
	}
	if mockRepo.Pedidos[0].StatusPagamento != "Pago" {
		t.Errorf("expected StatusPagamento 'Pago', got %s", mockRepo.Pedidos[0].StatusPagamento)
	}

	// Test refusal flow: reset to Pendente -> Recusado
	mockRepo.Pedidos[0].StatusPagamento = "Pendente"
	err = useCase.Run(context.Background(), 1, "Recusado")
	if err != nil {
		t.Fatalf("expected no error for 'Recusado', got %v", err)
	}
	if mockRepo.Pedidos[0].StatusPagamento != "Recusado" {
		t.Errorf("expected StatusPagamento 'Recusado', got %s", mockRepo.Pedidos[0].StatusPagamento)
	}
}
