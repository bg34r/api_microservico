package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type produtoMysqlRepository struct {
	database *sql.DB
}

func NewProdutoMysqlRepository(db *sql.DB) repository.ProdutoRepository {
	return &produtoMysqlRepository{
		database: db,
	}
}

func (pr *produtoMysqlRepository) AdicionarProduto(c context.Context, produto *entities.Produto) error {
	query := "INSERT INTO Produto (nomeProduto, descricaoProduto, precoProduto, categoriaProduto) VALUES (?, ?, ?, ?)"
	result, err := pr.database.ExecContext(c, query, produto.Nome, produto.Descricao, produto.Preco, produto.Categoria)
	if err != nil {
		return err
	}

	// Captura o ID gerado automaticamente
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Atribui o ID ao produto
	produto.ID = int(lastInsertID)
	fmt.Printf("DEBUG: Produto criado com ID: %d\n", produto.ID)

	return nil
}

func (pr *produtoMysqlRepository) BuscarProdutoPorId(c context.Context, id int) (*entities.Produto, error) {
	query := "SELECT idProduto, nomeProduto, descricaoProduto, precoProduto, categoriaProduto FROM Produto WHERE idProduto = ?"
	var produto entities.Produto
	err := pr.database.QueryRowContext(c, query, id).
		Scan(&produto.ID, &produto.Nome, &produto.Descricao, &produto.Preco, &produto.Categoria)
	fmt.Println("Repository Buscando produto:", produto.Nome)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produto não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar produto: %v", err)
	}
	fmt.Println("Repository Produto encontrado:", produto.Nome, produto.Descricao, produto.Preco, produto.Categoria)
	return &produto, nil
}

func (pr *produtoMysqlRepository) ListarTodosOsProdutos(c context.Context) ([]*entities.Produto, error) {
	query := "SELECT idProduto, nomeProduto, descricaoProduto, precoProduto, categoriaProduto FROM Produto"
	rows, err := pr.database.QueryContext(c, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %v", err)
	}
	defer rows.Close()

	var produtos []*entities.Produto
	for rows.Next() {
		var p entities.Produto
		if err := rows.Scan(&p.ID, &p.Nome, &p.Descricao, &p.Preco, &p.Categoria); err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %v", err)
		}
		produtos = append(produtos, &p)
	}
	return produtos, nil
}

func (pr *produtoMysqlRepository) EditarProduto(c context.Context, produto *entities.Produto) error {
	query := "UPDATE Produto SET nomeProduto = ?, descricaoProduto = ?, precoProduto = ?, categoriaProduto = ? WHERE nomeProduto = ?"
	fmt.Println("Repository Atualizando produto:", produto.Nome, produto.Descricao, produto.Preco, produto.Categoria)
	result, err := pr.database.ExecContext(c, query, produto.Nome, produto.Descricao, produto.Preco, produto.Categoria, produto.Nome)
	if err != nil {
		return fmt.Errorf("erro ao atualizar produto: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualização: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	return nil
}

func (pr *produtoMysqlRepository) RemoverProduto(c context.Context, id int) error {
	query := "DELETE FROM Produto WHERE idProduto = ?"
	result, err := pr.database.ExecContext(c, query, id)
	if err != nil {
		return fmt.Errorf("erro ao remover produto: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar remoção: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	return nil
}

func (pr *produtoMysqlRepository) ListarPorCategoria(c context.Context, categoria string) ([]*entities.Produto, error) {
	query := "SELECT idProduto, nomeProduto, descricaoProduto, precoProduto, categoriaProduto FROM Produto WHERE categoriaProduto = ?"
	rows, err := pr.database.QueryContext(c, query, categoria)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos por categoria: %v", err)
	}
	defer rows.Close()

	var produtos []*entities.Produto
	for rows.Next() {
		var p entities.Produto
		if err := rows.Scan(&p.ID, &p.Nome, &p.Descricao, &p.Preco, &p.Categoria); err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %v", err)
		}
		produtos = append(produtos, &p)
	}

	// Verificar se houve erro durante a iteração
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração dos produtos: %v", err)
	}

	return produtos, nil
}
