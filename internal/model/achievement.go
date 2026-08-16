package model

import "time"

// Achievement 成就徽章定义。
type Achievement struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Code          string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name          string    `gorm:"size:64;not null" json:"name"`
	ConditionType string    `gorm:"size:32;not null" json:"condition_type"`
	Threshold     int       `gorm:"not null" json:"threshold"`
	Icon          string    `gorm:"size:128" json:"icon"`
	Description   string    `gorm:"size:255" json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

// UserAchievement 用户已获得徽章。
type UserAchievement struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	UserID        uint         `gorm:"index;not null" json:"user_id"`
	AchievementID uint         `gorm:"index;not null" json:"achievement_id"`
	Achievement   *Achievement `gorm:"foreignKey:AchievementID" json:"achievement,omitempty"`
	AchievedAt    time.Time    `json:"achieved_at"`
}
