package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/model"
)

// ViolationRepository 违约记录仓储。
type ViolationRepository interface {
	Create(record *model.ViolationRecord) error
	CreateTx(tx *gorm.DB, record *model.ViolationRecord) error
	List(userID uint, page, pageSize int) ([]model.ViolationRecord, int64, error)
	ActiveBlockedUntil(userID uint) (*time.Time, error)
}

type violationRepository struct {
	db *gorm.DB
}

// NewViolationRepository 构造违约记录仓储。
func NewViolationRepository(db *gorm.DB) ViolationRepository {
	return &violationRepository{db: db}
}

func (r *violationRepository) Create(record *model.ViolationRecord) error {
	return r.CreateTx(nil, record)
}

func (r *violationRepository) CreateTx(tx *gorm.DB, record *model.ViolationRecord) error {
	if err := dbOrTx(r.db, tx).Create(record).Error; err != nil {
		return fmt.Errorf("create violation: %w", err)
	}
	return nil
}

func (r *violationRepository) List(userID uint, page, pageSize int) ([]model.ViolationRecord, int64, error) {
	var list []model.ViolationRecord
	var total int64
	q := r.db.Model(&model.ViolationRecord{})
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count violations: %w", err)
	}
	if err := q.Preload("User").Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list violations: %w", err)
	}
	return list, total, nil
}

// ActiveBlockedUntil 返回当前生效中的黑名单截止时间。
func (r *violationRepository) ActiveBlockedUntil(userID uint) (*time.Time, error) {
	var record model.ViolationRecord
	err := r.db.Where("user_id = ? AND blocked_until > ?", userID, time.Now()).Order("blocked_until desc").First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("active blocked until: %w", err)
	}
	return record.BlockedUntil, nil
}
