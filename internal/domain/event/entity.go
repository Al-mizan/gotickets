package event

import (
	"gotickets/internal/domain/event/dto"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title            string    `gorm:"type:varchar(150);not null"`
	Description      string    `gorm:"type:text"`
	Location         string    `gorm:"type:varchar(150);not null"`
	StartsAt         time.Time `gorm:"not null"`
	TotalTickets     int       `gorm:"not null"`
	AvailableTickets int       `gorm:"not null"`
	Price            int       `gorm:"not null"`
}

func (e *Event) ToResponse() *dto.Response {
	return &dto.Response{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		Location:         e.Location,
		StartsAt:         e.StartsAt,
		TotalTickets:     e.TotalTickets,
		AvailableTickets: e.AvailableTickets,
		Price:            e.Price,
		CreatedAt:        e.CreatedAt,
	}
}
