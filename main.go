package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"lanchonete/bootstrap"
	_ "lanchonete/docs"
	queue "lanchonete/infra/consumer"
	"lanchonete/internal/interfaces/http/server"
	"lanchonete/usecases"
)

// @title Lanchonete API - Tech Challenge 2
// @version 1.0
// @description API para o Tech Challenge 2 da FIAP - SOAT

// @host localhost:8080
// @BasePath /
//
//go:generate go run github.com/swaggo/swag/cmd/swag@latest init
func main() {
	fmt.Println("🔧 Iniciando aplicação...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize application
	app, err := bootstrap.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create and configure HTTP server
	srv := server.NewServer(app)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Inicia consumer da fila SQS de pagamento
	sqsConsumer, err := queue.NewSQSConsumer()
	if err != nil {
		log.Fatalf("Erro ao inicializar consumidor SQS: %v", err)
	}

	// Cria use-case e inicia consumo
	pagamentoUseCase := usecases.NewPedidoAtualizarStatusPagamentoUseCase(app.PedidoRepository)
	sqsConsumer.StartConsumingPagamento(app.Env.PagamentoQueueURL, pagamentoUseCase)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	// Perform any cleanup here if needed
}
