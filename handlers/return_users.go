package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omeiirr/dummy-data-api/helpers"
)

func ReturnUsers(c *fiber.Ctx) error {
	users := helpers.GenerateUsers(5)
	return c.JSON(users)
}
