package repository

import (
	"context"
	"lanchonete/internal/domain/entities"
)

type ProdutoRepository interface {
	AdicionarProduto(c context.Context, produto *entities.Produto) error
	BuscarProdutoPorId(c context.Context, id int) (*entities.Produto, error)
	ListarTodosOsProdutos(c context.Context) ([]*entities.Produto, error)
	EditarProduto(c context.Context, produto *entities.Produto) error
	RemoverProduto(c context.Context, id int) error
	ListarPorCategoria(c context.Context, categoria string) ([]*entities.Produto, error)
}
