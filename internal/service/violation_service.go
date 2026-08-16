package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/repository"
)

// ViolationService 违约记录业务逻辑。
type ViolationService interface {
	Create(userID, bookingID uint, violationType string, points int, blockedHours int, remark string) (*model.ViolationRecord, error)
	List(userID uint, page, pageSize int) ([]model.ViolationRecord, int64, error)
	BlockedUntil(userID uint) (*time.Time, error)
}

type violationService struct {
	violationRepo repository.ViolationRepository
	logger        *slog.Logger
}

// NewViolationService 构造违约记录服务。
func NewViolationService(violationRepo repository.ViolationRepository, logger *slog.Logger) ViolationService {
	return &violationService{violationRepo: violationRepo, logger: logger}
}

func (s *violationService) Create(userID, bookingID uint, violationType string, points int, blockedHours int, remark string) (*model.ViolationRecord, error) {
	var blockedUntil *time.Time
	if blockedHours > 0 {
		t := time.Now().Add(time.Duration(blockedHours) * time.Hour)
		blockedUntil = &t
	}
	record := &model.ViolationRecord{UserID: userID, BookingID: bookingID, ViolationType: violationType, Points: points, BlockedUntil: blockedUntil, Remark: remark}
	if err := s.violationRepo.Create(record); err != nil {
		return nil, fmt.Errorf("create violation: %w", err)
	}
	s.logger.Info(constants.LogViolationCreated, "user_id", userID, "type", violationType)
	return record, nil
}

func (s *violationService) List(userID uint, page, pageSize int) ([]model.ViolationRecord, int64, error) {
	list, total, err := s.violationRepo.List(userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list violations: %w", err)
	}
	s.logger.Info(constants.LogViolationListQueried, "user", userID, "total", total)
	return list, total, nil
}

func (s *violationService) BlockedUntil(userID uint) (*time.Time, error) {
	blocked, err := s.violationRepo.ActiveBlockedUntil(userID)
	if err != nil {
		return nil, fmt.Errorf("blocked until user[%d]: %w", userID, err)
	}
	if blocked != nil {
		s.logger.Warn(constants.LogUserBlocked, "user_id", userID, "until", blocked.Format("2006-01-02 15:04"))
	}
	return blocked, nil
}
