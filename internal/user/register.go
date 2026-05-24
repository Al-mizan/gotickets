package user

import (
	"gotickets/internal/auth"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userRepository := NewRepository(db)
	jwtService := auth.NewJWTService("") // you can pass a custom secret key here
	userService := NewService(userRepository, jwtService)
	userHandler := NewHandler(userService)

	api := e.Group("/api/v1/auth")

	api.POST("/register", userHandler.CreateUser) // api/v1/auth/register
	api.POST("/login", userHandler.LoginUser)     // api/v1/auth/login
}
