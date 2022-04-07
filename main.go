package main

import (
	"nihon-vocabulary/configs"
	"nihon-vocabulary/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	//setup .env
	configs.SetupEnv()

	//connect database
	configs.ConnectDB()

	//setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.AuthenticationRoutes(v1) // PATH api/v1/auth
	routes.UserRoutes(v1)           // PATH api/v1/user
	routes.VocabularyRoutes(v1)     // PATH api/v1/vocabulary

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
