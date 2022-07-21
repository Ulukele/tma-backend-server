package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	server, err := NewServer(
		os.Getenv("C_BACKEND_URL"),
		os.Getenv("A_BACKEND_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	apiGroup := app.Group("/api/")

	authGroup := apiGroup.Group("/auth/")
	authGroup.Post("/sign-up/", server.HandleSignUp)
	authGroup.Post("/sign-in/", server.HandleSignIn)

	contentGroup := apiGroup.Group("/content/", server.AuthHeaderValidationMiddleware)
	userGroup := contentGroup.Group("/user/:userId/", server.IdentityValidationMiddleware)
	userGroup.All("/*", server.AccessContentMiddleware)

	log.Fatal(app.Listen(os.Getenv("LISTEN_ON")))
}
