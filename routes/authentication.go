package routes

import (
	"nihon-vocabulary/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthenticationRoutes(app fiber.Router) {
	authentication := app.Group("/auth")

	authentication.Post("/login", handlers.Login)
	authentication.Post("/register", handlers.Register)
	authentication.Post("/reset-password", handlers.ResetPassword)
}
