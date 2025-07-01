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

type ExpenseController struct {
	repo           *repository.ExpenseRepository
	expenseShareRepo *repository.ExpenseShareRepository
}

func NewExpenseController(expenseColl *mongo.Collection, expenseShareColl *mongo.Collection) *ExpenseController {
	return &ExpenseController{
		repo:             repository.NewExpenseRepository(expenseColl),
		expenseShareRepo: repository.NewExpenseShareRepository(expenseShareColl),
	}
}

func (ctl *ExpenseController) CreateExpense(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(dto.CreateExpenseRequest)

	tripID, err := primitive.ObjectIDFromHex(body.TripID)
	if err != nil {
		return utils.BadRequest(c, "Invalid trip ID")
	}
	paidBy, err := primitive.ObjectIDFromHex(body.PaidBy)
	if err != nil {
		return utils.BadRequest(c, "Invalid paidBy ID")
	}

	now := time.Now()
	expenseID := primitive.NewObjectID()

	expense := models.Expense{
		ID:        expenseID,
		TripID:    tripID,
		Amount:    body.Amount,
		PaidBy:    paidBy,
		Category:  models.ExpenseCategory(body.Category),
		Note:      body.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert expense
	_, err = ctl.repo.Create(ctx, expense)
	if err != nil {
		return utils.Internal(c, "Failed to create expense")
	}

	// Build shares
	var shares []interface{}
	for _, s := range body.Shares {
		uid, err := primitive.ObjectIDFromHex(s.UserID)
		if err != nil {
			return utils.BadRequest(c, "Invalid share user ID")
		}
		share := models.ExpenseShare{
			ID:        primitive.NewObjectID(),
			ExpenseID: expenseID,
			UserID:    uid,
			Amount:    s.Amount,
			Settled:   false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		shares = append(shares, share)
	}

	// Insert shares
	if len(shares) > 0 {
		_, err = ctl.expenseShareRepo.CreateMany(ctx, shares)
		if err != nil {
			return utils.Internal(c, "Failed to create expense shares")
		}
	}

	return utils.Success(c, fiber.StatusCreated, "Expense created", expense, nil)
}

func (ctl *ExpenseController) GetExpensesByTrip(c *fiber.Ctx) error {
	tripIDParam := c.Params("tripId")
	tripID, err := primitive.ObjectIDFromHex(tripIDParam)
	if err != nil {
		return utils.BadRequest(c, "Invalid trip ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var expenses []models.Expense
	err = ctl.repo.FindByTrip(ctx, tripID, &expenses)
	if err != nil {
		return utils.Internal(c, "Failed to fetch expenses")
	}

	return utils.Success(c, fiber.StatusOK, "Expenses fetched", expenses, nil)
}

func (ctl *ExpenseController) UpdateExpense(c *fiber.Ctx) error {
	id := c.Params("id")
	body := c.Locals("validatedBody").(dto.UpdateExpenseRequest)

	update := bson.M{
		"updatedAt": time.Now(),
	}
	if body.Amount > 0 {
		update["amount"] = body.Amount
	}
	if body.Category != "" {
		update["category"] = body.Category
	}
	if body.Note != "" {
		update["note"] = body.Note
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := ctl.repo.Update(ctx, id, update)
	if err != nil {
		return utils.Internal(c, "Failed to update expense")
	}

	return utils.Success(c, fiber.StatusOK, "Expense updated", nil, nil)
}

func (ctl *ExpenseController) DeleteExpense(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := ctl.repo.Delete(ctx, id)
	if err != nil {
		return utils.Internal(c, "Failed to delete expense")
	}

	return utils.Success(c, fiber.StatusOK, "Expense deleted", nil, nil)
}
