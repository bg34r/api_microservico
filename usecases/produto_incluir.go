package usecases

import (
	"context"
	"fmt"
	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
	"lanchonete/internal/interfaces/publisher"
)

type ProdutoIncluirUseCase interface {
	Run(ctx context.Context, nome, categoria, descricao string, preco float32) (*entities.Produto, error)
}

type produtoIncluirUseCase struct {
	produtoRepository repository.ProdutoRepository
	eventPublisher    publisher.EventPublisher
}

func NewProdutoIncluirUseCase(produtoRepository repository.ProdutoRepository, publisher publisher.EventPublisher) ProdutoIncluirUseCase {
	return &produtoIncluirUseCase{
		produtoRepository: produtoRepository,
		eventPublisher:    publisher,
	}
}

func (pd *produtoIncluirUseCase) Run(c context.Context, nome string, categoria string, descricao string, preco float32) (*entities.Produto, error) {

	produto, err := entities.ProdutoNew(nome, categoria, descricao, preco)

	if err != nil {
		return nil, fmt.Errorf("criação de produto inválida: %w", err)
	}

	err = pd.produtoRepository.AdicionarProduto(c, produto)
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar produto: %w", err)
	}

	// ✨ Publicar evento no SQS
	payload := map[string]interface{}{
		"id_produto": produto.ID,
		"nome":       produto.Nome,
		"categoria":  produto.Categoria,
		"descricao":  produto.Descricao,
		"preco":      produto.Preco,
	}

	err = pd.eventPublisher.Publish("produto_criado", payload)
	if err != nil {
		fmt.Println("⚠️ Falha ao publicar evento do produto:", err)
	}

	return produto, nil
}
