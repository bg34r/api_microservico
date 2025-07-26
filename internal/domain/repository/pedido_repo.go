package repository

import (
	"context"
	"lanchonete/internal/domain/entities"
	"time"
)

// PedidoRepository define a interface para operações de dados de pedidos
type PedidoRepository interface {
	CriarPedido(c context.Context, pedido *entities.Pedido) error
	BuscarPedido(c context.Context, pedidoID int) (*entities.Pedido, error)
	AtualizarStatusPedido(c context.Context, pedidoID int, status string, UltimaAtualizacao time.Time) error
	AtualizarStatusPagamento(c context.Context, pedidoID int, statusPagamento string, UltimaAtualizacao time.Time) error
	ListarTodosOsPedidos(c context.Context) ([]*entities.Pedido, error)
}
