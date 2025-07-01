package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupTripMemberPaymentRoutes(api fiber.Router, paymentColl *mongo.Collection) {
	ctrl := controllers.NewTripMemberPaymentController(paymentColl)

	payment := api.Group("/trips/:tripId/payments")

	payment.Post("/", middleware.Protected(), middleware.ValidateBody[dto.CreateTripMemberPaymentRequest](), ctrl.Create)
	payment.Get("/", middleware.Protected(), ctrl.ListByTrip)
	payment.Put("/:id", middleware.Protected(), middleware.ValidateBody[dto.UpdateTripMemberPaymentRequest](), ctrl.Update)
	payment.Delete("/:id", middleware.Protected(), ctrl.Delete)
}
