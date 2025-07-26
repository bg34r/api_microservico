package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type ProdutoListarTodosUseCase interface {
	Run(ctx context.Context) ([]*entities.Produto, error)
}

type produtoListarTodosUseCase struct {
	produtoRepo repository.ProdutoRepository
}

func NewProdutoListarTodosUseCase(produtoRepo repository.ProdutoRepository) ProdutoListarTodosUseCase {
	return &produtoListarTodosUseCase{
		produtoRepo: produtoRepo,
	}
}

func (pd *produtoListarTodosUseCase) Run(c context.Context) ([]*entities.Produto, error) {
	produtos, err := pd.produtoRepo.ListarTodosOsProdutos(c)
	if err != nil {
		return nil, fmt.Errorf("não foi possível listar produtos: %w", err)
	}
	return produtos, nil
}
