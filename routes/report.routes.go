package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupTripReportRoutes(api fiber.Router, expenseColl, shareColl, paymentColl, userColl *mongo.Collection) {
	controller := controllers.NewTripReportController(expenseColl, shareColl, paymentColl, userColl)

	report := api.Group("/trips/:tripId/report")

	report.Get("/", middleware.Protected(), controller.GetTripReport)
}
