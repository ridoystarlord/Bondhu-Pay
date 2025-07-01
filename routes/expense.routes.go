package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/controllers"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupExpenseRoutes(api fiber.Router, expenseColl *mongo.Collection) {
	controller := controllers.NewExpenseController(expenseColl)

	expense := api.Group("/trips/:tripId/expenses")

	expense.Post("/new", 
		middleware.ValidateBody[dto.CreateExpenseRequest](), 
		middleware.Protected(), 
		controller.CreateExpense,
	)

	expense.Get("/", 
		middleware.Protected(), 
		controller.GetExpensesByTrip,
	)

	expense.Put("/:id", 
		middleware.ValidateBody[dto.UpdateExpenseRequest](), 
		middleware.Protected(), 
		controller.UpdateExpense,
	)

	expense.Delete("/:id", 
		middleware.Protected(), 
		controller.DeleteExpense,
	)
}
