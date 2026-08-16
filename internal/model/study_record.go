package model

import "time"

// StudyRecord 学习时长记录。
type StudyRecord struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	User            *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookingID       uint      `gorm:"index" json:"booking_id"`
	StudyDate       string    `gorm:"size:16;index" json:"study_date"`
	DurationMinutes int       `gorm:"not null" json:"duration_minutes"`
	CreatedAt       time.Time `json:"created_at"`
}
