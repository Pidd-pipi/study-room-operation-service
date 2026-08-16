package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/util"
)

// AchievementRepository 成就徽章仓储。
type AchievementRepository interface {
	Create(achievement *model.Achievement) error
	FindByID(id uint) (*model.Achievement, error)
	List() ([]model.Achievement, error)
	Update(achievement *model.Achievement) error
	CreateUserAchievement(ua *model.UserAchievement) error
	ListUserAchievements(userID uint) ([]model.UserAchievement, error)
	Exists(userID, achievementID uint) (bool, error)
}

type achievementRepository struct {
	db *gorm.DB
}

// NewAchievementRepository 构造成就徽章仓储。
func NewAchievementRepository(db *gorm.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) Create(achievement *model.Achievement) error {
	if err := r.db.Create(achievement).Error; err != nil {
		return fmt.Errorf("create achievement: %w", err)
	}
	return nil
}

func (r *achievementRepository) FindByID(id uint) (*model.Achievement, error) {
	var a model.Achievement
	err := r.db.First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find achievement: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find achievement: %w", err)
	}
	return &a, nil
}

func (r *achievementRepository) List() ([]model.Achievement, error) {
	var list []model.Achievement
	if err := r.db.Order("id asc").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list achievements: %w", err)
	}
	return list, nil
}

func (r *achievementRepository) Update(achievement *model.Achievement) error {
	if err := r.db.Save(achievement).Error; err != nil {
		return fmt.Errorf("update achievement: %w", err)
	}
	return nil
}

func (r *achievementRepository) CreateUserAchievement(ua *model.UserAchievement) error {
	if err := r.db.Create(ua).Error; err != nil {
		return fmt.Errorf("create user achievement: %w", err)
	}
	return nil
}

func (r *achievementRepository) ListUserAchievements(userID uint) ([]model.UserAchievement, error) {
	var list []model.UserAchievement
	if err := r.db.Preload("Achievement").Where("user_id = ?", userID).Order("id desc").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list user achievements: %w", err)
	}
	return list, nil
}

func (r *achievementRepository) Exists(userID, achievementID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&model.UserAchievement{}).Where("user_id = ? AND achievement_id = ?", userID, achievementID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check user achievement: %w", err)
	}
	return count > 0, nil
}
