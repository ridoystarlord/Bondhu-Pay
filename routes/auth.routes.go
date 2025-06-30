package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupAuthRoutes(api fiber.Router, userCollection *mongo.Collection) {
	userController := controllers.NewAuthController(userCollection)

	auth := api.Group("/auth")
	auth.Post("/register", middleware.ValidateBody[dto.RegisterRequest](), userController.Register)
	auth.Post("/login", middleware.ValidateBody[dto.LoginRequest](), userController.Login)
}
