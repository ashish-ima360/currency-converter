package controller

import (
	"context"
	"currency-converter/dto"
	"currency-converter/models"
	"currency-converter/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CurrencyService interface {
	CreateCurrency(ctx context.Context, req dto.CurrencyRequest) (*models.Currency, *utils.AppError)
	GetCurrencyByID(ctx context.Context, id int) (*models.Currency, *utils.AppError)
	GetAllCurrencies(ctx context.Context) ([]models.Currency, *utils.AppError)
	UpdateCurrency(ctx context.Context, id int, req dto.CurrencyUpdateRequest) *utils.AppError
	DeleteCurrency(ctx context.Context, id int) *utils.AppError
}

type CurrencyController struct {
	currencyService CurrencyService
}

func NewCurrencyController(currencyService CurrencyService) *CurrencyController {
	return &CurrencyController{
		currencyService: currencyService,
	}
}

func (h *CurrencyController) CreateCurrency(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	result, err := h.currencyService.CreateCurrency(ctx, req)
	if err != nil {
		c.JSON(err.Code, gin.H{
			"error": err.Message,
		})
		return
	}
	resp := dto.CurrencyResponse{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		Symbol:    result.Symbol,
		IsActive:  result.IsActive,
		Deleted:   result.Deleted,
		DeletedAt: result.DeletedAt.Format(time.RFC3339),
		UpdatedAt: result.UpdatedAt.Format(time.RFC3339),
		CreatedAt: result.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *CurrencyController) GetCurrencyByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := utils.ParseIDParam("id", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param",
		})
		return
	}

	result, apperr := h.currencyService.GetCurrencyByID(ctx, id)
	if apperr != nil {
		c.JSON(apperr.Code, gin.H{
			"error": apperr.Message,
		})
		return
	}

	resp := dto.CurrencyResponse{
		ID:        result.ID,
		Code:      result.Code,
		Name:      result.Name,
		Symbol:    result.Symbol,
		IsActive:  result.IsActive,
		Deleted:   result.Deleted,
		DeletedAt: result.DeletedAt.Format(time.RFC3339),
		UpdatedAt: result.UpdatedAt.Format(time.RFC3339),
		CreatedAt: result.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}

func (h *CurrencyController) GetCurrencies(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := h.currencyService.GetAllCurrencies(ctx)
	if err != nil {
		c.JSON(err.Code, gin.H{
			"error": err.Message,
		})
		return
	}

	currencies := make([]dto.CurrencyResponse, 0, len(result))
	for _, currency := range result {
		resp := dto.CurrencyResponse{
			ID:        currency.ID,
			Code:      currency.Code,
			Name:      currency.Name,
			Symbol:    currency.Symbol,
			IsActive:  currency.IsActive,
			Deleted:   currency.Deleted,
			DeletedAt: currency.DeletedAt.Format(time.RFC3339),
			UpdatedAt: currency.UpdatedAt.Format(time.RFC3339),
			CreatedAt: currency.CreatedAt.Format(time.RFC3339),
		}
		currencies = append(currencies, resp)
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "Currencies fetched successfully",
		"currencies": currencies,
	})
}

func (h *CurrencyController) UpdateCurrency(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ParseIDParam("id", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param",
		})
		return
	}

	var req dto.CurrencyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	apperr := h.currencyService.UpdateCurrency(ctx, id, req)
	if apperr != nil {
		c.JSON(apperr.Code, gin.H{
			"error": apperr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
		"message": "Currency updated successfully",
	})
}

func (h *CurrencyController) DeleteCurrency(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ParseIDParam("id", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param",
		})
		return
	}
	apperr := h.currencyService.DeleteCurrency(ctx, id)
	if apperr != nil {
		c.JSON(apperr.Code, gin.H{
			"error": apperr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Currency deleted successfully",
		"id":      id,
	})
}
