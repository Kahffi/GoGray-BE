package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	imagestore "github.com/Kahffi/GoGray-BE/internal/image_store"
	"github.com/Kahffi/GoGray-BE/models"
)

type TransactionService interface {
	CreateTransaction(request models.TransactionCreateRequest) error
}

type transactionService struct {
	imageStore *imagestore.ImageStore
	publisher  *sarama.AsyncProducer
}

func NewTransactionService(imageStore *imagestore.ImageStore, publisher *sarama.AsyncProducer) TransactionService {
	return &transactionService{
		imageStore: imageStore,
		publisher:  publisher,
	}
}

// send the response immediately, and run upload image in the background, uploaded image then publish it to kafka
func (ts *transactionService) CreateTransaction(request models.TransactionCreateRequest) error {
	const workerLimit = 5
	ctxWTO, cancel := context.WithTimeout(context.Background(), time.Second*10)

	imagesChan := make(chan string, len(request.Images))
	imageURlChan := make(chan string, len(request.Images))

	wg := sync.WaitGroup{}
	wg.Add(workerLimit)

	for range workerLimit {
		go func() {
			defer wg.Done()
			ts.uploadImageToStore(ctxWTO, imagesChan, imageURlChan)
		}()
	}

	// populate imagesChan channel
	for _, image := range request.Images {
		imagesChan <- image
	}
	// closing imagesChan so worker doesn't wait for new data
	close(imagesChan)

	// to cancel the context after upload workers done
	go func() {
		wg.Wait()
		close(imageURlChan)

		for im := range imageURlChan {
			msg := &sarama.ProducerMessage{
				Topic: "image_url_request",
				Value: sarama.StringEncoder(im),
			}
			log.Println("Publishing to Kafka")
			(*ts.publisher).Input() <- msg
		}

		log.Println("Routine Finished")
		cancel()
	}()

	return nil
}

func (ts *transactionService) uploadImageToStore(ctx context.Context, imagesChan <-chan string, imageURLChan chan<- string) {
	for image := range imagesChan {
		uploadedURL, err := (*ts.imageStore).UploadImage(ctx, string(image))
		if err != nil {
			log.Printf("Failed to Upload image: %v\n", err)
			continue
		}
		log.Println("Upload image success:", uploadedURL)
		imageURLChan <- uploadedURL
	}
	log.Println("Worker Stopped")
}
