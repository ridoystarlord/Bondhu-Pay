package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"
	"github.com/ridoystarlord/bondhu-pay/config"
	"github.com/ridoystarlord/bondhu-pay/routes"
	"github.com/ridoystarlord/bondhu-pay/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// Check if it's a known Fiber error (e.g., 404)
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Log unexpected errors
			if code == fiber.StatusInternalServerError {
				log.Printf("Unhandled error: %v", err)
				return utils.Internal(c, "Internal Server Error")
			}

			// For known errors, return your custom error response
			return utils.Error(c, code, err.Error())
		},
	})

	// Security headers
	app.Use(helmet.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5001",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

	// Request logging
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return utils.Error(c, fiber.StatusTooManyRequests, "Too many requests, please try again later.")
		},
	}))

	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	log.Fatal(app.Listen(":" + port))
}
