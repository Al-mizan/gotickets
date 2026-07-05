package booking

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, handler Handler, authMw echo.MiddlewareFunc) {
	api := e.Group("/api/v1/bookings", authMw)

	api.POST("", handler.CreateBooking)
	api.GET("/me", handler.GetMyBookings)
}
