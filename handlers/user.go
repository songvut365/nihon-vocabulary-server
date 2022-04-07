package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	return c.SendString("get user")
}

func UpdateUser(c *fiber.Ctx) error {
	return c.SendString("update user")
}

func ChangePasswordUser(c *fiber.Ctx) error {
	return c.SendString("change password user")
}

func DeleteUser(c *fiber.Ctx) error {
	return c.SendString("delete user")
}
