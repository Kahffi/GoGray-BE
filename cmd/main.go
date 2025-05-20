package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"

	"github.com/Kahffi/GoGray-BE/internal/http/controller"
	imagestore "github.com/Kahffi/GoGray-BE/internal/image_store"
	"github.com/Kahffi/GoGray-BE/internal/kafka"
	"github.com/Kahffi/GoGray-BE/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Cannot Load .env file")
	}

	imgurClientID := os.Getenv("IMGUR_CLIENT_ID")
	imgurRefreshToken := os.Getenv("IMGUR_REFRESH_TOKEN")
	imgurClientSecret := os.Getenv("IMGUR_CLIENT_SECRET")

	imgur, err := imagestore.NewImgurService(
		&imagestore.ImgurOpts{
			ClientID:     imgurClientID,
			RefreshToken: imgurRefreshToken,
			ClientSecret: imgurClientSecret,
			HttpClient:   &http.Client{},
		})

	if err != nil {	
		panic(fmt.Sprintf("Cannot grant access token for imagestore: %v", err))
	}
	server := echo.New()

	config := sarama.NewConfig()
	config.Producer.Return.Errors = true // Capture errors on the Errors() channel.
	config.Producer.Return.Successes = false

	asyncProducer, err := kafka.NewAsyncProducer(config, []string{"localhost:9092"})

	if err != nil {
		panic("Cannot initialize kafka")
	}

	transactionService := service.NewTransactionService(&imgur, &asyncProducer)

	transactionController := controller.NewTransactionController(transactionService)

	server.POST("/images", transactionController.Create)

	server.Logger.Fatal(server.Start(":1323"))
}
