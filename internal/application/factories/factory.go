package factory

// import (
// 	"lanchonete/internal/infrastructure/repository"
// 	"lanchonete/usecases"
// )

// // UseCaseFactory é responsável por criar instâncias de casos de uso
// type UseCaseFactory struct {
// 	repositorioFactory *repository.RepositorioFactory
// }

// // NovoUseCaseFactory cria uma nova instância de UseCaseFactory
// func NovoUseCaseFactory(repositorioFactory *repository.RepositorioFactory) *UseCaseFactory {
// 	return &UseCaseFactory{
// 		repositorioFactory: repositorioFactory,
// 	}
// }

// // CriarPedidoIncluirUseCase cria uma instância de PedidoIncluirUseCase
// func (f *UseCaseFactory) CriarPedidoIncluirUseCase() usecases.PedidoIncluirUseCase {
// 	return usecases.NewPedidoIncluirUseCase(f.repositorioFactory.CriarPedidoRepository())
// }

// // CriarPedidoBuscarPorIdUseCase cria uma instância de PedidoBuscarPorIdUseCase
// func (f *UseCaseFactory) CriarPedidoBuscarPorIdUseCase() usecases.PedidoBuscarPorIdUseCase {
// 	return usecases.NewPedidoBuscarPorIdUseCase(f.repositorioFactory.CriarPedidoRepository())
// }

// // CriarPedidoAtualizarStatusUseCase cria uma instância de PedidoAtualizarStatusUseCase
// func (f *UseCaseFactory) CriarPedidoAtualizarStatusUseCase() usecases.PedidoAtualizarStatusUseCase {
// 	return usecases.NewPedidoAtualizarStatusUseCase(f.repositorioFactory.CriarPedidoRepository())
// }

// // CriarPedidoListarTodosUseCase cria uma instância de PedidoListarTodosUseCase
// func (f *UseCaseFactory) CriarPedidoListarTodosUseCase() usecases.PedidoListarTodosUseCase {
// 	return usecases.NewPedidoListarTodosUseCase(f.repositorioFactory.CriarPedidoRepository())
// }

// // CriarProdutoBuscarPorIdUseCase cria uma instância de ProdutoBuscaPorIdUseCase
// func (f *UseCaseFactory) CriarProdutoBuscarPorIdUseCase() usecases.ProdutoBuscaPorIdUseCase {
// 	return usecases.NewProdutoBuscaPorIdUseCase(f.repositorioFactory.CriarProdutoRepository())
// }
