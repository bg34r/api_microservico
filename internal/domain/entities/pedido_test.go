package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPedidoNew_Success(t *testing.T) {
	produtos := []Produto{
		{ID: 1, Nome: "Big Mac", Categoria: Lanche, Preco: 25.0},
		{ID: 2, Nome: "Coca Cola", Categoria: Bebida, Preco: 5.0},
	}
	personalizacao := "Sem cebola"

	pedido, err := PedidoNew("João Silva", produtos, &personalizacao)

	assert.NoError(t, err)
	assert.NotNil(t, pedido)
	assert.Equal(t, "João Silva", pedido.ClienteNome)
	assert.Equal(t, Pendente, pedido.Status)
	assert.Equal(t, "Pendente", pedido.StatusPagamento)
	assert.Equal(t, float32(30.0), pedido.Total)
	assert.Equal(t, &personalizacao, pedido.Personalizacao)
	assert.Equal(t, produtos, pedido.Produtos)
	assert.NotZero(t, pedido.UltimaAtualizacao)
}

func TestPedidoNew_ErrorSemProdutos(t *testing.T) {
	produtos := []Produto{}

	pedido, err := PedidoNew("João Silva", produtos, nil)

	assert.Error(t, err)
	assert.Nil(t, pedido)
	assert.Equal(t, "o pedido precisa ter ao menos um produto", err.Error())
}

func TestPedidoNew_ErrorSemLanche(t *testing.T) {
	produtos := []Produto{
		{ID: 1, Nome: "Coca Cola", Categoria: Bebida, Preco: 5.0},
		{ID: 2, Nome: "Batata Frita", Categoria: Acompanhamento, Preco: 8.0},
	}

	pedido, err := PedidoNew("João Silva", produtos, nil)

	assert.Error(t, err)
	assert.Nil(t, pedido)
	assert.Equal(t, "o pedido precisa ter ao menos um lanche", err.Error())
}

func TestPedido_UpdateStatus_Success(t *testing.T) {
	pedido := &Pedido{
		Status:            Pendente,
		UltimaAtualizacao: time.Now().Add(-time.Hour),
	}
	oldTime := pedido.UltimaAtualizacao

	err := pedido.UpdateStatus(Recebido)

	assert.NoError(t, err)
	assert.Equal(t, Recebido, pedido.Status)
	assert.True(t, pedido.UltimaAtualizacao.After(oldTime))
}

func TestPedido_UpdateStatus_AllValidStatuses(t *testing.T) {
	pedido := &Pedido{Status: Pendente}

	validStatuses := []StatusPedido{Pendente, Recebido, EmPreparacao, Pronto, Finalizado}

	for _, status := range validStatuses {
		err := pedido.UpdateStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, status, pedido.Status)
	}
}

func TestPedido_UpdateStatus_InvalidStatus(t *testing.T) {
	pedido := &Pedido{Status: Pendente}

	err := pedido.UpdateStatus(StatusPedido("StatusInvalido"))

	assert.Error(t, err)
	assert.Equal(t, "status inválido", err.Error())
	assert.Equal(t, Pendente, pedido.Status) // Status não deve mudar
}

func TestPedido_UpdateStatusPagamento_Success(t *testing.T) {
	pedido := &Pedido{
		StatusPagamento:   "Pendente",
		UltimaAtualizacao: time.Now().Add(-time.Hour),
	}
	oldTime := pedido.UltimaAtualizacao

	err := pedido.UpdateStatusPagamento("Pago")

	assert.NoError(t, err)
	assert.Equal(t, "Pago", pedido.StatusPagamento)
	assert.True(t, pedido.UltimaAtualizacao.After(oldTime))
}

func TestPedido_UpdateStatusPagamento_AllValidStatuses(t *testing.T) {
	pedido := &Pedido{StatusPagamento: "Pendente"}

	validStatuses := []string{"Pendente", "Pago", "Recusado", "Cancelado"}

	for _, status := range validStatuses {
		err := pedido.UpdateStatusPagamento(status)
		assert.NoError(t, err)
		assert.Equal(t, status, pedido.StatusPagamento)
	}
}

func TestPedido_UpdateStatusPagamento_InvalidStatus(t *testing.T) {
	pedido := &Pedido{StatusPagamento: "Pendente"}

	err := pedido.UpdateStatusPagamento("StatusInvalido")

	assert.Error(t, err)
	assert.Equal(t, "status de pagamento inválido", err.Error())
	assert.Equal(t, "Pendente", pedido.StatusPagamento) // Status não deve mudar
}
