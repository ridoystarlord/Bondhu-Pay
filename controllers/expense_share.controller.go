package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/models"
	"github.com/ridoystarlord/bondhu-pay/repository"
	"github.com/ridoystarlord/bondhu-pay/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpenseShareController struct {
	repo *repository.ExpenseShareRepository
}

func NewExpenseShareController(coll *mongo.Collection) *ExpenseShareController {
	return &ExpenseShareController{
		repo: repository.NewExpenseShareRepository(coll),
	}
}

func (e *ExpenseShareController) CreateExpenseShare(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(dto.CreateExpenseShareRequest)
	expenseIDParam := c.Params("expenseId")

	expenseID, err := primitive.ObjectIDFromHex(expenseIDParam)
	if err != nil {
		return utils.BadRequest(c, "Invalid expense ID")
	}

	userID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		return utils.BadRequest(c, "Invalid user ID")
	}

	now := time.Now()
	share := models.ExpenseShare{
		ID:        primitive.NewObjectID(),
		ExpenseID: expenseID,
		UserID:    userID,
		Amount:    body.Amount,
		Settled:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = e.repo.Create(ctx, share)
	if err != nil {
		return utils.Internal(c, "Failed to create expense share")
	}

	return utils.Success(c, fiber.StatusCreated, "Expense share created", share, nil)
}

func (e *ExpenseShareController) GetExpenseShares(c *fiber.Ctx) error {
	expenseIDParam := c.Params("expenseId")

	expenseID, err := primitive.ObjectIDFromHex(expenseIDParam)
	if err != nil {
		return utils.BadRequest(c, "Invalid expense ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var shares []models.ExpenseShare
	// Fixed: method name to match repository method
	err = e.repo.FindByExpense(ctx, expenseID, &shares)
	if err != nil {
		return utils.Internal(c, "Failed to fetch expense shares")
	}

	return utils.Success(c, fiber.StatusOK, "Expense shares fetched", shares, nil)
}

func (e *ExpenseShareController) UpdateExpenseShare(c *fiber.Ctx) error {
	shareID := c.Params("id")
	body := c.Locals("validatedBody").(dto.UpdateExpenseShareRequest)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{}
	if body.Amount != 0 {
		update["amount"] = body.Amount
	}
	if body.Settled != nil {
		update["settled"] = *body.Settled
		if *body.Settled {
			now := time.Now()
			update["settledAt"] = &now
		} else {
			update["settledAt"] = nil
		}
	}
	if body.SettledVia != "" {
		update["settledVia"] = body.SettledVia
	}
	if body.TransactionID != "" {
		update["transactionId"] = body.TransactionID
	}
	update["updatedAt"] = time.Now()

	_, err := e.repo.Update(ctx, shareID, update)
	if err != nil {
		return utils.Internal(c, "Failed to update expense share")
	}

	return utils.Success(c, fiber.StatusOK, "Expense share updated", nil, nil)
}

func (e *ExpenseShareController) DeleteExpenseShare(c *fiber.Ctx) error {
	shareID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := e.repo.Delete(ctx, shareID)
	if err != nil {
		return utils.Internal(c, "Failed to delete expense share")
	}

	return utils.Success(c, fiber.StatusOK, "Expense share deleted", nil, nil)
}
