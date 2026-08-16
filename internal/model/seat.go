package model

import (
	"time"

	"github.com/ld/studyroom/internal/constants"
)

// Seat 自习室座位。
type Seat struct {
	ID        uint                 `gorm:"primaryKey" json:"id"`
	SeatNo    string               `gorm:"size:32;uniqueIndex;not null" json:"seat_no"`
	Floor     int                  `gorm:"not null;default:1" json:"floor"`
	Zone      constants.SeatZone   `gorm:"size:16;index;not null" json:"zone"`
	SeatType  string               `gorm:"size:32" json:"seat_type"`
	Status    constants.SeatStatus `gorm:"size:16;index;not null" json:"status"`
	X         int                  `json:"x"`
	Y         int                  `json:"y"`
	Remark    string               `gorm:"size:255" json:"remark"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
