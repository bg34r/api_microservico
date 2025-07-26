package handler

import (
	"encoding/json"
	"fmt"
	_ "lanchonete/docs"
	"lanchonete/internal/domain/entities"
	response "lanchonete/internal/interfaces/http/responses"
	"lanchonete/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PedidoHandler struct {
	PedidoIncluirUseCase                  usecases.PedidoIncluirUseCase
	PedidoBuscarPorIdUseCase              usecases.PedidoBuscarPorIdUseCase
	PedidoAtualizarStatusUseCase          usecases.PedidoAtualizarStatusUseCase
	PedidoAtualizarStatusPagamentoUseCase usecases.PedidoAtualizarStatusPagamentoUseCase
	ProdutoBuscarPorIdUseCase             usecases.ProdutoBuscaPorIdUseCase
	PedidoListarTodosUseCase              usecases.PedidoListarTodosUseCase
}

func NewPedidoHandler(pedidoIncluirUseCase usecases.PedidoIncluirUseCase,
	pedidoBuscarPorIdUseCase usecases.PedidoBuscarPorIdUseCase,
	pedidoAtualizarStatusUsecase usecases.PedidoAtualizarStatusUseCase,
	pedidoAtualizarStatusPagamentoUseCase usecases.PedidoAtualizarStatusPagamentoUseCase,
	produtoBuscarPorIdUseCase usecases.ProdutoBuscaPorIdUseCase,
	pedidoListarTodosUseCase usecases.PedidoListarTodosUseCase) *PedidoHandler {
	return &PedidoHandler{
		PedidoIncluirUseCase:                  pedidoIncluirUseCase,
		PedidoBuscarPorIdUseCase:              pedidoBuscarPorIdUseCase,
		PedidoAtualizarStatusUseCase:          pedidoAtualizarStatusUsecase,
		PedidoAtualizarStatusPagamentoUseCase: pedidoAtualizarStatusPagamentoUseCase,
		ProdutoBuscarPorIdUseCase:             produtoBuscarPorIdUseCase,
		PedidoListarTodosUseCase:              pedidoListarTodosUseCase,
	}
}

// CriarPedido godoc
// @Summary Cria um pedido
// @Description Cria um pedido
// @Tags pedido
// @Router /pedidos [post]
// @Accept  json
// @Produce  json
// @Param pedido body entities.Pedido true "Pedido"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) CriarPedido(r *gin.Context) {
	var pedido entities.Pedido
	fmt.Println("Handler Criando pedido", pedido)
	err := json.NewDecoder(r.Request.Body).Decode(&pedido)
	fmt.Println("Handler Criando Depois pedido", pedido)
	if err != nil {
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Substituir o array de produtos com os dados completos do banco
	produtosCompletos := []entities.Produto{}

	for _, produto := range pedido.Produtos {
		pBanco, err := h.ProdutoBuscarPorIdUseCase.Run(r, produto.ID)
		if err != nil {
			r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Produto não Cadastrado!"})
			return
		}
		produtosCompletos = append(produtosCompletos, *pBanco)
	}

	// Chamar PedidoNew com os produtos completos
	ped, err := h.PedidoIncluirUseCase.Run(r, pedido.ClienteNome, produtosCompletos, pedido.Personalizacao)
	if err != nil {
		r.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	r.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Pedido criado com sucesso" + strconv.Itoa(ped.ID),
	})
}

// BuscarPedido godoc
// @Summary Busca um pedido
// @Description Busca um pedido
// @Tags pedido
// @Router /pedidos/{ID} [get]
// @Accept  json
// @Produce  json
// @Param ID path string true "Número do pedido"
// @Success 200 {object} entities.Pedido
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) BuscarPedido(r *gin.Context) {
	nroPedido := r.Param("nroPedido")
	id, err := strconv.Atoi(nroPedido)
	if err != nil {
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Número do pedido inválido"})
		return
	}
	pedido, err := h.PedidoBuscarPorIdUseCase.Run(r, id)
	if err != nil {
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	r.JSON(http.StatusOK, pedido)

}

// AtualizarPedido godoc
// @Summary Atualiza um pedido a partir de sua Identificação
// @Description Atualizar um pedido
// @Tags pedido
// @Router /pedidos/{nroPedido}/status/{status} [put]
// @Accept  json
// @Produce  json
// @Param nroPedido path string true "Número do pedido"
// @Param status path string true "Novo Status do pedido"
// @Success 200 {object} entities.Pedido
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) AtualizarStatusPedido(r *gin.Context) {
	nroPedido := r.Param("nroPedido")
	id, err := strconv.Atoi(nroPedido)
	if err != nil {
		fmt.Printf("Erro ao converter ID do pedido: %v\n", err)
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Número do pedido inválido"})
		return
	}

	status := r.Param("status")
	fmt.Printf("Atualizando pedido ID: %d para status: '%s'\n", id, status)

	err = h.PedidoAtualizarStatusUseCase.Run(r, id, status)
	if err != nil {
		fmt.Printf("Erro ao atualizar status: %v\n", err)
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	fmt.Printf("Status do pedido %d atualizado com sucesso para '%s'\n", id, status)
	r.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Status do pedido atualizado com sucesso",
	})
}

// AtualizarStatusPagamento godoc
// @Summary Atualiza o status de pagamento de um pedido
// @Description Atualizar o status de pagamento de um pedido
// @Tags pedido
// @Router /pedidos/{nroPedido}/pagamento/{statusPagamento} [put]
// @Accept  json
// @Produce  json
// @Param nroPedido path string true "Número do pedido"
// @Param statusPagamento path string true "Novo Status de pagamento (Pendente, Pago, Recusado, Cancelado)"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) AtualizarStatusPagamento(r *gin.Context) {
	nroPedido := r.Param("nroPedido")
	id, err := strconv.Atoi(nroPedido)
	if err != nil {
		fmt.Printf("Erro ao converter ID do pedido: %v\n", err)
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Número do pedido inválido"})
		return
	}

	statusPagamento := r.Param("statusPagamento")
	fmt.Printf("Atualizando status de pagamento do pedido ID: %d para status: '%s'\n", id, statusPagamento)

	err = h.PedidoAtualizarStatusPagamentoUseCase.Run(r, id, statusPagamento)
	if err != nil {
		fmt.Printf("Erro ao atualizar status de pagamento: %v\n", err)
		r.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	fmt.Printf("Status de pagamento do pedido %d atualizado com sucesso para '%s'\n", id, statusPagamento)
	r.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Status de pagamento atualizado com sucesso",
	})
}

// ProdutoListarTodos godoc
// @Summary Lista todos os pedidos no banco
// @Description Lista todos os pedidos presentes no banco
// @Tags pedido
// @Router /pedidos/listartodos [GET]
// @Accept  json
// @Produce  json
// @Success 200 {object} []entities.Pedido
// @Failure 400 {object} response.ErrorResponse
func (h *PedidoHandler) ListarTodosOsPedidos(r *gin.Context) {
	pedidos, err := h.PedidoListarTodosUseCase.Run(r)
	if err != nil {
		r.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	r.JSON(http.StatusOK, pedidos)
}
