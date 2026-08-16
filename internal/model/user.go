package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/constants"
)

// User 学员/管理员。
type User struct {
	ID                uint               `gorm:"primaryKey" json:"id"`
	Username          string             `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash      string             `gorm:"size:255;not null" json:"-"`
	Nickname          string             `gorm:"size:64" json:"nickname"`
	Avatar            string             `gorm:"size:255" json:"avatar"`
	Role              constants.UserRole `gorm:"size:16;index;not null" json:"role"`
	Phone             string             `gorm:"size:32" json:"phone"`
	TotalStudyMinutes int                `gorm:"not null;default:0" json:"total_study_minutes"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if !u.Role.Valid() {
		return constants.ErrInvalidUserRole
	}
	return nil
}
