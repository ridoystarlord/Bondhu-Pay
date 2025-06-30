package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupTripRoutes(api fiber.Router, tripCollection *mongo.Collection) {
	tripController := controllers.NewTripController(tripCollection)

	trip := api.Group("/trips")

	trip.Post("/", middleware.ValidateBody[dto.CreateTripRequest](), middleware.Protected(), tripController.CreateTrip)
	trip.Get("/:id", middleware.Protected(), tripController.GetTrip)
	trip.Put("/:id", middleware.ValidateBody[dto.UpdateTripRequest](), middleware.Protected(), tripController.UpdateTrip)
	trip.Delete("/:id", middleware.Protected(), tripController.DeleteTrip)
}
