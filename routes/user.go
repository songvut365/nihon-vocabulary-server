package routes

import (
	"nihon-vocabulary/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router) {
	user := app.Group("/user")

	user.Get("/", handlers.GetUser)
	user.Put("/", handlers.UpdateUser)
	user.Put("/password", handlers.ChangePasswordUser)
	user.Delete("/", handlers.DeleteUser)
}
