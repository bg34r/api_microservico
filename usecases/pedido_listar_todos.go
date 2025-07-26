package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type PedidoListarTodosUseCase interface {
	Run(ctx context.Context) ([]*entities.Pedido, error)
}

type pedidoListarTodosUseCase struct {
	pedidoRepo repository.PedidoRepository
}

func NewPedidoListarTodosUseCase(pedidoRepo repository.PedidoRepository) PedidoListarTodosUseCase {
	return &pedidoListarTodosUseCase{
		pedidoRepo: pedidoRepo,
	}
}

func (pd *pedidoListarTodosUseCase) Run(c context.Context) ([]*entities.Pedido, error) {
	pedidos, err := pd.pedidoRepo.ListarTodosOsPedidos(c)
	if err != nil {
		return nil, fmt.Errorf("não foi possível listar pedidos: %w", err)
	}
	return pedidos, nil
}
