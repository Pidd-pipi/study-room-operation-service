package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/util"
)

// AnnouncementRepository 公告仓储。
type AnnouncementRepository interface {
	Create(announcement *model.Announcement) error
	FindByID(id uint) (*model.Announcement, error)
	List(page, pageSize int, activeOnly bool) ([]model.Announcement, int64, error)
	Update(announcement *model.Announcement) error
	Delete(id uint) error
}

type announcementRepository struct {
	db *gorm.DB
}

// NewAnnouncementRepository 构造公告仓储。
func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) Create(announcement *model.Announcement) error {
	if err := r.db.Create(announcement).Error; err != nil {
		return fmt.Errorf("create announcement: %w", err)
	}
	return nil
}

func (r *announcementRepository) FindByID(id uint) (*model.Announcement, error) {
	var a model.Announcement
	err := r.db.First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find announcement: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find announcement: %w", err)
	}
	return &a, nil
}

func (r *announcementRepository) List(page, pageSize int, activeOnly bool) ([]model.Announcement, int64, error) {
	var list []model.Announcement
	var total int64
	q := r.db.Model(&model.Announcement{})
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count announcements: %w", err)
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("published_at desc").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list announcements: %w", err)
	}
	return list, total, nil
}

func (r *announcementRepository) Update(announcement *model.Announcement) error {
	if err := r.db.Save(announcement).Error; err != nil {
		return fmt.Errorf("update announcement: %w", err)
	}
	return nil
}

func (r *announcementRepository) Delete(id uint) error {
	res := r.db.Delete(&model.Announcement{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete announcement: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete announcement: %w", util.ErrNotFound)
	}
	return nil
}
