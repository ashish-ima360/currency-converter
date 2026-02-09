package controller

import (
	"context"
	"currency-converter/dto"
	"currency-converter/models"
	"currency-converter/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ExchangeRateService interface {
	CreateExchangeRate(ctx context.Context, req dto.ExchangeRateRequest) (*models.ExchangeRate, *utils.AppError)
	GetExchangeRateByID(ctx context.Context, id int) (*models.ExchangeRate, *utils.AppError)
	GetAllExchangeRates(ctx context.Context) ([]models.ExchangeRate, *utils.AppError)
	UpdateExchangeRate(ctx context.Context, id int, rate float64) (*models.ExchangeRate, *utils.AppError)
	DeleteExchangeRate(ctx context.Context, id int) *utils.AppError
}

type ExchangeRateController struct {
	exchangeRateService ExchangeRateService
}

func NewExchangeRateController(exchangeRateService ExchangeRateService) *ExchangeRateController {
	return &ExchangeRateController{
		exchangeRateService: exchangeRateService,
	}
}

func (h *ExchangeRateController) CreateExchangeRate(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.ExchangeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}
	if req.ToCurrencyID == req.FromCurrencyID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "FromCurrencyID and ToCurrencyID cannot be the same",
		})
		return
	}

	exchangeRate, err := h.exchangeRateService.CreateExchangeRate(ctx, req)
	if err != nil {
		c.JSON(err.Code, gin.H{
			"error": err.Message,
		})
		return
	}
	resp := dto.ExchangeRateResponse{
		ID:             exchangeRate.ID,
		FromCurrencyID: exchangeRate.FromCurrencyID,
		ToCurrencyID:   exchangeRate.ToCurrencyID,
		Rate:           exchangeRate.Rate,
		IsActive:       exchangeRate.IsActive,
		Deleted:        exchangeRate.Deleted,
		DeletedAt:      exchangeRate.DeletedAt.Format(time.RFC3339),
		UpdatedAt:      exchangeRate.UpdatedAt.Format(time.RFC3339),
		CreatedAt:      exchangeRate.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *ExchangeRateController) GetExchangeRateByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := utils.ParseIDParam("id", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param",
		})
		return
	}

	exchangeRate, appErr := h.exchangeRateService.GetExchangeRateByID(ctx, id)
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{
			"error": appErr.Message,
		})
		return
	}

	resp := dto.ExchangeRateResponse{
		ID:             exchangeRate.ID,
		FromCurrencyID: exchangeRate.FromCurrencyID,
		ToCurrencyID:   exchangeRate.ToCurrencyID,
		Rate:           exchangeRate.Rate,
		IsActive:       exchangeRate.IsActive,
		Deleted:        exchangeRate.Deleted,
		DeletedAt:      exchangeRate.DeletedAt.Format(time.RFC3339),
		CreatedAt:      exchangeRate.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      exchangeRate.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ExchangeRateController) GetAllExchangeRates(c *gin.Context) {
	ctx := c.Request.Context()

	result, appErr := h.exchangeRateService.GetAllExchangeRates(ctx)
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{
			"error": appErr.Message,
		})
		return
	}
	exchangeRates := make([]dto.ExchangeRateResponse, 0, len(result))
	for _, rate := range result {
		exchangeRates = append(exchangeRates, dto.ExchangeRateResponse{
			ID:             rate.ID,
			FromCurrencyID: rate.FromCurrencyID,
			ToCurrencyID:   rate.ToCurrencyID,
			Rate:           rate.Rate,
			IsActive:       rate.IsActive,
			Deleted:        rate.Deleted,
			DeletedAt:      rate.DeletedAt.Format(time.RFC3339),
			CreatedAt:      rate.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      rate.UpdatedAt.Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Exchange rates fetched successfully",
		"exchange_rates": exchangeRates,
	})
}

func (h *ExchangeRateController) UpdateExchangeRate(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := utils.ParseIDParam("id", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param",
		})
		return
	}

	rate := c.Query("rate")
	if rate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Rate parameter is required",
		})
		return
	}

	rateFloat, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid rate value",
		})
		return
	}

	exchangeRate, appErr := h.exchangeRateService.UpdateExchangeRate(ctx, id, rateFloat)
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{
			"error": appErr.Message,
		})
		return
	}

	resp := dto.ExchangeRateResponse{
		ID:             exchangeRate.ID,
		FromCurrencyID: exchangeRate.FromCurrencyID,
		ToCurrencyID:   exchangeRate.ToCurrencyID,
		Rate:           exchangeRate.Rate,
		IsActive:       exchangeRate.IsActive,
		Deleted:        exchangeRate.Deleted,
		DeletedAt:      exchangeRate.DeletedAt.Format(time.RFC3339),
		CreatedAt:      exchangeRate.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      exchangeRate.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ExchangeRateController) DeleteExchangeRate(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := utils.ParseIDParam("id", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID param",
		})
		return
	}

	appErr := h.exchangeRateService.DeleteExchangeRate(ctx, id)
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{
			"error": appErr.Message,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Exchange rate deleted successfully",
	})
}
