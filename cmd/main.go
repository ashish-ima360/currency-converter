package main

import (
	"fmt"
	"log"

	"currency-converter/config"
	"currency-converter/controller"
	"currency-converter/db"
	"currency-converter/middleware"
	"currency-converter/repository"
	"currency-converter/router"
	"currency-converter/security"
	"currency-converter/service"
)
// main function acts as orchestration layer where 
// we initialize all the dependencies, 
// handle all fatal errors and start the server
func main() {    
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error in loading config: %v", err)
	}

	// Connect to DB
	dbConn, err := db.ConnectDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("error in connecting to DB: %v", err)
	}

	// do all migrations in DB
	if err := db.Migrate(dbConn); err != nil {
		log.Fatalf("error in migration in DB: %v", err)
	}

	// Inject dependencies -> 
	// Currently we are injecting dependencies in main, but in future we use a DI container for better management of dependencies
	
	// create repositories
	userRepo := repository.NewUserRepository(dbConn)
	currencyRepo := repository.NewCurrencyRepository(dbConn)
	exchangeRateRepo := repository.NewExchangeRateRepository(dbConn)

	// create services
	tokenService := security.NewTokenService(&cfg.AuthConfig)

	userService := service.NewUserService(userRepo, tokenService)
	currencyService := service.NewCurrencyService(currencyRepo)
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo)
	conversionService := service.NewConversionService(currencyRepo, exchangeRateRepo)

	// create controllers
	userController := controller.NewUserController(userService)
	currencyController := controller.NewCurrencyController(currencyService)
	exchangeRateController := controller.NewExchangeRateController(exchangeRateService)
	conversionController := controller.NewConversionController(conversionService)

	// create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// Setup Routes
	r := router.SetupRouter(authMiddleware, userController, currencyController, exchangeRateController, conversionController)

	// Run Server 
	// In production we should do graceful shutdown of server
	log.Printf("server listening on Port : %v", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%v", cfg.Port)); err != nil {
		log.Fatalf("error in runnin server: %v", err)
	}
}
