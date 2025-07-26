package entities

import (
	"time"
)

// Order represents the core order entity in the domain
type Order struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	Items       []Item    `json:"items"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Item represents an item in an order
type Item struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// NovoPedido creates a new Order instance
func NovoPedido(customerID string, items []Item) *Order {
	return &Order{
		CustomerID:  customerID,
		Items:       items,
		Status:      "pendente",
		TotalAmount: calcularTotal(items),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// calcularTotal calculates the total amount of the order
func calcularTotal(items []Item) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// AtualizarStatus updates the status of the order
func (o *Order) AtualizarStatus(status string) {
	o.Status = status
	o.UpdatedAt = time.Now()
} 