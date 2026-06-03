package booking

import (
	"gotickets/internal/event"
)

type service struct {
	bookingRepo Repository
	eventRepo   event.Repository
}

func NewService(bookingRepo Repository) *service {
	return &service{
		bookingRepo: bookingRepo,
	}
}
