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

type TripController struct {
	repo *repository.TripRepository
	memberCollection *mongo.Collection
}

func NewTripController(coll *mongo.Collection,memberColl *mongo.Collection) *TripController {
	return &TripController{
		repo: repository.NewTripRepository(coll),
		memberCollection: memberColl,
	}
}
func (ctl *TripController) CreateTrip(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(dto.CreateTripRequest)

	userID := c.Locals("userID").(string)
	createdBy, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return utils.BadRequest(c, "Invalid user id")
	}

	tripID := primitive.NewObjectID()
	now := time.Now()

	trip := models.Trip{
		ID:          tripID,
		Name:        body.Name,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
		CoverPhoto:  body.CoverPhoto,
		CreatedByID: createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create the trip
	_, err = ctl.repo.Create(ctx, trip)
	if err != nil {
		return utils.InternalWrap(c, err)
	}

	// Add the creator as a member with role Admin
	memberRepo := repository.NewTripMemberRepository(ctl.memberCollection)
	member := models.TripMember{
		ID:        primitive.NewObjectID(),
		TripID:    tripID,
		UserID:    createdBy,
		Role:      models.TripMemberRoleAdmin,
		CreatedAt: now,
		UpdatedAt: now,
		JoinedAt:  now,
	}

	_, err = memberRepo.Create(ctx, member)
	if err != nil {
		return utils.InternalWrap(c, err)
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

func (t *TripController) GetMyTrips(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return utils.BadRequest(c, "Invalid user ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get pagination params from query
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	perPage := c.QueryInt("limit", 10)
	if perPage < 1 {
		perPage = 10
	}

	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	filter := bson.M{"createdById": userID}

	// Count total matching documents
	count, err := t.repo.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return utils.Internal(c, "Failed to count trips")
	}

	// Fetch paginated results
	var trips []models.Trip
	err = t.repo.FindMany(ctx, filter, limit, skip, &trips)
	if err != nil {
		return utils.Internal(c, "Failed to fetch trips")
	}

	// Calculate total pages
	totalPages := int((count + int64(perPage) - 1) / int64(perPage))

	return utils.Success(c, fiber.StatusOK, "Trips fetched successfully", trips, &utils.Pagination{
		Total:      int(count),
		Page:       page,
		Limit:    perPage,
		TotalPages: totalPages,
	})
}



