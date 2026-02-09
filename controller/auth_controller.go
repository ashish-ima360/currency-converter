package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"currency-converter/dto"
	"currency-converter/utils"
)

type UserService interface {
	Login(context.Context, dto.LoginRequest) (dto.LoginResult, *utils.AppError)
	Register(context.Context, dto.RegisterRequest) (int, *utils.AppError)
}

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	result, err := h.userService.Register(ctx, req)
	if err != nil {
		c.JSON(err.Code, gin.H{
			"error": err.Message,
		})
		return
	}
	resp := dto.RegisterResponse{
		UserId:  result,
		Message: "User Registered Successfully",
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	result, err := h.userService.Login(ctx, req)
	if err != nil {
		c.JSON(err.Code, gin.H{
			"error": err.Message,
		})
		return
	}

	resp := dto.LoginResponse{
		UserId:  result.ID,
		Token:   result.Token,
		Message: "User Logged in Successfully",
	}

	c.JSON(http.StatusOK, resp)
}
