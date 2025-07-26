package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type ProdutoListarPorCategoriaUseCase interface {
	Run(ctx context.Context, categoria string) ([]*entities.Produto, error)
}
type produtoListarPorCategoriaUseCase struct {
	produtoRepo repository.ProdutoRepository
}

func NewProdutoListarPorCategoriaUseCase(produtoRepo repository.ProdutoRepository) ProdutoListarPorCategoriaUseCase {
	return &produtoListarPorCategoriaUseCase{
		produtoRepo: produtoRepo,
	}
}

func (pd *produtoListarPorCategoriaUseCase) Run(c context.Context, categoria string) ([]*entities.Produto, error) {
	produtos, err := pd.produtoRepo.ListarPorCategoria(c, categoria)
	if err != nil {
		return nil, fmt.Errorf("não foi possível listar produtos: %w", err)
	}
	return produtos, nil
}
