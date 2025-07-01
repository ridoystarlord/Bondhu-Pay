package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/utils"
)

// Generic validation middleware for Go 1.18+
func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			return utils.BadRequest(c, "Invalid request body")
		}

		if err := utils.ValidateStruct(body); err != nil {
			return utils.BadRequest(c, err.Error())
		}

		c.Locals("validatedBody", body)
		return c.Next()
	}
}
