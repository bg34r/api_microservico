package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type ProdutoIncluirUseCase interface {
	Run(ctx context.Context, nome, categoria, descricao string, preco float32) (*entities.Produto, error)
}

type produtoIncluirUseCase struct {
	produtoRepository repository.ProdutoRepository
}

func NewProdutoIncluirUseCase(produtoRepository repository.ProdutoRepository) ProdutoIncluirUseCase {
	return &produtoIncluirUseCase{
		produtoRepository: produtoRepository,
	}
}

func (pd *produtoIncluirUseCase) Run(c context.Context, nome string, categoria string, descricao string, preco float32) (*entities.Produto, error) {

	produto, err := entities.ProdutoNew(nome, categoria, descricao, preco)

	if err != nil {
		return nil, fmt.Errorf("criação de produto inválida: %w", err)
	}

	err = pd.produtoRepository.AdicionarProduto(c, produto)
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar produto: %w", err)
	}

	return produto, nil
}
