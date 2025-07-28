package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/repository"
	"lanchonete/internal/interfaces/publisher"
)

type ProdutoRemoverUseCase interface {
	Run(ctx context.Context, id int) error
}

type produtoRemoverUseCase struct {
	produtoGateway repository.ProdutoRepository
	eventPublisher publisher.EventPublisher
}

func NewProdutoRemoverUseCase(
	produtoGateway repository.ProdutoRepository,
	eventPublisher publisher.EventPublisher,
) ProdutoRemoverUseCase {
	return &produtoRemoverUseCase{
		produtoGateway: produtoGateway,
		eventPublisher: eventPublisher,
	}
}

func (pruc *produtoRemoverUseCase) Run(c context.Context, id int) error {
	_, err := pruc.produtoGateway.BuscarProdutoPorId(c, id)
	if err != nil {
		return fmt.Errorf("produto não existe no banco de dados: %w", err)
	}

	err = pruc.produtoGateway.RemoverProduto(c, id)
	if err != nil {
		return fmt.Errorf("não foi possível remover o produto: %w", err)
	}

	// ✨ Publicar evento de remoção
	payload := map[string]interface{}{
		"id_produto": id,
	}

	err = pruc.eventPublisher.Publish("produto_removido", payload)
	if err != nil {
		fmt.Println("⚠️ Falha ao publicar evento de remoção do produto:", err)
	}

	return nil
}
