-- Microserviço de Produtos e Pedidos

DROP TABLE IF EXISTS `Produto`;
CREATE TABLE `Produto` (
  `idProduto` int NOT NULL AUTO_INCREMENT,
  `nomeProduto` varchar(45) NOT NULL,
  `descricaoProduto` varchar(125) NOT NULL,
  `precoProduto` float NOT NULL,
  `categoriaProduto` enum('Lanche','Bebida','Acompanhamento','Sobremesa') DEFAULT NULL,
  PRIMARY KEY (`idProduto`),
  UNIQUE KEY `idProduto_UNIQUE` (`idProduto`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `Produto` VALUES (1,'X-Salada','Lanche com tomate, alface, hambúrguer e maionese',22.5,'Lanche'),(2,'Coca-cola','Refrigerante gelado de cola',6,'Bebida'),(3,'Batata-frita','Porção de batata-frita palito crocante',18,'Acompanhamento'),(4,'Mousse de chocolate','Chocolate cremoso ao leite',12.5,'Sobremesa'),(5,'X-Frango','Lanche com frango desfiado e bacon',26,'Lanche'),(6,'X-Tudo','Calabresa, Bacon, 2 ovos, maionese e queijo',28.5,'Lanche'),(7,'Cachorro-quente','2 salsichas, purê, milho  ervilha',18,'Lanche'),(8,'Cachorrão especial','2 salsichas, calabresa, bacon, purê, milho e ervilha',22,'Lanche'),(9,'Fanta','Refrigerante sabor laranja gelado',5.5,'Bebida'),(10,'Sprite','Refrigerante sabor limão gelado',5.5,'Bebida');

DROP TABLE IF EXISTS `Pedido`;
CREATE TABLE `Pedido` (
  `idPedido` INT NOT NULL AUTO_INCREMENT,
  `clienteNome` VARCHAR(100) DEFAULT 'Cliente',
  `totalPedido` FLOAT NOT NULL DEFAULT 0,
  `tempoEstimado` TIME NOT NULL DEFAULT '00:15:00',
  `ultimaAtualizacao` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status` VARCHAR(50) DEFAULT 'Pendente',
  `statusPagamento` VARCHAR(50) DEFAULT 'Pendente',
  `personalizacao` VARCHAR(255) DEFAULT NULL,
  PRIMARY KEY (`idPedido`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `Pedido_Produto`;
CREATE TABLE `Pedido_Produto` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `idPedido` INT NOT NULL,
  `idProduto` INT NOT NULL,
  `quantidade` INT DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_pedido` (`idPedido`),
  KEY `idx_produto` (`idProduto`),
  CONSTRAINT `fk_pedido` FOREIGN KEY (`idPedido`) REFERENCES `Pedido` (`idPedido`) ON DELETE CASCADE,
  CONSTRAINT `fk_produto` FOREIGN KEY (`idProduto`) REFERENCES `Produto` (`idProduto`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `Pagamento`;
CREATE TABLE `Pagamento` (
  `idPagamento` int NOT NULL AUTO_INCREMENT,
  `dataCriacao` datetime NOT NULL,
  `Status` enum('Pendente','Recebido','Em Preparação','Pronto','Finalizado') NOT NULL DEFAULT 'Pendente',
  `idPedido` int NOT NULL,
  PRIMARY KEY (`idPagamento`),
  UNIQUE KEY `idPagamento_UNIQUE` (`idPagamento`),
  KEY `idPedido_idx` (`idPedido`),
  CONSTRAINT `idPedido` FOREIGN KEY (`idPedido`) REFERENCES `Pedido` (`idPedido`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `cliente`;
CREATE TABLE `cliente` (
  `id` int NOT NULL AUTO_INCREMENT,
  `cpf` varchar(11) NOT NULL,
  `nome` varchar(100) NOT NULL,
  `email` varchar(100) DEFAULT NULL,
  `telefone` varchar(20) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cpf` (`cpf`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `cliente` VALUES (1,'12345678901','Test User 1','test1@example.com','11999999999','2025-04-30 12:18:06','2025-04-30 12:18:06'),(2,'98765432101','Test User 2','test2@example.com','11988888888','2025-04-30 12:18:06','2025-04-30 12:18:06');

DROP TABLE IF EXISTS `Acompanhamento`;
CREATE TABLE `Acompanhamento` (
  `idAcompanhamento` INT AUTO_INCREMENT PRIMARY KEY,
  `tempoEstimado` TIME NOT NULL DEFAULT '00:15:00',
  `ultimaAtualizacao` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `FilaPedidos`;
CREATE TABLE `FilaPedidos` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `idAcompanhamento` INT NOT NULL,
  `idPedido` INT NOT NULL,
  `ordem` INT NOT NULL,
  FOREIGN KEY (`idAcompanhamento`) REFERENCES `Acompanhamento` (`idAcompanhamento`) ON DELETE CASCADE,
  FOREIGN KEY (`idPedido`) REFERENCES `Pedido` (`idPedido`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
