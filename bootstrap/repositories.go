package bootstrap

import (
	"database/sql"
	"lanchonete/infra/database/repositories"
	"lanchonete/internal/domain/repository"
)

func NewRepositories(db *sql.DB) (
	acomp interface{}, // placeholder para manter compatibilidade
	pedido repository.PedidoRepository,
	produto repository.ProdutoRepository,
	cliente interface{}, // placeholder para manter compatibilidade
	pagamento interface{}, // placeholder para manter compatibilidade
) {
	// Retornando nil para repositórios não utilizados neste microserviço
	acomp = nil
	pedido = repositories.NewPedidoMysqlRepository(db)
	produto = repositories.NewProdutoMysqlRepository(db)
	cliente = nil
	pagamento = nil

	return
}
