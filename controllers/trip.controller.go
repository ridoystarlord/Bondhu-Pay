package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/dto"
	"github.com/ridoystarlord/bondhu-pay/models"
	"github.com/ridoystarlord/bondhu-pay/repository"
	"github.com/ridoystarlord/bondhu-pay/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripController struct {
	repo *repository.TripRepository
}

func NewTripController(coll *mongo.Collection) *TripController {
	return &TripController{
		repo: repository.NewTripRepository(coll),
	}
}
func (ctl *TripController) CreateTrip(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(dto.CreateTripRequest)

	userID := c.Locals("userId").(string)
	createdBy, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return utils.BadRequest(c,"Invalid user id")
	}

	trip := models.Trip{
		ID:        primitive.NewObjectID(),
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		CoverPhoto: body.CoverPhoto,
		CreatedByID: createdBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = ctl.repo.Create(ctx, trip)
	if err != nil {
		return utils.InternalWrap(c,err)
	}

	return utils.Success(c, fiber.StatusCreated, "Trip created successfully", trip, nil)
}
func (ctl *TripController) GetTrip(c *fiber.Ctx) error {
	id := c.Params("id")
	var trip models.Trip

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ctl.repo.FindByID(ctx, id, &trip)
	if err != nil {
		return utils.NotFound(c,"Trip not found")
	}

	return utils.Success(c, fiber.StatusOK, "Trip fetched successfully", trip, nil)
}
func (ctl *TripController) UpdateTrip(c *fiber.Ctx) error {
	id := c.Params("id")
	var update dto.UpdateTripRequest
	if err := c.BodyParser(&update); err != nil {
		return utils.BadRequest(c,"Invalid input")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get existing trip
	var existing models.Trip
	err := ctl.repo.FindByID(ctx, id, &existing)
	if err != nil {
		return utils.NotFound(c,"Trip not found")
	}

	// Prepare update document
	updateDoc := make(map[string]interface{})
	if update.Name != "" {
		updateDoc["name"] = update.Name
	}
	if !update.StartDate.IsZero() {
		updateDoc["startDate"] = update.StartDate
	}
	if !update.EndDate.IsZero() {
		updateDoc["endDate"] = update.EndDate
	}
	if update.CoverPhoto != "" {
		updateDoc["coverPhoto"] = update.CoverPhoto
	}
	updateDoc["updatedAt"] = time.Now()

	_, err = ctl.repo.Update(ctx, id, updateDoc)
	if err != nil {
		return utils.InternalWrap(c,err)
	}

	return utils.Success(c, fiber.StatusOK, "Trip updated successfully", nil, nil)
}
func (ctl *TripController) DeleteTrip(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := ctl.repo.Delete(ctx, id)
	if err != nil {
		return utils.NotFound(c,"Trip not found")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
