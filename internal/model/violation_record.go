package model

import "time"

// ViolationRecord 违约黑名单记录。
type ViolationRecord struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"index;not null" json:"user_id"`
	User          *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookingID     uint       `gorm:"index" json:"booking_id"`
	ViolationType string     `gorm:"size:32;not null" json:"violation_type"`
	Points        int        `gorm:"not null;default:0" json:"points"`
	BlockedUntil  *time.Time `json:"blocked_until"`
	Remark        string     `gorm:"size:255" json:"remark"`
	CreatedAt     time.Time  `json:"created_at"`
}
