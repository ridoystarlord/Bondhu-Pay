package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Pagination struct
type Pagination struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalPages int `json:"totalPages"`
}

// SUCCESS RESPONSE (Unified)
func Success(c *fiber.Ctx, statusCode int, message string, data interface{}, pagination *Pagination) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success":    true,
		"statusCode": statusCode,
		"message":    message,
		"data":       data,
		"error":      nil,
		"pagination": pagination,
	})
}

// ERROR RESPONSES (standard helpers)
func BadRequest(c *fiber.Ctx, message string) error {
	return errorResponse(c, 400, message)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return errorResponse(c, 401, message)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return errorResponse(c, 403, message)
}

func NotFound(c *fiber.Ctx, message string) error {
	return errorResponse(c, 404, message)
}

func Internal(c *fiber.Ctx, message string) error {
	return errorResponse(c, 500, message)
}

func InternalWrap(c *fiber.Ctx, err error) error {
	if err != nil {
		return errorResponse(c, 500, fmt.Sprintf("Internal error: %s", err.Error()))
	}
	return errorResponse(c, 500, "Internal error")
}

// GENERIC CUSTOM ERROR RESPONSE
func Error(c *fiber.Ctx, statusCode int, message string) error {
	return errorResponse(c, statusCode, message)
}

// Core private function
func errorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success":    false,
		"statusCode": statusCode,
		"message":    message,
		"data":       nil,
		"error":      message,
		"pagination": nil,
	})
}
