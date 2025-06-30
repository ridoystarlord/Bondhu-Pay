package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/config"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/middleware"
)

func SetupRoutes(app *fiber.App) {
	
	api := app.Group("/api")
	
	auth := api.Group("/auth")
	auth.Post("/register", middleware.ValidateBody[dto.RegisterRequest](), controllers.Register)
	auth.Post("/login", middleware.ValidateBody[dto.LoginRequest](), controllers.Login)

	tripCollection := config.GetCollection("trips")
	SetupTripRoutes(api, tripCollection)

}
