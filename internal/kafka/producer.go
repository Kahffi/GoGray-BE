package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func NewSyncProducer(cfg *sarama.Config, serverAdress []string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(serverAdress, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize producer: %v", err)
		return nil, err
	}
	log.Printf("Kafka Producer Initialized")
	return producer, nil
}

func NewAsyncProducer(cfg *sarama.Config, serverAdress []string) (sarama.AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(serverAdress, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize producer: %v", err)
		return nil, err
	}
	log.Printf("Kafka Producer Initialized")

	go func() {
		for err := range producer.Errors() {
			log.Println("Kafka publish error: ", err)
		}

	}()

	go func() {
		for suc := range producer.Successes() {
			log.Println("Publishing Success: ", suc)
		}
	}()

	return producer, nil
}
