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

type UserController struct {
	repo *repository.UserRepository
}

func NewAuthController(coll *mongo.Collection) *UserController {
	return &UserController{
		repo: repository.NewUserRepository(coll),
	}
}

// Register handles user registration
func (uc *UserController) Register(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(dto.RegisterRequest)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if phone already exists
	existing, err := uc.repo.FindByPhone(ctx, body.Phone)
	if err == nil && existing != nil {
		return utils.BadRequest(c, "Phone already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return utils.Internal(c, "Failed to hash password")
	}

	// Create new user
	user := models.User{
		ID:           primitive.NewObjectID(),
		Name:         body.Name,
		Phone:        body.Phone,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = uc.repo.Create(ctx, user)
	if err != nil {
		return utils.Internal(c, "Failed to create user")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		return utils.Internal(c, "Failed to generate token")
	}

	return utils.Success(c, fiber.StatusCreated, "User registered successfully", fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"phone": user.Phone,
		},
	}, nil)
}

// Login handles user authentication
func (uc *UserController) Login(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(dto.LoginRequest)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by phone
	user, err := uc.repo.FindByPhone(ctx, body.Phone)
	if err != nil || user == nil {
		return utils.Unauthorized(c, "Invalid phone or password")
	}

	// Check password
	if !utils.CheckPasswordHash(body.Password, user.PasswordHash) {
		return utils.Unauthorized(c, "Invalid phone or password")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		return utils.Internal(c, "Failed to generate token")
	}

	return utils.Success(c, fiber.StatusOK, "User logged in successfully", fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"phone": user.Phone,
		},
	}, nil)
}
