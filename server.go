package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Server struct {
	contentURL string
	authURL    string
}

// Validator
var validate = validator.New()

func NewServer(contentURL string, authURL string) (*Server, error) {
	s := &Server{}
	s.contentURL = contentURL
	s.authURL = authURL

	return s, nil
}

func (s *Server) HandleSignIn(c *fiber.Ctx) error {
	log.Printf("handle sign-in at %s", c.Path())
	req := AuthRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect username and password")
	}
	err := validate.Struct(req)
	if err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	res, err := s.getUserFull(req.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no such user")
	}
	if res.Username != req.Username || res.Password != req.Password {
		return fiber.NewError(fiber.StatusBadRequest, "wrong username or password")
	}

	auth, err := s.requestSignIn(res.Username, res.UserId, res.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't create jwt")
	}

	return c.JSON(auth)
}

func (s *Server) HandleSignUp(c *fiber.Ctx) error {
	log.Printf("handle sign-up at %s", c.Path())
	req := AuthRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect username and password")
	}
	err := validate.Struct(req)
	if err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	user, err := s.registerUser(req.Username, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't register user")
	}

	auth, err := s.requestSignIn(user.Username, user.UserId, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't create jwt")
	}

	return c.JSON(auth)
}
