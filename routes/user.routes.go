package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupUserRoutes(api fiber.Router, tripCollection *mongo.Collection) {
	// tripController := controllers.NewUserController(tripCollection)

	// trip := api.Group("/user")
}
