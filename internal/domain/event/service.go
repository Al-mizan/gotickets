package event

import (
	"errors"
	"gotickets/internal/apperror"
	"gotickets/internal/domain/event/dto"
)

type Service interface {
	CreateEvent(req dto.CreateRequest) (*dto.Response, error)
	GetEvents() ([]dto.Response, error)
	GetEventByID(eventId uint) (*dto.Response, error)
	UpdateEvent(eventId uint, req dto.UpdateRequest) (*dto.Response, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateEvent(req dto.CreateRequest) (*dto.Response, error) {
	event := Event{
		Title:            req.Title,
		Description:      req.Description,
		Location:         req.Location,
		StartsAt:         req.StartsAt,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		Price:            req.Price,
	}

	if err := s.repo.Create(&event); err != nil {
		return nil, apperror.NewInternal(err, "failed to create event")
	}

	return event.ToResponse(), nil

}

func (s *service) GetEvents() ([]dto.Response, error) {
	events, err := s.repo.GetAll()
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch events")
	}

	responses := make([]dto.Response, 0, len(events))

	for _, event := range events {
		responses = append(responses, *event.ToResponse())
	}

	return responses, nil
}

func (s *service) GetEventByID(eventId uint) (*dto.Response, error) {
	event, err := s.repo.GetByID(eventId)
	if err != nil {
		if errors.Is(err, ErrEventNotFound) {
			return nil, apperror.NewNotFound(err, "event not found")
		}
		return nil, apperror.NewInternal(err, "failed to fetch event")
	}

	return event.ToResponse(), nil
}

func (s *service) UpdateEvent(eventId uint, req dto.UpdateRequest) (*dto.Response, error) {
	event, err := s.repo.GetByID(eventId) // getting existing event by the id first
	if err != nil {
		if errors.Is(err, ErrEventNotFound) {
			return nil, apperror.NewNotFound(err, "event not found")
		}
		return nil, apperror.NewInternal(err, "failed to fetch event")
	}

	if req.Title != "" {
		event.Title = req.Title
	}

	if req.Description != "" {
		event.Description = req.Description
	}

	if req.Location != "" {
		event.Location = req.Location
	}

	if !req.StartsAt.IsZero() {
		event.StartsAt = req.StartsAt
	}

	if req.Price != 0 {
		event.Price = req.Price
	}

	if err := s.repo.Update(event); err != nil {
		return nil, apperror.NewInternal(err, "failed to update event")
	}

	return event.ToResponse(), nil
}
