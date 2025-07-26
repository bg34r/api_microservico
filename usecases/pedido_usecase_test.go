package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"time"
)

// MockPedidoRepository implements repository.PedidoRepository for testing
// This is a shared mock that can be used by multiple test files if needed
type MockPedidoRepository struct {
	Pedidos []*entities.Pedido
}

func (m *MockPedidoRepository) CriarPedido(ctx context.Context, pedido *entities.Pedido) error {
	// Simulate duplicate check
	for _, p := range m.Pedidos {
		if pedido.ID == 0 {
			pedido.ID = len(m.Pedidos) + 1
		}

		if p.ID == pedido.ID {
			return errors.New("pedido já existe")
		}
	}
	m.Pedidos = append(m.Pedidos, pedido)
	return nil
}

func (m *MockPedidoRepository) BuscarPedido(ctx context.Context, id int) (*entities.Pedido, error) {
	for _, p := range m.Pedidos {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("pedido não encontrado")
}

func (m *MockPedidoRepository) AtualizarStatusPedido(ctx context.Context, pedidoID int, status string, ultimaAtualizacao time.Time) error {
	for _, p := range m.Pedidos {
		if p.ID == pedidoID {
			p.Status = entities.StatusPedido(status)
			p.UltimaAtualizacao = ultimaAtualizacao
			return nil
		}
	}
	return errors.New("pedido não encontrado")
}

func (m *MockPedidoRepository) ListarTodosOsPedidos(ctx context.Context) ([]*entities.Pedido, error) {
	return m.Pedidos, nil
}

func (m *MockPedidoRepository) AtualizarStatusPagamento(ctx context.Context, pedidoID int, statusPagamento string, ultimaAtualizacao time.Time) error {
	for _, p := range m.Pedidos {
		if p.ID == pedidoID {
			p.StatusPagamento = statusPagamento
			p.UltimaAtualizacao = ultimaAtualizacao
			return nil
		}
	}
	return errors.New("pedido não encontrado")
}
