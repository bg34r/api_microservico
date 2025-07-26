package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/repository"
)

type ProdutoRemoverUseCase interface {
	Run(ctx context.Context, id int) error
}

type produtoRemoverUseCase struct {
	produtoGateway repository.ProdutoRepository
}

func NewProdutoRemoverUseCase(produtoGateway repository.ProdutoRepository) ProdutoRemoverUseCase {
	return &produtoRemoverUseCase{
		produtoGateway: produtoGateway,
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

	return nil
}
