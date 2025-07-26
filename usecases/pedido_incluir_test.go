package usecases

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"testing"
	"time"
)

// MockPedidoRepository implements repository.PedidoRepository for testing
type MockPedidoRepositoryIncluir struct {
	Pedidos []*entities.Pedido
}

func (m *MockPedidoRepositoryIncluir) CriarPedido(ctx context.Context, pedido *entities.Pedido) error {
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

func (m *MockPedidoRepositoryIncluir) BuscarPedido(ctx context.Context, id int) (*entities.Pedido, error) {
	for _, p := range m.Pedidos {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("pedido não encontrado")
}

func (m *MockPedidoRepositoryIncluir) AtualizarStatusPedido(ctx context.Context, pedidoID int, status string, ultimaAtualizacao time.Time) error {
	for _, p := range m.Pedidos {
		if p.ID == pedidoID {
			p.Status = entities.StatusPedido(status)
			p.UltimaAtualizacao = ultimaAtualizacao
			return nil
		}
	}
	return errors.New("pedido não encontrado")
}

func (m *MockPedidoRepositoryIncluir) ListarTodosOsPedidos(ctx context.Context) ([]*entities.Pedido, error) {
	return m.Pedidos, nil
}

func (m *MockPedidoRepositoryIncluir) AtualizarStatusPagamento(ctx context.Context, pedidoID int, statusPagamento string, ultimaAtualizacao time.Time) error {
	for _, p := range m.Pedidos {
		if p.ID == pedidoID {
			p.StatusPagamento = statusPagamento
			p.UltimaAtualizacao = ultimaAtualizacao
			return nil
		}
	}
	return errors.New("pedido não encontrado")
}

func TestPedidoIncluirUseCase_Run_MultiplePedidos(t *testing.T) {
	mockRepo := &MockPedidoRepositoryIncluir{}
	useCase := NewPedidoIncluirUseCase(mockRepo)

	// Produtos base
	produtos := []entities.Produto{
		{Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Hamburguer artesanal", Preco: 25.0},
		{Nome: "Batata Frita", Categoria: entities.Acompanhamento, Descricao: "Batata frita crocante", Preco: 10.0},
		{Nome: "Refrigerante", Categoria: entities.Bebida, Descricao: "Coca-Cola lata", Preco: 7.5},
	}

	pedidos := []struct {
		ClienteNome    string
		Produtos       []entities.Produto
		Personalizacao *string
	}{
		{
			ClienteNome:    "João",
			Produtos:       []entities.Produto{produtos[0], produtos[1]},
			Personalizacao: nil,
		},
		{
			ClienteNome:    "Maria",
			Produtos:       []entities.Produto{produtos[0], produtos[2]},
			Personalizacao: nil,
		},
		{
			ClienteNome:    "Pedro",
			Produtos:       []entities.Produto{produtos[0], produtos[1], produtos[2]},
			Personalizacao: nil,
		},
	}

	for _, p := range pedidos {
		pedido, err := useCase.Run(context.Background(), p.ClienteNome, p.Produtos, p.Personalizacao)
		if err != nil {
			t.Fatalf("unexpected error for pedido %+v: %v\n", p, err)
		}
		if pedido == nil {
			t.Fatalf("expected pedido to be created for %+v\n", p)
		}
	}

	if len(mockRepo.Pedidos) != 3 {
		t.Errorf("expected 3 pedidos in repository, got %d", len(mockRepo.Pedidos))
	}

	// Check attributes of each created pedido
	for i, pedido := range mockRepo.Pedidos {
		expected := pedidos[i]
		if pedido.ClienteNome != expected.ClienteNome {
			t.Errorf("pedido cliente mismatch: got %+v, want %+v", pedido.ClienteNome, expected.ClienteNome)
		}
		if len(pedido.Produtos) != len(expected.Produtos) {
			t.Errorf("pedido produtos count mismatch: got %d, want %d", len(pedido.Produtos), len(expected.Produtos))
		}
	}
}

func TestPedidoIncluirUseCase_Run_WithPersonalizacao(t *testing.T) {
	mockRepo := &MockPedidoRepositoryIncluir{}
	useCase := NewPedidoIncluirUseCase(mockRepo)

	produtos := []entities.Produto{
		{Nome: "Hamburguer", Categoria: entities.Lanche, Descricao: "Hamburguer artesanal", Preco: 25.0},
	}

	personalizacao := "Sem cebola e com molho extra"
	pedido, err := useCase.Run(context.Background(), "João", produtos, &personalizacao)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pedido == nil {
		t.Fatal("expected pedido to be created")
	}
	if pedido.Personalizacao == nil || *pedido.Personalizacao != personalizacao {
		t.Errorf("expected personalizacao '%s', got %v", personalizacao, pedido.Personalizacao)
	}
}

func TestPedidoIncluirUseCase_Run_EmptyProductList(t *testing.T) {
	mockRepo := &MockPedidoRepositoryIncluir{}
	useCase := NewPedidoIncluirUseCase(mockRepo)

	pedido, err := useCase.Run(context.Background(), "João", []entities.Produto{}, nil)

	if err == nil {
		t.Fatal("expected error for empty product list, got nil")
	}
	if pedido != nil {
		t.Errorf("expected nil pedido for empty product list, got %+v", pedido)
	}
}
