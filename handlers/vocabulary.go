package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func CreateVocabulary(c *fiber.Ctx) error {
	return c.SendString("create vocabulary")
}

func GetVocabularies(c *fiber.Ctx) error {
	return c.SendString("get vocabularies")
}

func GetVocabulary(c *fiber.Ctx) error {
	return c.SendString("get vocabulary")
}

func UpdateVocabulary(c *fiber.Ctx) error {
	return c.SendString("update vocabulary")
}

func DeleteVocabulary(c *fiber.Ctx) error {
	return c.SendString("delete vocabulary")
}
