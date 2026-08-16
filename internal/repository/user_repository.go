package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/model"
)

// ErrDuplicate 唯一键冲突。
var ErrDuplicate = errors.New("duplicate key")

// UserRepository 用户仓储。
type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByIDTx(tx *gorm.DB, id uint) (*model.User, error)
	Update(user *model.User) error
	UpdateTx(tx *gorm.DB, user *model.User) error
	ListRanking(period string, limit int) ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if isDuplicate(err) {
			return fmt.Errorf("create user: %w", ErrDuplicate)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return &u, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	return r.FindByIDTx(nil, id)
}

func (r *userRepository) FindByIDTx(tx *gorm.DB, id uint) (*model.User, error) {
	var u model.User
	err := dbOrTx(r.db, tx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.UpdateTx(nil, user)
}

func (r *userRepository) UpdateTx(tx *gorm.DB, user *model.User) error {
	if err := dbOrTx(r.db, tx).Save(user).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// ListRanking 按累计学习时长排名。
func (r *userRepository) ListRanking(_ string, limit int) ([]model.User, error) {
	var users []model.User
	if err := r.db.Where("role = ?", "user").Order("total_study_minutes desc").Limit(limit).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("list ranking: %w", err)
	}
	return users, nil
}
