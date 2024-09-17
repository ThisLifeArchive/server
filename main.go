package main

import (
	"fmt"
	"log"

	"github.com/ThisLifeArchive/server/episodes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/episodes", func(c *fiber.Ctx) error {
		eps, err := episodes.List()
		if err != nil {
			return fmt.Errorf("failed to list episodes: %w", err)
		}
		return c.JSON(eps)
	})
	if err := app.Listen(":8888"); err != nil {
		log.Fatal(err)
	}
}
