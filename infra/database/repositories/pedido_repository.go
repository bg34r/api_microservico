package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"lanchonete/internal/domain/entities"
	"lanchonete/internal/domain/repository"
)

type pedidoMysqlRepository struct {
	db *sql.DB
}

func NewPedidoMysqlRepository(db *sql.DB) repository.PedidoRepository {
	return &pedidoMysqlRepository{db: db}
}

func (pr *pedidoMysqlRepository) CriarPedido(c context.Context, pedido *entities.Pedido) error {
	tx, err := pr.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	query := `INSERT INTO Pedido (clienteNome, totalPedido, tempoEstimado, status, statusPagamento, personalizacao) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := tx.ExecContext(c, query,
		pedido.ClienteNome,
		pedido.Total,
		"00:15:00",
		pedido.Status,
		pedido.StatusPagamento,
		pedido.Personalizacao,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir pedido: %w", err)
	}

	pedidoID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao obter ID do pedido: %w", err)
	}
	pedido.ID = int(pedidoID)

	// Inserir produtos relacionados
	prodQuery := `INSERT INTO Pedido_Produto (idPedido, idProduto, quantidade) VALUES (?, ?, ?)`
	for _, prod := range pedido.Produtos {
		_, err := tx.ExecContext(c, prodQuery, pedidoID, prod.ID, 1)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao inserir produto no pedido: %w", err)
		}
	}

	return tx.Commit()
}

func (pr *pedidoMysqlRepository) BuscarPedido(c context.Context, identificacao int) (*entities.Pedido, error) {
	query := `SELECT idPedido, clienteNome, totalPedido, tempoEstimado, status, statusPagamento, personalizacao FROM Pedido WHERE idPedido = ?`

	var pedido entities.Pedido
	var clienteNome string
	var tempoEstimado string
	var personalizacao *string

	fmt.Println("TimeStamp: ", pedido.TimeStamp, "ID: ", identificacao)
	err := pr.db.QueryRowContext(c, query, identificacao).Scan(
		&pedido.ID,
		&clienteNome,
		&pedido.Total,
		&tempoEstimado,
		&pedido.Status,
		&pedido.StatusPagamento,
		&personalizacao,
	)
	fmt.Println("Repository pedido: ", pedido.TimeStamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pedido não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar pedido: %w", err)
	}

	pedido.Produtos = []entities.Produto{}
	pedido.ClienteNome = clienteNome
	pedido.Personalizacao = personalizacao

	// Buscar produtos
	prodQuery := `SELECT p.idProduto, p.nomeProduto, p.descricaoProduto, p.precoProduto, p.categoriaProduto FROM Produto p JOIN Pedido_Produto pp ON pp.idProduto = p.idProduto WHERE pp.idPedido = ?`

	rows, err := pr.db.QueryContext(c, prodQuery, identificacao)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos do pedido: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p entities.Produto
		if err := rows.Scan(&p.ID, &p.Nome, &p.Descricao, &p.Preco, &p.Categoria); err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %v", err)
		}
		pedido.Produtos = append(pedido.Produtos, p)
	}

	return &pedido, nil
}

func (pr *pedidoMysqlRepository) AtualizarStatusPedido(c context.Context, identificacao int, status string, ultimaAtualizacao time.Time) error {
	query := `UPDATE Pedido SET status = ?, ultimaAtualizacao = ? WHERE idPedido = ?`
	result, err := pr.db.ExecContext(c, query, status, ultimaAtualizacao, identificacao)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do pedido: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualização: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("pedido não encontrado ou status não foi alterado")
	}

	return nil
}

func (pr *pedidoMysqlRepository) AtualizarStatusPagamento(c context.Context, identificacao int, statusPagamento string, ultimaAtualizacao time.Time) error {
	query := `UPDATE Pedido SET statusPagamento = ?, ultimaAtualizacao = ? WHERE idPedido = ?`
	result, err := pr.db.ExecContext(c, query, statusPagamento, ultimaAtualizacao, identificacao)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status de pagamento do pedido: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualização: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("pedido não encontrado ou status de pagamento não foi alterado")
	}

	return nil
}

func (pr *pedidoMysqlRepository) ListarTodosOsPedidos(c context.Context) ([]*entities.Pedido, error) {
	query := `SELECT idPedido, clienteNome, totalPedido, tempoEstimado, status, statusPagamento, personalizacao FROM Pedido`

	rows, err := pr.db.QueryContext(c, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos: %w", err)
	}
	defer rows.Close()

	var pedidos []*entities.Pedido
	for rows.Next() {
		var p entities.Pedido
		var clienteNome string
		var tempoEstimadoStr string
		var personalizacao *string

		if err := rows.Scan(
			&p.ID,
			&clienteNome,
			&p.Total,
			&tempoEstimadoStr,
			&p.Status,
			&p.StatusPagamento,
			&personalizacao,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear pedido: %w", err)
		}

		// pega a data atual (ano, mês e dia)
		now := time.Now()
		dataHoje := now.Format("2006-01-02") // ex: "2025-05-28"

		// concatena uma data fixa com o horário recebido
		datetimeStr := dataHoje + " " + tempoEstimadoStr // data fixa arbitrária

		// parse para time.Time completo
		_, err := time.Parse("2006-01-02 15:04:05", datetimeStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter tempoEstimado: %w", err)
		}

		p.TimeStamp = "00:15:00" // Definindo um valor fixo para o TimeStamp
		p.ClienteNome = clienteNome
		p.Personalizacao = personalizacao
		p.Produtos = []entities.Produto{}

		// Buscar ids dos produtos do pedido
		prodIDsQuery := `SELECT idProduto FROM Pedido_Produto WHERE idPedido = ?`
		prodIDRows, err := pr.db.QueryContext(c, prodIDsQuery, p.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar ids dos produtos: %w", err)
		}

		var produtos []entities.Produto

		for prodIDRows.Next() {
			var idProduto int
			if err := prodIDRows.Scan(&idProduto); err != nil {
				prodIDRows.Close()
				return nil, fmt.Errorf("erro ao escanear idProduto: %w", err)
			}

			// Busca dados completos do produto
			produtoQuery := `SELECT idProduto, nomeProduto, descricaoProduto, precoProduto, categoriaProduto FROM Produto WHERE idProduto = ?`
			var produto entities.Produto
			err = pr.db.QueryRowContext(c, produtoQuery, idProduto).Scan(
				&produto.ID,
				&produto.Nome,
				&produto.Descricao,
				&produto.Preco,
				&produto.Categoria,
			)
			if err != nil {
				prodIDRows.Close()
				return nil, fmt.Errorf("erro ao buscar produto para id %d: %w", idProduto, err)
			}

			produtos = append(produtos, produto)
		}
		prodIDRows.Close()

		if err := prodIDRows.Err(); err != nil {
			return nil, fmt.Errorf("erro na iteração dos ids dos produtos: %w", err)
		}

		// Atribui a lista de produtos ao pedido
		p.Produtos = produtos

		pedidos = append(pedidos, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração dos pedidos: %w", err)
	}

	return pedidos, nil
}
