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

type TripMemberController struct {
	repo *repository.TripMemberRepository
}

func NewTripMemberController(coll *mongo.Collection) *TripMemberController {
	return &TripMemberController{
		repo: repository.NewTripMemberRepository(coll),
	}
}

func (t *TripMemberController) CreateTripMember(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(dto.CreateTripMemberRequest)
	tripIDParam := c.Params("tripId")

	tripID, err := primitive.ObjectIDFromHex(tripIDParam)
	if err != nil {
		return utils.BadRequest(c, "Invalid trip ID")
	}

	userID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		return utils.BadRequest(c, "Invalid user ID")
	}

	now := time.Now()
	member := models.TripMember{
		ID:        primitive.NewObjectID(),
		TripID:    tripID,
		UserID:    userID,
		Role:      models.TripMemberRole(body.Role),
		CreatedAt: now,
		UpdatedAt: now,
		JoinedAt:  now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = t.repo.Create(ctx, member)
	if err != nil {
		return utils.Internal(c, "Failed to create trip member")
	}

	return utils.Success(c, fiber.StatusCreated, "Trip member added", member, nil)
}

func (t *TripMemberController) GetTripMembers(c *fiber.Ctx) error {
	tripIDParam := c.Params("tripId")

	tripID, err := primitive.ObjectIDFromHex(tripIDParam)
	if err != nil {
		return utils.BadRequest(c, "Invalid trip ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var members []models.TripMember
	err = t.repo.FindManyByTrip(ctx, tripID, &members)
	if err != nil {
		return utils.Internal(c, "Failed to fetch trip members")
	}

	return utils.Success(c, fiber.StatusOK, "Trip members fetched", members, nil)
}

func (t *TripMemberController) UpdateTripMember(c *fiber.Ctx) error {
	memberID := c.Params("id")
	body := c.Locals("validatedBody").(dto.UpdateTripMemberRequest)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"role":      body.Role,
		"updatedAt": time.Now(),
	}

	_, err := t.repo.Update(ctx, memberID, update)
	if err != nil {
		return utils.Internal(c, "Failed to update trip member")
	}

	return utils.Success(c, fiber.StatusOK, "Trip member updated", nil, nil)
}

func (t *TripMemberController) DeleteTripMember(c *fiber.Ctx) error {
	memberID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := t.repo.Delete(ctx, memberID)
	if err != nil {
		return utils.Internal(c, "Failed to delete trip member")
	}

	return utils.Success(c, fiber.StatusOK, "Trip member deleted", nil, nil)
}
