package booking

import (
	"gotickets/internal/domain/booking/dto"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
)

type Booking struct {
	gorm.Model
	UserID      uint          `gorm:"not null;constraint:OnDelete:CASCADE"`
	EventID     uint          `gorm:"not null;constraint:OnDelete:CASCADE"`
	Quantity    int           `gorm:"not null"`
	TotalPrice  int           `gorm:"not null"`
	Status      BookingStatus `gorm:"type:varchar(50);not null"`
	BookingCode string        `gorm:"uniqueIndex;not null"`
}

func (b *Booking) ToResponse() *dto.Response {
	return &dto.Response{
		ID:          b.ID,
		UserID:      b.UserID,
		EventID:     b.EventID,
		Quantity:    b.Quantity,
		TotalPrice:  b.TotalPrice,
		Status:      string(b.Status),
		BookingCode: b.BookingCode,
		CreatedAt:   b.CreatedAt,
	}
}
