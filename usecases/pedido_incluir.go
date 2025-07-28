package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"lanchonete/internal/interfaces/publisher"
)

type PedidoIncluirUseCase interface {
	Run(ctx context.Context, clienteNome string, produtos []entities.Produto, personalizacao *string) (*entities.Pedido, error)
}

type pedidoIncluirUseCase struct {
	pedidoRepository repository.PedidoRepository
	eventPublisher   publisher.EventPublisher
}

func NewPedidoIncluirUseCase(pedidoRepository repository.PedidoRepository, publisher publisher.EventPublisher) PedidoIncluirUseCase {
	return &pedidoIncluirUseCase{
		pedidoRepository: pedidoRepository,
		eventPublisher:   publisher,
	}
}

func serializeProdutos(produtos []entities.Produto) []map[string]interface{} {
	var lista []map[string]interface{}
	for _, p := range produtos {
		lista = append(lista, map[string]interface{}{
			"id":    p.ID,
			"nome":  p.Nome,
			"preco": p.Preco,
		})
	}
	return lista
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

	// ✨ Publicar evento "pedido_criado"
	payload := map[string]interface{}{
		"id_pedido":      pedido.ID,
		"cliente":        clienteNome,
		"status":         pedido.Status,
		"personalizacao": personalizacao,
		"criado_em":      pedido.UltimaAtualizacao,
		"produtos":       serializeProdutos(produtos), // transformar []Produto em dados simples
	}

	err = pduc.eventPublisher.Publish("pedido_criado", payload)
	if err != nil {
		fmt.Println("⚠️ Falha ao publicar evento de pedido:", err)
	}

	return pedido, nil
}
