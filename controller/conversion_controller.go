package controller

import (
	"context"
	"currency-converter/dto"
	"currency-converter/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ConversionService interface {
	ConvertCurrency(ctx context.Context, cmd dto.ConversionCmd) (dto.ConversionResult, *utils.AppError)
}

type ConversionController struct {
	conversionService ConversionService
}

func NewConversionController(conversionService ConversionService) *ConversionController {
	return &ConversionController{
		conversionService: conversionService,
	}
}

func (h *ConversionController) ConvertCurrency(c *gin.Context) {
	ctx := c.Request.Context()

	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")

	if from == "" || to == "" || amountStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required query parameters: from, to, amount",
		})
		return
	}

	// make to and from uppercase
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if from == to {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "From and To currencies cannot be the same",
		})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid amount",
		})
		return
	}

	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Amount must be greater than zero",
		})
		return
	}

	result, appError := h.conversionService.ConvertCurrency(ctx, dto.ConversionCmd{
		From:   from,
		To:     to,
		Amount: amount,
	})
	if appError != nil {
		c.JSON(appError.Code, gin.H{
			"error": appError.Message,
		})
		return
	}
	resp := dto.CurrencyConversionResponse{
		From:            from,
		To:              to,
		Amount:          amount,
		Rate:            result.Rate,
		ConvertedAmount: result.ConvertedAmount,
	}

	c.JSON(http.StatusOK, resp)
}
