package controllers

import (
	"context"
	"sort"
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
	repo             *repository.TripRepository
	memberCollection *mongo.Collection
}

func NewTripController(coll *mongo.Collection, memberColl *mongo.Collection) *TripController {
	return &TripController{
		repo:             repository.NewTripRepository(coll),
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
	tripObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.BadRequest(c, "Invalid trip ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var trip models.Trip
	err = ctl.repo.FindByID(ctx, id, &trip)
	if err != nil {
		return utils.NotFound(c, "Trip not found")
	}

	memberRepo := repository.NewTripMemberRepository(ctl.memberCollection)
	membersWithUser, err := memberRepo.FindMembersWithUserInfo(ctx, tripObjectID)
	if err != nil {
		return utils.InternalWrap(c, err)
	}

	response := fiber.Map{
		"trip":    trip,
		"members": membersWithUser,
	}

	return utils.Success(c, fiber.StatusOK, "Trip fetched successfully", response, nil)
}

func (ctl *TripController) UpdateTrip(c *fiber.Ctx) error {
	id := c.Params("id")
	var update dto.UpdateTripRequest
	if err := c.BodyParser(&update); err != nil {
		return utils.BadRequest(c, "Invalid input")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get existing trip
	var existing models.Trip
	err := ctl.repo.FindByID(ctx, id, &existing)
	if err != nil {
		return utils.NotFound(c, "Trip not found")
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
		return utils.InternalWrap(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "Trip updated successfully", nil, nil)
}
func (ctl *TripController) DeleteTrip(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := ctl.repo.Delete(ctx, id)
	if err != nil {
		return utils.NotFound(c, "Trip not found")
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

	// 1. Get pagination params
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	perPage := c.QueryInt("limit", 10)
	if perPage < 1 {
		perPage = 10
	}
	skip := (page - 1) * perPage

	// 2. Find trips created by user
	var createdTrips []models.Trip
	createdFilter := bson.M{"createdById": userID}
	if err := t.repo.FindMany(ctx, createdFilter, 0, 0, &createdTrips); err != nil {
		return utils.Internal(c, "Failed to fetch created trips")
	}

	// 3. Find trips where user is a member
	memberRepo := repository.NewTripMemberRepository(t.memberCollection)
	var memberships []models.TripMember
	if err := memberRepo.FindByUserID(ctx, userID, &memberships); err != nil {
		return utils.Internal(c, "Failed to fetch memberships")
	}

	// 4. Collect trip IDs from memberships
	tripIDsMap := make(map[primitive.ObjectID]bool)
	for _, t := range createdTrips {
		tripIDsMap[t.ID] = true
	}
	for _, m := range memberships {
		tripIDsMap[m.TripID] = true
	}

	// 5. Convert to a slice of unique trip IDs
	var allTripIDs []primitive.ObjectID
	for id := range tripIDsMap {
		allTripIDs = append(allTripIDs, id)
	}

	if len(allTripIDs) == 0 {
		// No trips found
		return utils.Success(c, fiber.StatusOK, "No trips found", []models.Trip{}, &utils.Pagination{
			Total:      0,
			Page:       page,
			Limit:      perPage,
			TotalPages: 0,
		})
	}

	// 6. Query trips by allTripIDs
	filter := bson.M{"_id": bson.M{"$in": allTripIDs}}
	var trips []models.Trip
	if err := t.repo.FindMany(ctx, filter, 0, 0, &trips); err != nil {
		return utils.Internal(c, "Failed to fetch trips")
	}

	// 7. Sort trips by CreatedAt descending
	sort.Slice(trips, func(i, j int) bool {
		return trips[i].CreatedAt.After(trips[j].CreatedAt)
	})

	// 8. Pagination
	total := len(trips)
	totalPages := (total + perPage - 1) / perPage
	start := skip
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}
	paginatedTrips := trips[start:end]

	// 9. Return response
	return utils.Success(c, fiber.StatusOK, "Trips fetched successfully", paginatedTrips, &utils.Pagination{
		Total:      total,
		Page:       page,
		Limit:      perPage,
		TotalPages: totalPages,
	})
}
