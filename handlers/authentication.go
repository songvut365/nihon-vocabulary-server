package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	return c.SendString("login")
}

func Logout(c *fiber.Ctx) error {
	return c.SendString("logout")
}

func Register(c *fiber.Ctx) error {
	return c.SendString("register")
}

func ResetPassword(c *fiber.Ctx) error {
	return c.SendString("reset-password")
}
