package booking

import (
	"gotickets/internal/apperror"
	"gotickets/internal/ctxkeys"
	"gotickets/internal/domain/booking/dto"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	CreateBooking(c *echo.Context) error
	GetMyBookings(c *echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func getCurrentUserID(c *echo.Context) (uint, bool) {
	userId, ok := c.Get(string(ctxkeys.UserID)).(uint)
	return userId, ok
}

func (h *handler) CreateBooking(c *echo.Context) error {
	userId, ok := getCurrentUserID(c)
	if !ok {
		return apperror.NewUnauthorized(nil, "Unauthorized")
	}

	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.CreateBooking(userId, req)
	if err != nil {
		return err
	}

	return c.JSON(201, response)
}

func (h *handler) GetMyBookings(c *echo.Context) error {
	userId, ok := getCurrentUserID(c)
	if !ok {
		return apperror.NewUnauthorized(nil, "Unauthorized")
	}

	bookings, err := h.service.GetMyBookings(userId)
	if err != nil {
		return err
	}

	return c.JSON(200, bookings)
}
