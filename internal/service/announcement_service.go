package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/repository"
	"github.com/ld/studyroom/internal/util"
)

// AnnouncementService 公告业务逻辑。
type AnnouncementService interface {
	Create(title, content string, createdBy uint) (*model.Announcement, error)
	Update(id uint, title, content string, isActive bool) (*model.Announcement, error)
	Delete(id uint) error
	List(page, pageSize int, activeOnly bool) ([]model.Announcement, int64, error)
}

type announcementService struct {
	announcementRepo repository.AnnouncementRepository
	logger           *slog.Logger
}

// NewAnnouncementService 构造公告服务。
func NewAnnouncementService(announcementRepo repository.AnnouncementRepository, logger *slog.Logger) AnnouncementService {
	return &announcementService{announcementRepo: announcementRepo, logger: logger}
}

func (s *announcementService) Create(title, content string, createdBy uint) (*model.Announcement, error) {
	if title == "" || content == "" {
		return nil, fmt.Errorf("create announcement: %w", util.ErrValidation)
	}
	announcement := &model.Announcement{Title: title, Content: content, IsActive: true, PublishedAt: time.Now(), CreatedBy: createdBy}
	if err := s.announcementRepo.Create(announcement); err != nil {
		return nil, fmt.Errorf("create announcement: %w", err)
	}
	s.logger.Info(constants.LogAnnouncementCreated, "announcement_id", announcement.ID, "title", title)
	return announcement, nil
}

func (s *announcementService) Update(id uint, title, content string, isActive bool) (*model.Announcement, error) {
	announcement, err := s.announcementRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update announcement[id=%d]: %w", id, err)
	}
	if title != "" {
		announcement.Title = title
	}
	if content != "" {
		announcement.Content = content
	}
	announcement.IsActive = isActive
	if err := s.announcementRepo.Update(announcement); err != nil {
		return nil, fmt.Errorf("update announcement[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogAnnouncementUpdated, "announcement_id", id)
	return announcement, nil
}

func (s *announcementService) Delete(id uint) error {
	if err := s.announcementRepo.Delete(id); err != nil {
		return fmt.Errorf("delete announcement[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogAnnouncementDeleted, "announcement_id", id)
	return nil
}

func (s *announcementService) List(page, pageSize int, activeOnly bool) ([]model.Announcement, int64, error) {
	list, total, err := s.announcementRepo.List(page, pageSize, activeOnly)
	if err != nil {
		return nil, 0, fmt.Errorf("list announcements: %w", err)
	}
	s.logger.Info(constants.LogAnnouncementList, "total", total)
	return list, total, nil
}
