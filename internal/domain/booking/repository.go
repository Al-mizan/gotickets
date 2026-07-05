package booking

import (
	"errors"
	"gotickets/internal/domain/event"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrBookingNotFound         = errors.New("booking not found")
	ErrNotEnoughTickets        = errors.New("not enough tickets available")
	ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
	ErrForbiddenBookingAccess  = errors.New("you do not own this booking")
)

type Repository interface {
	Create(booking *Booking) error
	GetByID(bookingId uint) (*Booking, error)
	GetByUserID(userId uint) ([]*Booking, error)
	Update(booking *Booking) error
	CreateWithTicketUpdate(booking *Booking) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(booking *Booking) error {
	return r.db.Create(booking).Error
}

func (r *repository) GetByID(bookingId uint) (*Booking, error) {
	var booking Booking

	err := r.db.First(&booking, bookingId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}

		return nil, err
	}

	return &booking, nil
}

func (r *repository) GetByUserID(userId uint) ([]*Booking, error) {
	var bookings []*Booking

	err := r.db.Where("user_id = ?", userId).Find(&bookings).Error
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *repository) Update(booking *Booking) error {
	return r.db.Save(booking).Error
}

func (r *repository) CreateWithTicketUpdate(booking *Booking) error {
	// start transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		var eventData event.Event

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&eventData, booking.EventID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return event.ErrEventNotFound
			}
			return err
		}

		if eventData.AvailableTickets < booking.Quantity {
			return ErrNotEnoughTickets
		}

		booking.TotalPrice = booking.Quantity * eventData.Price

		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		eventData.AvailableTickets -= booking.Quantity
		if err := tx.Save(&eventData).Error; err != nil {
			return err
		}

		return nil
	})
}
