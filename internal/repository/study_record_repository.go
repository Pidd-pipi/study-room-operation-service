package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/model"
)

// StudyRecordRepository 学习时长记录仓储。
type StudyRecordRepository interface {
	Create(record *model.StudyRecord) error
	CreateTx(tx *gorm.DB, record *model.StudyRecord) error
	SumByUser(userID uint) (int, error)
	DailyRanking(date string, limit int) ([]model.StudyRecord, error)
	WeeklyRanking(limit int) ([]model.StudyRecord, error)
	MonthlyRanking(limit int) ([]model.StudyRecord, error)
}

type studyRecordRepository struct {
	db *gorm.DB
}

// NewStudyRecordRepository 构造学习时长记录仓储。
func NewStudyRecordRepository(db *gorm.DB) StudyRecordRepository {
	return &studyRecordRepository{db: db}
}

func (r *studyRecordRepository) Create(record *model.StudyRecord) error {
	return r.CreateTx(nil, record)
}

func (r *studyRecordRepository) CreateTx(tx *gorm.DB, record *model.StudyRecord) error {
	if err := dbOrTx(r.db, tx).Create(record).Error; err != nil {
		return fmt.Errorf("create study record: %w", err)
	}
	return nil
}

func (r *studyRecordRepository) SumByUser(userID uint) (int, error) {
	var total int64
	if err := r.db.Model(&model.StudyRecord{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration_minutes),0)").Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("sum study minutes: %w", err)
	}
	return int(total), nil
}

func (r *studyRecordRepository) DailyRanking(date string, limit int) ([]model.StudyRecord, error) {
	var list []model.StudyRecord
	err := r.db.Preload("User").
		Where("study_date = ?", date).
		Order("duration_minutes desc").Limit(limit).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("daily ranking: %w", err)
	}
	return list, nil
}

func (r *studyRecordRepository) WeeklyRanking(limit int) ([]model.StudyRecord, error) {
	var list []model.StudyRecord
	err := r.db.Preload("User").
		Where("study_date >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)").
		Order("duration_minutes desc").Limit(limit).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("weekly ranking: %w", err)
	}
	return list, nil
}

func (r *studyRecordRepository) MonthlyRanking(limit int) ([]model.StudyRecord, error) {
	var list []model.StudyRecord
	err := r.db.Preload("User").
		Where("study_date >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)").
		Order("duration_minutes desc").Limit(limit).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("monthly ranking: %w", err)
	}
	return list, nil
}
