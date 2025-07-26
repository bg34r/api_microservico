package route

import (
	"database/sql"
	"lanchonete/bootstrap"
	repo "lanchonete/infra/database/repositories"
	handler "lanchonete/internal/interfaces/http/handlers"
	"lanchonete/usecases"

	"github.com/gin-gonic/gin"
)

// NovoPedidoRouter cria um novo roteador para pedidos
func NovoPedidoRouter(env *bootstrap.Env, db *sql.DB, router *gin.RouterGroup) {

	// Criar casos de uso
	pedidoIncluirUseCase := usecases.NewPedidoIncluirUseCase(repo.NewPedidoMysqlRepository(db))
	pedidoBuscarPorIdUseCase := usecases.NewPedidoBuscarPorIdUseCase(repo.NewPedidoMysqlRepository(db))
	pedidoAtualizarStatusUseCase := usecases.NewPedidoAtualizarStatusUseCase(repo.NewPedidoMysqlRepository(db))
	produtoBuscarPorIdUseCase := usecases.NewProdutoBuscaPorIdUseCase(repo.NewProdutoMysqlRepository(db))
	pedidoListarTodosUseCase := usecases.NewPedidoListarTodosUseCase(repo.NewPedidoMysqlRepository(db))

	// Criar handler
	puc := &handler.PedidoHandler{
		PedidoIncluirUseCase:         pedidoIncluirUseCase,
		PedidoBuscarPorIdUseCase:     pedidoBuscarPorIdUseCase,
		PedidoAtualizarStatusUseCase: pedidoAtualizarStatusUseCase,
		ProdutoBuscarPorIdUseCase:    produtoBuscarPorIdUseCase,
		PedidoListarTodosUseCase:     pedidoListarTodosUseCase,
	}

	// Configurar rotas
	router.POST("/pedidos", puc.CriarPedido)
	router.GET("/pedidos/:nroPedido", puc.BuscarPedido)
	router.PUT("/pedidos/:nroPedido/status/:status", puc.AtualizarStatusPedido)
	router.GET("/pedidos/listartodos", puc.ListarTodosOsPedidos)
}
