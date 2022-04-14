package main

import (
	"nihon-vocabulary/configs"
	"nihon-vocabulary/routes"

	_ "nihon-vocabulary/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title Nihon Vocabulary
// @version 1.0
// @description This is an API for Nihon Vocabulary Application

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name ApiKeyAuth

// @contact.name Songvut Nakrong
// @contact.email songvut.nakrong@gmail.com

// @BasePath /api/v1/
func main() {
	app := fiber.New()
	app.Use(logger.New())

	//setup .env
	configs.SetupEnv()

	//connect database
	configs.ConnectDB()

	//setup routesapp
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/docs/*", swagger.HandlerDefault) //swagger PATH api/v1/docs/

	routes.AuthenticationRoutes(v1) // PATH api/v1/auth
	routes.UserRoutes(v1)           // PATH api/v1/user
	routes.VocabularyRoutes(v1)     // PATH api/v1/vocabulary

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
