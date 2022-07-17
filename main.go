package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"strconv"
	"strings"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	server, err := NewServer("http://localhost:8081/api", "http://localhost:8080/api")
	if err != nil {
		log.Fatal(err)
	}

	apiGroup := app.Group("/api/")

	authGroup := apiGroup.Group("/auth/")
	authGroup.Post("/sign-up/", server.HandleSignUp)
	authGroup.Post("/sign-in/", server.HandleSignIn)

	contentGroup := apiGroup.Group("/content/", func(c *fiber.Ctx) error {
		jwt := c.Get("Bearer", "")
		userInfo, err := server.requestValidation(jwt)
		if err != nil {
			return fiber.NewError(fiber.StatusForbidden, "invalid jwt")
		}
		c.Locals("userId", userInfo.Id)
		return c.Next()
	})
	userGroup := contentGroup.Group("/user/:userId/", func(c *fiber.Ctx) error {
		userIdAuth := c.Locals("userId").(uint)
		userIdProvided, err := strconv.Atoi(c.Params("userId", ""))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("no user with id=%s", c.Params("userId", "")))
		}
		if userIdAuth != uint(userIdProvided) {
			return fiber.NewError(fiber.StatusForbidden, "jwt mismatch with user id")
		}

		return c.Next()
	})
	userGroup.All("/*", func(c *fiber.Ctx) error {
		keys := strings.Split(c.Path(), "/")
		keys[2] = "public"
		newPath := "/" + strings.Join(keys[2:], "/") // TODO <- do this better
		content, err := server.requestContent(c.Method(), newPath, c.Body(), "")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "unable to forward")
		}
		return c.Send(content)
	})

	log.Fatal(app.Listen(":8082"))
}
