package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"strings"
)

type Server struct {
	contentURL string
	authURL    string
}

// Validator
var validate = validator.New()

func NewServer(contentURL string, authURL string) (*Server, error) {
	s := &Server{}
	log.Printf("set content api url to %s", contentURL)
	s.contentURL = contentURL
	log.Printf("set auth api url to %s", authURL)
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

	auth, err := s.requestSignIn(res.Username, res.UserId)
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

	auth, err := s.requestSignIn(user.Username, user.UserId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't create jwt")
	}

	return c.JSON(auth)
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

func (s *Server) IdentityValidationMiddleware(c *fiber.Ctx) error {
	userIdAuth := c.Locals("userId").(uint)
	userIdProvided, err := strconv.Atoi(c.Params("userId", ""))
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Sprintf("no user with id=%s", c.Params("userId", "")),
		)
	}
	if userIdAuth != uint(userIdProvided) {
		return fiber.NewError(fiber.StatusForbidden, "jwt mismatch with user id")
	}

	return c.Next()
}

func (s *Server) AccessContentMiddleware(c *fiber.Ctx) error {
	// /api/content/...
	keys := strings.Split(c.Path(), "/")
	keys[1] = "api"
	keys[2] = "public"

	// api/public/...
	newPath := strings.Join(keys[1:], "/")
	content, err := s.requestContent(c.Method(), newPath, c.Body(), c.GetReqHeaders())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "unable to forward")
	}
	return c.Send(content)
}
