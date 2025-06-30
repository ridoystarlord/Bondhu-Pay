package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/config"
)

func SetupRoutes(app *fiber.App) {
	
	api := app.Group("/api")
	
	userCollection := config.GetCollection("users")
	tripCollection := config.GetCollection("trips")
	tripMemberCollection := config.GetCollection("trip_members")

	SetupAuthRoutes(api, userCollection)
	SetupUserRoutes(api, userCollection)
	SetupTripRoutes(api, tripCollection,tripMemberCollection)
	SetupTripMemberRoutes(api, tripMemberCollection)

}
