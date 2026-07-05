package event

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, handler Handler, authMw echo.MiddlewareFunc) {
	api := e.Group("/api/v1/events")

	api.POST("", handler.CreateEvent, authMw)
	api.GET("", handler.GetEvents)
	api.GET("/:id", handler.GetEventsByID)
	api.PATCH("/:id", handler.UpdateEvent, authMw)
}
