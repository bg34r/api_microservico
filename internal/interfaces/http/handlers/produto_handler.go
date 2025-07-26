package handler

import (
	"fmt"
	_ "lanchonete/docs"
	"lanchonete/internal/application/presenters"
	"lanchonete/internal/domain/entities"
	response "lanchonete/internal/interfaces/http/responses"
	"lanchonete/usecases"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProdutoHandler struct {
	ProdutoIncluirUseCase            usecases.ProdutoIncluirUseCase
	ProdutoBuscarPorIdUseCase        usecases.ProdutoBuscaPorIdUseCase
	ProdutoListarTodosUseCase        usecases.ProdutoListarTodosUseCase
	ProdutoEditarUseCase             usecases.ProdutoEditarUseCase
	ProdutoRemoverUseCase            usecases.ProdutoRemoverUseCase
	ProdutoListarPorCategoriaUseCase usecases.ProdutoListarPorCategoriaUseCase
}

func NewProdutoHandler(produtoIncluirUseCase usecases.ProdutoIncluirUseCase,
	produtoBuscarPorIdUseCase usecases.ProdutoBuscaPorIdUseCase,
	produtoListarTodosUseCase usecases.ProdutoListarTodosUseCase,
	produtoEditarUseCase usecases.ProdutoEditarUseCase,
	produtoRemoverUseCase usecases.ProdutoRemoverUseCase,
	produtoListarPorCategoriaUseCase usecases.ProdutoListarPorCategoriaUseCase) *ProdutoHandler {
	return &ProdutoHandler{
		ProdutoIncluirUseCase:            produtoIncluirUseCase,
		ProdutoBuscarPorIdUseCase:        produtoBuscarPorIdUseCase,
		ProdutoEditarUseCase:             produtoEditarUseCase,
		ProdutoListarTodosUseCase:        produtoListarTodosUseCase,
		ProdutoRemoverUseCase:            produtoRemoverUseCase,
		ProdutoListarPorCategoriaUseCase: produtoListarPorCategoriaUseCase,
	}
}

// CriarProduto godoc
// @Summary Cria um produto
// @Description Cria um produto
// @Tags produto
// @Router /produto [post]
// @Accept  json
// @Produce  json
// @Param produto body entities.Produto true "Produto"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
func (ph *ProdutoHandler) ProdutoIncluir(c *gin.Context) {

	var produto entities.Produto

	err := c.ShouldBindJSON(&produto)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	prd, err := ph.ProdutoIncluirUseCase.Run(c, produto.Nome, string(produto.Categoria), produto.Descricao, produto.Preco)
	fmt.Println("Entrando no if erro Handler")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Produto incluído com sucesso: " + prd.Nome,
	})
}

// BuscarProduto godoc
// @Summary Busca um produto
// @Description Busca um produto
// @Tags produto
// @Router /produto/{id} [get]
// @Accept  json
// @Produce  json
// @Param id path int true "ID do produto"
// @Success 200 {object} presenters.ProdutoDTO
// @Failure 400 {object} response.ErrorResponse
func (ph *ProdutoHandler) ProdutoBuscarPorId(c *gin.Context) {
	id := c.Param("id")

	idnt, _ := strconv.Atoi(id)
	prd, err := ph.ProdutoBuscarPorIdUseCase.Run(c, idnt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "ID inválido"})
		return
	}

	prdDTO := presenters.NewProdutoDTO(prd)

	c.JSON(http.StatusOK, prdDTO)
}

// ProdutoListarTodos godoc
// @Summary Lista todos os produtos no banco
// @Description Lista todos os produtos cadastrados
// @Tags produto
// @Router /produtos [GET]
// @Accept  json
// @Produce  json
// @Success 200 {object} []presenters.ProdutoDTO
// @Failure 400 {object} response.ErrorResponse
func (ph *ProdutoHandler) ProdutoListarTodos(c *gin.Context) {
	produtos, err := ph.ProdutoListarTodosUseCase.Run(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Convert to DTOs
	var produtosDTO []*presenters.ProdutoDTO
	for _, produto := range produtos {
		produtosDTO = append(produtosDTO, presenters.NewProdutoDTO(produto))
	}

	c.JSON(http.StatusOK, produtosDTO)
}

// EditarProduto godoc
// @Summary Edita um produto
// @Description Edita um produto existente pelo ID
// @Tags produto
// @Accept  json
// @Produce  json
// @Param produto body entities.Produto true "Produto"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /produto/editar [put]
func (ph *ProdutoHandler) ProdutoEditar(c *gin.Context) {
	var produto entities.Produto

	err := c.ShouldBindJSON(&produto)

	if err != nil {
		fmt.Println("Entrando no primeiro erro")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	prd, err := ph.ProdutoEditarUseCase.Run(c, produto.ID, produto.Nome, string(produto.Categoria), produto.Descricao, produto.Preco)
	if err != nil {
		fmt.Println("Entrando no segundo erro")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	fmt.Println("Produto editado com sucesso:", prd.Nome, prd.Descricao, prd.Preco, prd.Categoria)
	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Produto editado com sucesso: " + prd.Nome,
	})
}

// RemoverProduto godoc
// @Summary Remove um produto
// @Description Remove um produto
// @Tags produto
// @Router /produto/delete/{id} [DELETE]
// @Accept  json
// @Produce  json
// @Param id path int true "ID do produto"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
func (ph *ProdutoHandler) ProdutoRemover(c *gin.Context) {
	id := c.Param("id")

	idnt, _ := strconv.Atoi(id)
	err := ph.ProdutoRemoverUseCase.Run(c, idnt)
	if err != nil {
		if strings.Contains(err.Error(), "produto não encontrado") {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Message: err.Error()})
			return
		}

		// Outros erros (erro no banco, etc)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Produto removido com sucesso",
	})
}

// ProdutoListarPorCategoria godoc
// @Summary Lista os produtos por categoria
// @Description Lista todos os produtos por categoria
// @Tags produto
// @Router /produtos/{categoria} [GET]
// @Accept  json
// @Produce  json
// @Param categoria path string true "Categoria de produtos"
// @Success 200 {object} []presenters.ProdutoDTO
// @Failure 400 {object} response.ErrorResponse
func (ph *ProdutoHandler) ProdutoListarPorCategoria(c *gin.Context) {
	categoria := c.Param("categoria")

	produtos, err := ph.ProdutoListarPorCategoriaUseCase.Run(c, categoria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Convert to DTOs
	var produtosDTO []*presenters.ProdutoDTO
	for _, produto := range produtos {
		produtosDTO = append(produtosDTO, presenters.NewProdutoDTO(produto))
	}

	c.JSON(http.StatusOK, produtosDTO)
}
