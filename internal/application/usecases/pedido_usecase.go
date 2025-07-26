package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"lanchonete/usecases"
	"time"
)

type pedidoUseCase struct {
	pedidoRepository repository.PedidoRepository
}

func NewPedidoUseCase(pedidoRepository repository.PedidoRepository) usecases.PedidoUseCase {
	return &pedidoUseCase{
		pedidoRepository: pedidoRepository,
	}
}

func (puc *pedidoUseCase) CriarPedido(c context.Context, pedido *entities.Pedido) error {
	return puc.pedidoRepository.CriarPedido(c, pedido)
}

func (puc *pedidoUseCase) BuscarPedido(c context.Context, identificacao int) (*entities.Pedido, error) {
	return puc.pedidoRepository.BuscarPedido(c, identificacao)
}

func (puc *pedidoUseCase) AtualizarStatusPedido(c context.Context, identificacao int, status string, ultimaAtualizacao time.Time) error {
	return puc.pedidoRepository.AtualizarStatusPedido(c, identificacao, status, ultimaAtualizacao)
}
