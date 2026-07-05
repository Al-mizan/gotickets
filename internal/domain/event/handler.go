package event

import (
	"gotickets/internal/apperror"
	"gotickets/internal/domain/event/dto"
	"strconv"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	CreateEvent(c *echo.Context) error
	GetEvents(c *echo.Context) error
	GetEventsByID(c *echo.Context) error
	UpdateEvent(c *echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func (h *handler) CreateEvent(c *echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.CreateEvent(req)
	if err != nil {
		return err
	}

	return c.JSON(201, response)
}

func (h *handler) GetEvents(c *echo.Context) error {
	events, err := h.service.GetEvents()
	if err != nil {
		return err
	}

	return c.JSON(200, events)
}

func (h *handler) GetEventsByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return apperror.NewBadRequest(err, "Invalid event id")
	}

	response, err := h.service.GetEventByID(uint(id))

	if err != nil {
		return err
	}

	return c.JSON(200, response)
}

func (h *handler) UpdateEvent(c *echo.Context) error {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.NewBadRequest(err, "Invalid event id")
	}

	var req dto.UpdateRequest

	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.UpdateEvent(uint(eventId), req)
	if err != nil {
		return err
	}

	return c.JSON(200, response)
}
