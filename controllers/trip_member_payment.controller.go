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

type TripMemberPaymentController struct {
	repo *repository.TripMemberPaymentRepository
}

func NewTripMemberPaymentController(coll *mongo.Collection) *TripMemberPaymentController {
	return &TripMemberPaymentController{
		repo: repository.NewTripMemberPaymentRepository(coll),
	}
}

func (ctl *TripMemberPaymentController) Create(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(dto.CreateTripMemberPaymentRequest)
	tripIDParam := c.Params("tripId")

	tripID, err := primitive.ObjectIDFromHex(tripIDParam)
	if err != nil {
		return utils.BadRequest(c, "Invalid trip ID")
	}
	memberID, err := primitive.ObjectIDFromHex(body.MemberID)
	if err != nil {
		return utils.BadRequest(c, "Invalid Member ID")
	}

	now := time.Now()
	payment := models.TripMemberPayment{
		ID:        primitive.NewObjectID(),
		TripID:    tripID,
		Amount:    body.Amount,
		Note:      body.Note,
		CreatedAt: now,
		UpdatedAt: now,
		PaidAt:    now,
		MemberID:  memberID,
		Method:   models.PaymentMethod(body.Method),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = ctl.repo.Create(ctx, payment)
	if err != nil {
		return utils.InternalWrap(c, err)
	}

	return utils.Success(c, fiber.StatusCreated, "Payment recorded", payment, nil)
}

func (ctl *TripMemberPaymentController) ListByTrip(c *fiber.Ctx) error {
	tripIDParam := c.Params("tripId")

	tripID, err := primitive.ObjectIDFromHex(tripIDParam)
	if err != nil {
		return utils.BadRequest(c, "Invalid trip ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payments []models.TripMemberPayment
	err = ctl.repo.FindManyByTrip(ctx, tripID, &payments)
	if err != nil {
		return utils.Internal(c, "Failed to fetch payments")
	}

	return utils.Success(c, fiber.StatusOK, "Payments fetched", payments, nil)
}

func (ctl *TripMemberPaymentController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	body := c.Locals("validatedBody").(dto.UpdateTripMemberPaymentRequest)

	update := bson.M{"updatedAt": time.Now()}
	if body.Amount > 0 {
		update["amount"] = body.Amount
	}
	if body.Method != "" {
		update["method"] = body.Method
	}
	if body.Note != "" {
		update["note"] = body.Note
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := ctl.repo.Update(ctx, id, update)
	if err != nil {
		return utils.Internal(c, "Failed to update payment")
	}

	return utils.Success(c, fiber.StatusOK, "Payment updated", nil, nil)
}

func (ctl *TripMemberPaymentController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := ctl.repo.Delete(ctx, id)
	if err != nil {
		return utils.Internal(c, "Failed to delete payment")
	}

	return utils.Success(c, fiber.StatusOK, "Payment deleted", nil, nil)
}
