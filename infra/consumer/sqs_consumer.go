package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"lanchonete/usecases"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConsumer struct {
	client *sqs.Client
}

func NewSQSConsumer() (*SQSConsumer, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return &SQSConsumer{
		client: sqs.NewFromConfig(cfg),
	}, nil
}

func (c *SQSConsumer) StartConsumingPagamento(queueURL string, useCase usecases.PedidoAtualizarStatusPagamentoUseCase) {
	go c.consume(queueURL, func(msgBody []byte) {
		var envelope struct {
			EventType string `json:"event_type"`
			Data      struct {
				IDPagamento int     `json:"id_pagamento"`
				IDPedido    string  `json:"id_pedido"` // vem como string no JSON
				Valor       float64 `json:"valor"`
				Status      string  `json:"status"`
				DataCriacao string  `json:"data_criacao"`
			} `json:"data"`
		}

		if err := json.Unmarshal(msgBody, &envelope); err != nil {
			log.Printf("Erro ao deserializar mensagem de pagamento: %v", err)
			return
		}

		// Converte id_pedido string → int
		var pedidoID int
		if _, err := fmt.Sscanf(envelope.Data.IDPedido, "%d", &pedidoID); err != nil {
			log.Printf("Erro ao converter id_pedido: %v", err)
			return
		}

		log.Printf("📥 Evento '%s' recebido: pedidoID=%d status=%s", envelope.EventType, pedidoID, envelope.Data.Status)

		// Executa o use-case
		err := useCase.Run(context.Background(), pedidoID, envelope.Data.Status)
		if err != nil {
			log.Printf("❌ Erro ao atualizar status do pagamento: %v", err)
		} else {
			log.Printf("✅ Pagamento atualizado com sucesso: Pedido %d → %s", pedidoID, envelope.Data.Status)
		}
	})
}

func (c *SQSConsumer) consume(queueURL string, handler func([]byte)) {
	for {
		output, err := c.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			log.Printf("Erro ao ler da fila %s: %v\n", queueURL, err)
			continue
		}

		for _, msg := range output.Messages {
			handler([]byte(*msg.Body))

			// Deleta da fila
			_, err := c.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Printf("Erro ao deletar mensagem: %v\n", err)
			}
		}
	}
}
