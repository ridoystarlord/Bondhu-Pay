package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/models"
	"github.com/ridoystarlord/bondhu-pay/repository"
	"github.com/ridoystarlord/bondhu-pay/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripReportController struct {
	expenseRepo      *repository.ExpenseRepository
	expenseShareRepo *repository.ExpenseShareRepository
	paymentRepo      *repository.TripMemberPaymentRepository
	userRepo         *repository.UserRepository
}

func NewTripReportController(
	expenseColl *mongo.Collection,
	shareColl *mongo.Collection,
	paymentColl *mongo.Collection,
	userColl *mongo.Collection,
) *TripReportController {
	return &TripReportController{
		expenseRepo:      repository.NewExpenseRepository(expenseColl),
		expenseShareRepo: repository.NewExpenseShareRepository(shareColl),
		paymentRepo:      repository.NewTripMemberPaymentRepository(paymentColl),
		userRepo:         repository.NewUserRepository(userColl),
	}
}

func (ctl *TripReportController) GetTripReport(c *fiber.Ctx) error {
	tripIDParam := c.Params("tripId")
	tripID, err := primitive.ObjectIDFromHex(tripIDParam)
	if err != nil {
		return utils.BadRequest(c, "Invalid trip ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1️⃣ Get all expenses
	var expenses []models.Expense
	if err := ctl.expenseRepo.FindByTrip(ctx, tripID, &expenses); err != nil {
		return utils.Internal(c, "Failed to load expenses")
	}

	// 2️⃣ Get all shares
	var shares []models.ExpenseShare
	if err := ctl.expenseShareRepo.FindManyByTrip(ctx, tripID, &shares); err != nil {
		return utils.Internal(c, "Failed to load expense shares")
	}

	// 3️⃣ Get all payments
	var payments []models.TripMemberPayment
	if err := ctl.paymentRepo.FindManyByTrip(ctx, tripID, &payments); err != nil {
		return utils.Internal(c, "Failed to load payments")
	}

	// Collect unique user IDs from shares and payments and expenses' paidBy
	userIDSet := map[primitive.ObjectID]struct{}{}

	for _, s := range shares {
		userIDSet[s.UserID] = struct{}{}
	}

	for _, p := range payments {
		userIDSet[p.MemberID] = struct{}{}
	}

	for _, e := range expenses {
		userIDSet[e.PaidBy] = struct{}{}
	}

	userIDs := make([]primitive.ObjectID, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	// 4️⃣ Fetch users info
	users, err := ctl.userRepo.FindByIDs(ctx, userIDs)
	if err != nil {
		return utils.Internal(c, "Failed to load user info")
	}

	userMap := make(map[primitive.ObjectID]models.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	// 5️⃣ Build per-user aggregation
	type memberDataStruct struct {
		TotalOwed float64
		TotalPaid float64
	}

	memberData := make(map[primitive.ObjectID]*memberDataStruct)

	// Sum shares (owed)
	for _, s := range shares {
		if memberData[s.UserID] == nil {
			memberData[s.UserID] = &memberDataStruct{}
		}
		memberData[s.UserID].TotalOwed += s.Amount
	}

	// Sum payments
	for _, p := range payments {
		if memberData[p.MemberID] == nil {
			memberData[p.MemberID] = &memberDataStruct{}
		}
		memberData[p.MemberID].TotalPaid += p.Amount
	}

	// Sum direct expenses paid by members
	for _, e := range expenses {
		if memberData[e.PaidBy] == nil {
			memberData[e.PaidBy] = &memberDataStruct{}
		}
		memberData[e.PaidBy].TotalPaid += e.Amount
	}

	// 6️⃣ Build response
	totalExpense := 0.0
	for _, e := range expenses {
		totalExpense += e.Amount
	}

	var members []fiber.Map
	for userID, data := range memberData {
		user := userMap[userID]

		members = append(members, fiber.Map{
			"userId":       userID.Hex(),
			"name":         user.Name,
			"mobileNumber": user.MobileNumber,
			"totalOwed":    data.TotalOwed,
			"totalPaid":    data.TotalPaid,
			"netBalance":   data.TotalPaid - data.TotalOwed,
		})
	}

	response := fiber.Map{
		"tripId":       tripID.Hex(),
		"totalExpense": totalExpense,
		"members":      members,
	}

	return utils.Success(c, fiber.StatusOK, "Trip report generated", response, nil)
}
