package model

import "time"

// Announcement 公告。
type Announcement struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:128;not null" json:"title"`
	Content     string    `gorm:"type:text" json:"content"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	PublishedAt time.Time `json:"published_at"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
