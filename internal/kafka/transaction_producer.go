package kafka

import (
	"github.com/IBM/sarama"
)

type TransactionProducer interface {
	Send(message []byte)
	BulkSend(messages [][]byte) error
}

type transactionProducer struct {
	SyncProducer *sarama.SyncProducer
}

func NewTransactionProducer(syncProducer *sarama.SyncProducer) TransactionProducer {
	return &transactionProducer{
		SyncProducer: syncProducer,
	}
}

func (tp *transactionProducer) Send(message []byte) {
	(*tp.SyncProducer).SendMessage(&sarama.ProducerMessage{
		Topic:     "image_process_request",
		Value:     sarama.StringEncoder(string(message)),
		Partition: 1,
	})
}

func (tp *transactionProducer) BulkSend(messages [][]byte) error {
	// initialize messages
	producerMessages := make([]*sarama.ProducerMessage, len(messages))

	// populate messages
	for idx, message := range messages {
		producerMessages[idx] = &sarama.ProducerMessage{
			Topic:     "image_process_request",
			Value:     sarama.StringEncoder(string(message)),
			Partition: 1,
		}
	}
	return (*tp.SyncProducer).SendMessages(producerMessages)
}
