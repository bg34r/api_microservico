package repository

// import (
// 	"lanchonete/internal/domain/repository"
// 	 "lanchonete/infra/database/mongo"
// 	 "lanchonete/infra/database/mongo/repositories"
// )

// // RepositorioFactory é responsável por criar instâncias de repositórios
// type RepositorioFactory struct {
// 	db mongo.Database
// }

// // NovoRepositorioFactory cria uma nova instância de RepositorioFactory
// func NovoRepositorioFactory(db mongo.Database) *RepositorioFactory {
// 	return &RepositorioFactory{
// 		db: db,
// 	}
// }

// // CriarPedidoRepository cria uma instância de PedidoRepository
// func (f *RepositorioFactory) CriarPedidoRepository() repository.PedidoRepository {
// 	return repositories.NewPedidoMongoRepository(f.db, mongo.CollectionPedido)
// }

// // CriarProdutoRepository cria uma instância de ProdutoRepository
// func (f *RepositorioFactory) CriarProdutoRepository() repository.ProdutoRepository {
// 	return repositories.NewProdutoMongoRepository(f.db, mongo.CollectionProduto)
// }

// // CriarClienteRepository cria uma instância de ClienteRepository
// func (f *RepositorioFactory) CriarClienteRepository() repository.ClienteRepository {
// 	return repositories.NewClienteMongoRepository(f.db, mongo.CollectionCliente)
// }

// // CriarPagamentoRepository cria uma instância de PagamentoRepository
// func (f *RepositorioFactory) CriarPagamentoRepository() repository.PagamentoRepository {
// 	return repositories.NewPagamentoMongoRepository(f.db, mongo.CollectionPagamento)
// }

// // CriarAcompanhamentoRepository cria uma instância de AcompanhamentoRepository
// func (f *RepositorioFactory) CriarAcompanhamentoRepository() repository.AcompanhamentoRepository {
// 	return repositories.NewAcompanhamentoMongoRepository(f.db, mongo.CollectionAcompanhamento)
// }
