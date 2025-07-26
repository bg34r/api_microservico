package repository

import (
	"context"
	"lanchonete/internal/domain/entities"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	CriarPedido(ctx context.Context, order *entities.Order) error
	BuscarPedidoPorID(ctx context.Context, id string) (*entities.Order, error)
	AtualizarPedido(ctx context.Context, order *entities.Order) error
	DeletarPedido(ctx context.Context, id string) error
	ListarTodosPedidos(ctx context.Context) ([]*entities.Order, error)
} 