package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"lanchonete/internal/interfaces/publisher"
)

type PedidoAtualizarStatusUseCase interface {
	Run(ctx context.Context, pedidoID int, novo_status string) error
}

type pedidoAtualizarStatusUseCase struct {
	pedidoGateway  repository.PedidoRepository
	eventPublisher publisher.EventPublisher
}

func NewPedidoAtualizarStatusUseCase(pedidoGateway repository.PedidoRepository, publisher publisher.EventPublisher) PedidoAtualizarStatusUseCase {
	return &pedidoAtualizarStatusUseCase{
		pedidoGateway:  pedidoGateway,
		eventPublisher: publisher,
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

	// ✨ Publicar evento no SQS
	payload := map[string]interface{}{
		"id_pedido":     pedidoID,
		"status":        status,
		"atualizado_em": pedido.UltimaAtualizacao,
	}

	return pduc.eventPublisher.Publish("pedido_status_atualizado", payload)
}
