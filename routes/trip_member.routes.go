package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupTripMemberRoutes(api fiber.Router, tripMemberColl *mongo.Collection) {
	controller := controllers.NewTripMemberController(tripMemberColl)

	member := api.Group("/trips/:tripId/members")

	member.Post("/new", middleware.ValidateBody[dto.CreateTripMemberRequest](), middleware.Protected(), controller.CreateTripMember)
	member.Get("/", middleware.Protected(), controller.GetTripMembers)
	member.Put("/:id", middleware.ValidateBody[dto.UpdateTripMemberRequest](), middleware.Protected(), controller.UpdateTripMember)
	member.Delete("/:id", middleware.Protected(), controller.DeleteTripMember)
}
