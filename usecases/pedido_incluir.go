package usecases

import (
	"context"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type PedidoIncluirUseCase interface {
	Run(ctx context.Context, clienteNome string, produtos []entities.Produto, personalizacao *string) (*entities.Pedido, error)
}

type pedidoIncluirUseCase struct {
	pedidoRepository repository.PedidoRepository
}

func NewPedidoIncluirUseCase(pedidoRepository repository.PedidoRepository) PedidoIncluirUseCase {
	return &pedidoIncluirUseCase{
		pedidoRepository: pedidoRepository,
	}
}

func (pduc *pedidoIncluirUseCase) Run(c context.Context, clienteNome string, produtos []entities.Produto, personalizacao *string) (*entities.Pedido, error) {
	pedido, err := entities.PedidoNew(clienteNome, produtos, personalizacao)
	if err != nil {
		return nil, err
	}
	err = pduc.pedidoRepository.CriarPedido(c, pedido)
	if err != nil {
		return nil, err
	}
	return pedido, nil
}
