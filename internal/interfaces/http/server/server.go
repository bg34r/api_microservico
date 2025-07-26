package server

import (
	"fmt"
	"sync"

	"lanchonete/bootstrap"
	handler "lanchonete/internal/interfaces/http/handlers"
	"lanchonete/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	app       *bootstrap.App
	router    *gin.Engine
	setupOnce sync.Once
}

func NewServer(app *bootstrap.App) *Server {
	router := gin.Default()

	fmt.Printf("🆕 Instância de Server criada: %p\n", router)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return &Server{
		app:    app,
		router: router,
	}
}

func (s *Server) SetupRoutes() {
	s.setupOnce.Do(func() {
		fmt.Println("🔁 Registrando rotas HTTP...")

		api := s.router.Group("")

		// Produto
		produtoRepo := s.app.ProdutoRepository
		produtoIncluir := usecases.NewProdutoIncluirUseCase(produtoRepo)
		produtoEditar := usecases.NewProdutoEditarUseCase(produtoRepo)
		produtoRemover := usecases.NewProdutoRemoverUseCase(produtoRepo)
		produtoBuscar := usecases.NewProdutoBuscaPorIdUseCase(produtoRepo)
		produtoListarTodos := usecases.NewProdutoListarTodosUseCase(produtoRepo)
		produtoListarPorCategoria := usecases.NewProdutoListarPorCategoriaUseCase(produtoRepo)

		produtoHandler := handler.NewProdutoHandler(
			produtoIncluir,
			produtoBuscar,
			produtoListarTodos,
			produtoEditar,
			produtoRemover,
			produtoListarPorCategoria,
		)
		api.POST("/produto", produtoHandler.ProdutoIncluir)
		api.GET("/produto/:id", produtoHandler.ProdutoBuscarPorId)
		api.GET("/produtos", produtoHandler.ProdutoListarTodos)
		api.GET("/produtos/:categoria", produtoHandler.ProdutoListarPorCategoria)
		api.PUT("/produto/editar", produtoHandler.ProdutoEditar)
		api.DELETE("/produto/delete/:id", produtoHandler.ProdutoRemover)

		// Pedido
		pedidoRepo := s.app.PedidoRepository
		pedidoIncluir := usecases.NewPedidoIncluirUseCase(pedidoRepo)
		pedidoBuscar := usecases.NewPedidoBuscarPorIdUseCase(pedidoRepo)
		pedidoAtualizar := usecases.NewPedidoAtualizarStatusUseCase(pedidoRepo)
		pedidoAtualizarPagamento := usecases.NewPedidoAtualizarStatusPagamentoUseCase(pedidoRepo)
		pedidoListarTodos := usecases.NewPedidoListarTodosUseCase(pedidoRepo)
		produtoBuscaPorId := usecases.NewProdutoBuscaPorIdUseCase(produtoRepo)

		pedidoHandler := handler.NewPedidoHandler(
			pedidoIncluir,
			pedidoBuscar,
			pedidoAtualizar,
			pedidoAtualizarPagamento,
			produtoBuscaPorId,
			pedidoListarTodos,
		)
		api.POST("/pedidos", pedidoHandler.CriarPedido)
		api.GET("/pedidos/:nroPedido", pedidoHandler.BuscarPedido)
		api.PUT("/pedidos/:nroPedido/status/:status", pedidoHandler.AtualizarStatusPedido)
		api.PUT("/pedidos/:nroPedido/pagamento/:statusPagamento", pedidoHandler.AtualizarStatusPagamento)
		api.GET("/pedidos/listartodos", pedidoHandler.ListarTodosOsPedidos)

		// Health check e Swagger
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	})
}

func (s *Server) Start() error {
	s.SetupRoutes()
	return s.router.Run(s.app.Env.ServerAddress)
}
