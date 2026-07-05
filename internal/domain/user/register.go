package user

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, handler Handler, authMw echo.MiddlewareFunc) {
	api := e.Group("/api/v1/auth")

	api.POST("/register", handler.CreateUser)
	api.POST("/login", handler.LoginUser)
	api.GET("/me", handler.GetMe, authMw)
}
