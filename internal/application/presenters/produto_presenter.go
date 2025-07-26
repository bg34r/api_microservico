package presenters

import "lanchonete/internal/domain/entities"

type ProdutoDTO struct {
	Identificacao int                 `json:"identificacao"`
	Nome          string              `json:"nome"`
	Categoria     entities.CatProduto `json:"categoria"`
	Descricao     string              `json:"descricao"`
	Preco         float32             `json:"preco"`
}

func NewProdutoDTO(produto *entities.Produto) *ProdutoDTO {
	if produto == nil {
		return nil
	}

	return &ProdutoDTO{
		Identificacao: produto.ID,
		Nome:          produto.Nome,
		Categoria:     produto.Categoria,
		Descricao:     produto.Descricao,
		Preco:         produto.Preco,
	}
}
