package model

import (
	"time"

	"github.com/ld/studyroom/internal/constants"
)

// Booking 座位预约。
type Booking struct {
	ID           uint                    `gorm:"primaryKey" json:"id"`
	BookingNo    string                  `gorm:"size:32;uniqueIndex;not null" json:"booking_no"`
	UserID       uint                    `gorm:"index;not null" json:"user_id"`
	User         *User                   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SeatID       uint                    `gorm:"index;not null" json:"seat_id"`
	Seat         *Seat                   `gorm:"foreignKey:SeatID" json:"seat,omitempty"`
	BookingDate  string                  `gorm:"size:16;index" json:"booking_date"`
	StartHour    int                     `gorm:"not null" json:"start_hour"`
	EndHour      int                     `gorm:"not null" json:"end_hour"`
	Status       constants.BookingStatus `gorm:"size:16;index;not null" json:"status"`
	QRCode       string                  `gorm:"size:64" json:"qr_code"`
	CheckInTime  *time.Time              `json:"check_in_time"`
	CheckOutTime *time.Time              `json:"check_out_time"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}
