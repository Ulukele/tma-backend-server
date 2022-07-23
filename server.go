package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
	"strconv"
)

type Server struct {
	contentURL string
	authURL    string
}

func NewServer(contentURL string, authURL string) (*Server, error) {
	s := &Server{}
	log.Printf("set content api url to %s", contentURL)
	s.contentURL = contentURL
	log.Printf("set auth api url to %s", authURL)
	s.authURL = authURL

	return s, nil
}

func (s *Server) StartApp() error {
	app := fiber.New()
	app.Use(cors.New())

	apiGroup := app.Group("/api/v1/")

	authGroup := apiGroup.Group("/auth/")
	authGroup.All("/*", s.AccessAuthMiddleware)

	contentGroup := apiGroup.Group("/content/", s.AuthHeaderValidationMiddleware)
	userGroup := contentGroup.Group("/user/")
	teamGroup := contentGroup.Group("/team/")

	userGroup.All("/*", s.AccessAuthMiddleware)
	teamGroup.All("/*", s.AccessContentMiddleware)

	return app.Listen(os.Getenv("LISTEN_ON"))
}

func (s *Server) AuthHeaderValidationMiddleware(c *fiber.Ctx) error {
	jwt := c.Get("Bearer", "")
	userInfo, err := s.requestValidation(jwt)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, "invalid jwt")
	}
	c.Locals("userId", userInfo.Id)
	return c.Next()
}

func (s *Server) AccessContentMiddleware(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	headers["UserId"] = strconv.Itoa(int(c.Locals("userId").(uint)))
	content, statusCode, err := s.requestContent(c.Method(), c.Path()[1:], c.Body(), headers)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "unable to forward")
	}
	if statusCode != 200 {
		return fiber.NewError(statusCode, string(content))
	}
	return c.Send(content)
}

func (s *Server) AccessAuthMiddleware(c *fiber.Ctx) error {
	content, statusCode, err := s.requestAuth(c.Method(), c.Path()[1:], c.Body(), c.GetReqHeaders())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "unable to forward")
	}
	if statusCode != 200 {
		return fiber.NewError(statusCode, string(content))
	}
	return c.Send(content)
}
