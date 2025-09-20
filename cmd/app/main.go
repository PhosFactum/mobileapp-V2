// @title ClinicHub API
// @version 1.0.0
// @description API для работы с приёмами пациентов

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer your_jwt_token_here"

// @host localhost:8080
// @BasePath /api/v1
package main

import (
	_ "github.com/AlexanderMorozov1919/mobileapp/docs"
	_ "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/handlers"
	"github.com/AlexanderMorozov1919/mobileapp/internal/app"
)

func main() {
	app.New().Run()
}
