package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupTripRoutes(api fiber.Router, tripCollection *mongo.Collection,tripMemberCollection *mongo.Collection) {
	tripController := controllers.NewTripController(tripCollection,tripMemberCollection)

	trip := api.Group("/trips")

	trip.Post("/new", middleware.ValidateBody[dto.CreateTripRequest](), middleware.Protected(), tripController.CreateTrip)
	trip.Get("/", middleware.Protected(), tripController.GetMyTrips)
	trip.Get("/:id", middleware.Protected(), tripController.GetTrip)
	trip.Put("/:id", middleware.ValidateBody[dto.UpdateTripRequest](), middleware.Protected(), tripController.UpdateTrip)
	trip.Delete("/:id", middleware.Protected(), tripController.DeleteTrip)
}
