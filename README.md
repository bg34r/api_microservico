# 🍔 Microserviço de Produtos e Pedidos - Lanchonete

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Test Coverage](https://img.shields.io/badge/Coverage-90.1%25-brightgreen.svg)](./usecases)
[![Docker](https://img.shields.io/badge/Docker-Compose-blue.svg)](docker-compose.yml)

## 📖 Visão Geral

Este microserviço foi desenvolvido seguindo princípios de **Clean Architecture** e focado especificamente no gerenciamento de **produtos** e **pedidos** de uma lanchonete. O sistema possui **90.1% de cobertura de testes** e está totalmente containerizado.

---

## Estrutura de Diretórios

- **bootstrap/**  
  Inicialização da aplicação, configuração de ambiente e injeção de dependências.

- **db/**  
  Scripts de inicialização do banco de dados (ex: `init.sql`).

- **docs/**  
  Documentação Swagger gerada automaticamente.

- **infra/**  
  Implementações de infraestrutura, como repositórios de banco de dados.

- **internal/**  
  Código de domínio e regras de negócio:
  - **application/**: Casos de uso e presenters.
  - **domain/**: Entidades e interfaces de repositórios.
  - **infrastructure/**: Implementações específicas de infraestrutura.
  - **interfaces/**: Camada de entrada (HTTP, handlers, rotas).

- **usecases/**  
  Casos de uso específicos.

- **main.go**  
  Ponto de entrada da aplicação.

---

## Esquema do Banco de Dados

Exemplo simplificado das principais tabelas:

```sql
CREATE TABLE Cliente (
    cpfCliente VARCHAR(11) PRIMARY KEY,
    nomeCliente VARCHAR(100),
    emailCliente VARCHAR(100)
);

CREATE TABLE Produto (
    idProduto INT AUTO_INCREMENT PRIMARY KEY,
    nomeProduto VARCHAR(100),
    descricaoProduto TEXT,
    precoProduto FLOAT,
    personalizacaoProduto VARCHAR(255),
    categoriaProduto VARCHAR(50)
);

CREATE TABLE Pedido (
    idPedido INT AUTO_INCREMENT PRIMARY KEY,
    cliente VARCHAR(11),
    totalPedido FLOAT,
    tempoEstimado VARCHAR(8),
    status VARCHAR(50),
    statusPagamento VARCHAR(50),
    FOREIGN KEY (cliente) REFERENCES Cliente(cpfCliente)
);

CREATE TABLE Pedido_Produto (
    idPedido INT,
    idProduto INT,
    quantidade INT,
    FOREIGN KEY (idPedido) REFERENCES Pedido(idPedido),
    FOREIGN KEY (idProduto) REFERENCES Produto(idProduto)
);

CREATE TABLE Pagamento (
    idPagamento INT AUTO_INCREMENT PRIMARY KEY,
    dataCriacao DATETIME,
    Status VARCHAR(50),
    idPedido INT,
    FOREIGN KEY (idPedido) REFERENCES Pedido(idPedido)
);

CREATE TABLE Acompanhamento (
    idAcompanhamento INT AUTO_INCREMENT PRIMARY KEY,
    tempoEstimado VARCHAR(8)
);

CREATE TABLE Acompanhamento_Pedido (
    idAcompanhamento INT,
    idPedido INT,
    FOREIGN KEY (idAcompanhamento) REFERENCES Acompanhamento(idAcompanhamento),
    FOREIGN KEY (idPedido) REFERENCES Pedido(idPedido)
);
```
## 🚀 Como Rodar a Aplicação

### Pré-requisitos

- Go 1.21+
- Docker e Docker Compose
- MySQL (se não usar Docker)

### Passos

1. **Clone o repositório:**
   ```bash
   git clone <seu-repositorio>
   cd CleanArch-main
   ```

2. **Suba o ambiente com Docker:**
   ```bash
   docker-compose up -d
   ```

3. **Instale as dependências (desenvolvimento local):**
   ```bash
   go mod tidy
   ```

4. **Acesse os serviços:**
   - **API**: http://localhost:8080
   - **Documentação Swagger**: http://localhost:8080/docs
   - **Adminer (DB)**: http://localhost:8081

### Desenvolvimento Local (sem Docker)

1. **Configure as variáveis de ambiente:**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_USER=root
   export DB_PASS=password
   export DB_NAME=lanchonete
   ```

2. **Execute a aplicação:**
   ```bash
   go run main.go
   ```

---

## 🧪 Testes

O projeto possui **90.1% de cobertura de testes** com testes unitários organizados por domínio:

```bash
# Executar todos os testes
go test ./usecases/ -v

# Executar testes com coverage
go test ./usecases/ -cover

# Testes específicos por domínio
go test -run TestProduto ./usecases/ -v  # Produtos
go test -run TestPedido ./usecases/ -v   # Pedidos
```

### 📊 Estrutura de Testes

**Produtos** (6 use cases, 26 testes):
- `produto_buscar_por_id_test.go` (4 testes)
- `produto_editar_test.go` (4 testes)
- `produto_incluir_test.go` (4 testes)
- `produto_remover_test.go` (5 testes)
- `produto_listar_todos_test.go` (5 testes)
- `produto_listar_por_categoria_test.go` (6 testes)

**Pedidos** (5 use cases, 21 testes):
- `pedido_incluir_test.go` (3 testes)
- `pedido_buscar_por_id_test.go` (4 testes)
- `pedido_listar_todos_test.go` (3 testes)
- `pedido_atualizar_status_test.go` (5 testes)
- `pedido_atualizar_status_pagamento_test.go` (6 testes)

---

## 📈 Métricas de Qualidade

- **Cobertura de Testes**: 90.1%
- **Total de Testes**: 47 testes
- **Use Cases Cobertos**: 11/11 (100%)
- **Arquitetura**: Clean Architecture
- **Padrões**: Repository Pattern, Dependency Injection

---

## 🤝 Como Contribuir

1. Fork o projeto
2. Crie uma branch para sua feature a partir da `desenv` (`git checkout -b feature/AmazingFeature desenv`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request para a branch `desenv`

> **Nota**: O desenvolvimento principal acontece na branch `desenv`. Use ela como base para novas features.

---

## 📄 Licença

Este projeto possui fins educacionais e de demonstração de boas práticas em Go e Clean Architecture.
