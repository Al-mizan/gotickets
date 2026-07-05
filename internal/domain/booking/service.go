package booking

import (
	"errors"
	"gotickets/internal/apperror"
	"gotickets/internal/domain/booking/dto"
	"gotickets/internal/domain/event"

	"github.com/google/uuid"
)

type Service interface {
	CreateBooking(userId uint, req dto.CreateRequest) (*dto.Response, error)
	GetMyBookings(userId uint) ([]*dto.Response, error)
}

type service struct {
	bookingRepo Repository
	eventRepo   event.Repository
}

func NewService(bookingRepo Repository, eventRepo event.Repository) Service {
	return &service{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func generateBookingCode() string {
	return "GT-" + uuid.New().String()
}

func (s *service) CreateBooking(userId uint, req dto.CreateRequest) (*dto.Response, error) {
	booking := &Booking{
		UserID:      userId,
		EventID:     req.EventID,
		Quantity:    req.Quantity,
		Status:      BookingConfirmed,
		BookingCode: generateBookingCode(),
	}

	err := s.bookingRepo.CreateWithTicketUpdate(booking)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			return nil, apperror.NewNotFound(err, "event not found")
		}
		if errors.Is(err, ErrNotEnoughTickets) {
			return nil, apperror.NewConflict(err, "not enough tickets available")
		}
		return nil, apperror.NewInternal(err, "failed to create booking")
	}

	return booking.ToResponse(), nil
}

func (s *service) GetMyBookings(userId uint) ([]*dto.Response, error) {
	bookings, err := s.bookingRepo.GetByUserID(userId)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch bookings")
	}

	responses := make([]*dto.Response, len(bookings)) // Initialize the slice with the correct length

	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	return responses, nil
}
