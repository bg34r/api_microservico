package entities

import (
	"errors"
	"fmt"
	"strings"
)

type CatProduto string

const (
	Lanche         CatProduto = "Lanche"
	Acompanhamento CatProduto = "Acompanhamento"
	Bebida         CatProduto = "Bebida"
	Sobremesa      CatProduto = "Sobremesa"
)

type Produto struct {
	ID        int        `json:"idProduto"`
	Nome      string     `json:"nomeProduto"`
	Categoria CatProduto `json:"categoriaProduto"`
	Descricao string     `json:"descricaoProduto"`
	Preco     float32    `json:"precoProduto"`
}

func ProdutoNew(nome string, categoria string, descricao string, preco float32) (*Produto, error) {
	fmt.Println("Nome:", nome, "\nCategoria: ", categoria, "\nDescrição:", descricao, "\nPreço:", preco)
	if strings.TrimSpace(nome) == "" || preco <= 0 || strings.TrimSpace(categoria) == "" {
		return nil, errors.New("todos os campos são obrigatórios e o preço maior que zero")
	}

	var cat_prod CatProduto

	switch CatProduto(categoria) {
	case Lanche, Acompanhamento, Bebida, Sobremesa:
		cat_prod = CatProduto(categoria)
	default:
		return nil, errors.New("categoria inválida")
	}

	return &Produto{
		Nome:      nome,
		Categoria: cat_prod,
		Descricao: descricao,
		Preco:     preco,
	}, nil
}
