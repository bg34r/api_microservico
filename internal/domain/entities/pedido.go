package entities

import (
	"errors"
	"fmt"
	"time"
)

type StatusPedido string

const (
	Pendente     StatusPedido = "Pendente"
	Recebido     StatusPedido = "Recebido"
	EmPreparacao StatusPedido = "Em preparação"
	Pronto       StatusPedido = "Pronto"
	Finalizado   StatusPedido = "Finalizado"
)

type Pedido struct {
	ID                int          `json:"id,omitempty"`
	ClienteNome       string       `json:"cliente_nome,omitempty"` // Opcional: apenas nome do cliente
	Status            StatusPedido `json:"status"`
	StatusPagamento   string       `json:"status_pagamento"`
	TimeStamp         string       `json:"time_stamp"`
	UltimaAtualizacao time.Time    `json:"ultima_atualizacao"`
	Total             float32      `json:"total"`
	Personalizacao    *string      `json:"personalizacao,omitempty"` // Personalização específica do pedido
	Produtos          []Produto    `json:"produtos"`
}

func PedidoNew(clienteNome string, produtos []Produto, personalizacao *string) (*Pedido, error) {
	fmt.Println("Pedido Entity: ", produtos)
	if len(produtos) == 0 {
		return nil, errors.New("o pedido precisa ter ao menos um produto")
	}

	temLanche := false
	total := float32(0)
	for _, produto := range produtos {
		total += produto.Preco
		if produto.Categoria == Lanche {
			temLanche = true
		}
	}

	if !temLanche {
		return nil, errors.New("o pedido precisa ter ao menos um lanche")
	}

	now := time.Now()

	return &Pedido{
		ClienteNome:       clienteNome,
		Status:            Pendente,
		StatusPagamento:   "Pendente",
		TimeStamp:         "00:15:00",
		UltimaAtualizacao: now,
		Total:             total,
		Personalizacao:    personalizacao,
		Produtos:          produtos,
	}, nil
}

func (p *Pedido) UpdateStatus(status StatusPedido) error {
	switch status {
	case Pendente, Recebido, EmPreparacao, Pronto, Finalizado:
		p.Status = status
		p.UltimaAtualizacao = time.Now()
		return nil
	default:
		return errors.New("status inválido")
	}
}

func (p *Pedido) UpdateStatusPagamento(statusPagamento string) error {
	// Validar os status de pagamento válidos
	switch statusPagamento {
	case "Pendente", "Pago", "Recusado", "Cancelado":
		p.StatusPagamento = statusPagamento
		p.UltimaAtualizacao = time.Now()
		return nil
	default:
		return errors.New("status de pagamento inválido")
	}
}
