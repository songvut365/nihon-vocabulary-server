package routes

import (
	"nihon-vocabulary/handlers"

	"github.com/gofiber/fiber/v2"
)

func VocabularyRoutes(app fiber.Router) {
	vocabulary := app.Group("/vocabulary")

	vocabulary.Post("/", handlers.CreateVocabulary)
	vocabulary.Get("/", handlers.GetVocabularies)
	vocabulary.Get("/:id", handlers.GetVocabulary)
	vocabulary.Put("/:id", handlers.UpdateVocabulary)
	vocabulary.Delete("/:id", handlers.DeleteVocabulary)
}
