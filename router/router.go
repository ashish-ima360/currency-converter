package router

import (
	"currency-converter/controller"
	"currency-converter/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authMiddleware *middleware.AuthMiddleware,
	userController *controller.UserController,
	currencyController *controller.CurrencyController,
	exchangeRateController *controller.ExchangeRateController,
	conversionController *controller.ConversionController,
) *gin.Engine {

	r := gin.Default()

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	r.Use(authMiddleware.Handle())

	r.POST("/currencies", currencyController.CreateCurrency)
	r.GET("/currencies", currencyController.GetCurrencies)
	r.GET("/currencies/:id", currencyController.GetCurrencyByID)
	r.PATCH("/currencies/:id", currencyController.UpdateCurrency)
	r.DELETE("/currencies/:id", currencyController.DeleteCurrency)

	r.POST("/exchange-rates", exchangeRateController.CreateExchangeRate)
	r.GET("/exchange-rates", exchangeRateController.GetAllExchangeRates)
	r.GET("/exchange-rates/:id", exchangeRateController.GetExchangeRateByID)
	r.PATCH("/exchange-rates/:id", exchangeRateController.UpdateExchangeRate) 
	r.DELETE("/exchange-rates/:id", exchangeRateController.DeleteExchangeRate)
	r.POST("/exchange-rates/sync/:code", exchangeRateController.SyncExchangeRates) // /exchange-rates/sync/USD

	r.GET("/convert", conversionController.ConvertCurrency) // ?from=USD&to=INR&amount=100

	return r
}
