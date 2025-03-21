package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App) {
	api := app.Group("/api", func(c *fiber.Ctx) error { return nil })

	images := api.Group("/images")

	images.Post("/uploadImages", func(c *fiber.Ctx) error {
		log.Println("hello")
		return nil
	})

}
