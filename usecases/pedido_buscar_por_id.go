package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type PedidoBuscarPorIdUseCase interface {
	Run(ctx context.Context, pedidoID int) (*entities.Pedido, error)
}

type pedidoBuscarPorIdUseCase struct {
	pedidoRepository repository.PedidoRepository
}

func NewPedidoBuscarPorIdUseCase(pedidoRepository repository.PedidoRepository) PedidoBuscarPorIdUseCase {
	return &pedidoBuscarPorIdUseCase{
		pedidoRepository: pedidoRepository,
	}
}

func (pduc *pedidoBuscarPorIdUseCase) Run(c context.Context, pedidoID int) (*entities.Pedido, error) {

	pedido, err := pduc.pedidoRepository.BuscarPedido(c, pedidoID)
	if err != nil {
		return nil, err
	}
	return pedido, nil
}
