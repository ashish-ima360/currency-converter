package service

import (
	"context"
	"currency-converter/dto"
	"currency-converter/models"
	"currency-converter/security"
	"currency-converter/utils"
	"net/http"
)

type UserRepository interface {
	GetUserByEmail(context.Context, string) (*models.User, error)
	CreateUser(context.Context, *models.User) (int, error)
}

type userService struct {
	userRepo     UserRepository
	tokenService *security.TokenService
}

func NewUserService(userRepo UserRepository, tokenService *security.TokenService) *userService {
	return &userService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *userService) Register(ctx context.Context, req dto.RegisterRequest) (int, *utils.AppError) {
	// If user doesn't exist, proceed with registration
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return 0, utils.New(http.StatusInternalServerError, "Failed to hash password")
	}

	newUser := models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	userID, err := s.userRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return 0, utils.New(http.StatusInternalServerError, "Failed to create user")
	}

	return userID, nil
}

func (s *userService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResult, *utils.AppError) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResult{}, utils.New(http.StatusUnauthorized, "Invalid credentials")
	}

	// log.Printf("we are here: %v", user)

	ok, err := security.ComparePassword(req.Password, user.PasswordHash)
	if err != nil || !ok {
		return dto.LoginResult{}, utils.New(http.StatusUnauthorized, "Invalid credentials")
	}

	payload := security.RequestClaims{
		UserID: user.ID,
	}

	token, err := s.tokenService.GenerateAccessToken(payload)
	if err != nil {
		return dto.LoginResult{}, utils.New(http.StatusInternalServerError, "Internal server error")
	}

	return dto.LoginResult{
		ID:    user.ID,
		Token: token,
	}, nil
}
