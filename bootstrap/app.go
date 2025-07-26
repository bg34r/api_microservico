package bootstrap

import (
	"context"
	"database/sql"
	"log"

	"lanchonete/infra/database"
	"lanchonete/internal/domain/repository"
)

type App struct {
	Env               *Env
	DB                *sql.DB
	PedidoRepository  repository.PedidoRepository
	ProdutoRepository repository.ProdutoRepository
}

func NewApp(ctx context.Context) (*App, error) {
	// Load environment variables
	env := NewEnv()

	db, err := database.NewMySQLConnection(
		env.DBUser,
		env.DBPass,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	if err != nil {
		log.Fatalf("erro ao conectar ao MySQL: %v", err)
	}

	// Initialize repositories
	_, pedidoRepo, produtoRepo, _, _ := NewRepositories(db)

	return &App{
		Env:               env,
		DB:                db,
		PedidoRepository:  pedidoRepo,
		ProdutoRepository: produtoRepo,
	}, nil
}
