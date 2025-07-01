package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupExpenseShareRoutes(api fiber.Router, expenseShareColl *mongo.Collection) {
	controller := controllers.NewExpenseShareController(expenseShareColl)

	share := api.Group("/expenses/:expenseId/shares")

	share.Post("/new", middleware.ValidateBody[dto.CreateExpenseShareRequest](), middleware.Protected(), controller.CreateExpenseShare)
	share.Get("/", middleware.Protected(), controller.GetExpenseShares)
	share.Put("/:id", middleware.ValidateBody[dto.UpdateExpenseShareRequest](), middleware.Protected(), controller.UpdateExpenseShare)
	share.Delete("/:id", middleware.Protected(), controller.DeleteExpenseShare)
}
