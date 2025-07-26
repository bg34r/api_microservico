package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type ProdutoBuscaPorIdUseCase interface {
	Run(ctx context.Context, id int) (*entities.Produto, error)
}

type produtoBuscaPorIdUseCase struct {
	produtoRepo repository.ProdutoRepository
}

func NewProdutoBuscaPorIdUseCase(produtoRepo repository.ProdutoRepository) ProdutoBuscaPorIdUseCase {
	return &produtoBuscaPorIdUseCase{
		produtoRepo: produtoRepo,
	}
}

func (pd *produtoBuscaPorIdUseCase) Run(c context.Context, id int) (*entities.Produto, error) {
	produto, err := pd.produtoRepo.BuscarProdutoPorId(c, id)

	if err != nil {
		return nil, fmt.Errorf("não foi possível buscar produto: %w", err)
	}
	return produto, nil
}
