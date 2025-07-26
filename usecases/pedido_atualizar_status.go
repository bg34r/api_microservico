package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type PedidoAtualizarStatusUseCase interface {
	Run(ctx context.Context, pedidoID int, novo_status string) error
}

type pedidoAtualizarStatusUseCase struct {
	pedidoGateway repository.PedidoRepository
}

func NewPedidoAtualizarStatusUseCase(pedidoGateway repository.PedidoRepository) PedidoAtualizarStatusUseCase {
	return &pedidoAtualizarStatusUseCase{
		pedidoGateway: pedidoGateway,
	}
}

func (pduc *pedidoAtualizarStatusUseCase) Run(c context.Context, pedidoID int, status string) error {

	pedido, err := pduc.pedidoGateway.BuscarPedido(c, pedidoID)
	if err != nil {
		return err
	}

	err = pedido.UpdateStatus(entities.StatusPedido(status))
	if err != nil {
		return err
	}

	err = pduc.pedidoGateway.AtualizarStatusPedido(c, pedidoID, status, pedido.UltimaAtualizacao)

	if err != nil {
		return err
	}
	return nil
}
