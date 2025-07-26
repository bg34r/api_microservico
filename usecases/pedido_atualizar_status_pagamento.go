package usecases

import (
	"context"
	"lanchonete/internal/domain/repository"
)

type PedidoAtualizarStatusPagamentoUseCase interface {
	Run(ctx context.Context, pedidoID int, statusPagamento string) error
}

type pedidoAtualizarStatusPagamentoUseCase struct {
	pedidoGateway repository.PedidoRepository
}

func NewPedidoAtualizarStatusPagamentoUseCase(pedidoGateway repository.PedidoRepository) PedidoAtualizarStatusPagamentoUseCase {
	return &pedidoAtualizarStatusPagamentoUseCase{
		pedidoGateway: pedidoGateway,
	}
}

func (pduc *pedidoAtualizarStatusPagamentoUseCase) Run(c context.Context, pedidoID int, statusPagamento string) error {
	// Buscar o pedido para validar se existe e para pegar o timestamp atual
	pedido, err := pduc.pedidoGateway.BuscarPedido(c, pedidoID)
	if err != nil {
		return err
	}

	// Validar o status de pagamento usando o método da entidade
	err = pedido.UpdateStatusPagamento(statusPagamento)
	if err != nil {
		return err
	}

	// Atualizar no banco de dados
	err = pduc.pedidoGateway.AtualizarStatusPagamento(c, pedidoID, statusPagamento, pedido.UltimaAtualizacao)
	if err != nil {
		return err
	}

	return nil
}
