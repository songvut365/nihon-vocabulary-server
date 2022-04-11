package routes

import (
	"nihon-vocabulary/handlers"
	"nihon-vocabulary/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router) {
	user := app.Group("/user", middlewares.Protected())

	user.Get("/", handlers.GetUser)
	user.Put("/", handlers.UpdateUser)
	user.Put("/password", handlers.ChangePasswordUser)
	user.Delete("/", handlers.DeleteUser)
}
