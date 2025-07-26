package repository

import (
	"context"
	"errors"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"sync"
)

var (
	ErrPedidoNaoEncontrado = errors.New("pedido não encontrado")
	ErrIDObrigatorio       = errors.New("ID do pedido é obrigatório")
	ErrPedidoJaExiste     = errors.New("pedido já existe")
)

// PedidoRepositoryMemoria implementa a interface PedidoRepository usando armazenamento em memória
type PedidoRepositoryMemoria struct {
	pedidos map[string]*entities.Order
	mutex   sync.RWMutex
}

// NovoPedidoRepository cria uma nova instância do PedidoRepositoryMemoria
func NovoPedidoRepository() repository.OrderRepository {
	return &PedidoRepositoryMemoria{
		pedidos: make(map[string]*entities.Order),
	}
}

// CriarPedido implementa o método CriarPedido do PedidoRepository
func (r *PedidoRepositoryMemoria) CriarPedido(ctx context.Context, pedido *entities.Order) error {
	if err := r.validarContexto(ctx); err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if pedido.ID == "" {
		return ErrIDObrigatorio
	}

	if _, existe := r.pedidos[pedido.ID]; existe {
		return ErrPedidoJaExiste
	}

	r.pedidos[pedido.ID] = r.clonarPedido(pedido)
	return nil
}

// BuscarPedidoPorID implementa o método BuscarPedidoPorID do PedidoRepository
func (r *PedidoRepositoryMemoria) BuscarPedidoPorID(ctx context.Context, id string) (*entities.Order, error) {
	if err := r.validarContexto(ctx); err != nil {
		return nil, err
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pedido, existe := r.pedidos[id]
	if !existe {
		return nil, ErrPedidoNaoEncontrado
	}

	return r.clonarPedido(pedido), nil
}

// AtualizarPedido implementa o método AtualizarPedido do PedidoRepository
func (r *PedidoRepositoryMemoria) AtualizarPedido(ctx context.Context, pedido *entities.Order) error {
	if err := r.validarContexto(ctx); err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, existe := r.pedidos[pedido.ID]; !existe {
		return ErrPedidoNaoEncontrado
	}

	r.pedidos[pedido.ID] = r.clonarPedido(pedido)
	return nil
}

// DeletarPedido implementa o método DeletarPedido do PedidoRepository
func (r *PedidoRepositoryMemoria) DeletarPedido(ctx context.Context, id string) error {
	if err := r.validarContexto(ctx); err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, existe := r.pedidos[id]; !existe {
		return ErrPedidoNaoEncontrado
	}

	delete(r.pedidos, id)
	return nil
}

// ListarTodosPedidos implementa o método ListarTodosPedidos do PedidoRepository
func (r *PedidoRepositoryMemoria) ListarTodosPedidos(ctx context.Context) ([]*entities.Order, error) {
	if err := r.validarContexto(ctx); err != nil {
		return nil, err
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pedidos := make([]*entities.Order, 0, len(r.pedidos))
	for _, pedido := range r.pedidos {
		pedidos = append(pedidos, r.clonarPedido(pedido))
	}

	return pedidos, nil
}

// Métodos auxiliares privados
func (r *PedidoRepositoryMemoria) validarContexto(ctx context.Context) error {
	if ctx == nil {
		return errors.New("contexto é obrigatório")
	}
	return nil
}

func (r *PedidoRepositoryMemoria) clonarPedido(pedido *entities.Order) *entities.Order {
	if pedido == nil {
		return nil
	}
	clone := *pedido
	return &clone
} 