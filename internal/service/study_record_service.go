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

// RankingItem 排行榜条目。
type RankingItem struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Duration int    `json:"duration"`
}

// StudyRecordService 学习时长业务逻辑。
type StudyRecordService interface {
	Ranking(period string, limit int) ([]RankingItem, error)
	Stats(userID uint) (*util.StudyStats, error)
}

type studyRecordService struct {
	studyRepo repository.StudyRecordRepository
	logger    *slog.Logger
}

// NewStudyRecordService 构造学习时长服务。
func NewStudyRecordService(studyRepo repository.StudyRecordRepository, logger *slog.Logger) StudyRecordService {
	return &studyRecordService{studyRepo: studyRepo, logger: logger}
}

func (s *studyRecordService) Ranking(period string, limit int) ([]RankingItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var records []model.StudyRecord
	var err error
	switch period {
	case "day":
		records, err = s.studyRepo.DailyRanking(time.Now().Format("2006-01-02"), limit)
	case "week":
		records, err = s.studyRepo.WeeklyRanking(limit)
	default:
		records, err = s.studyRepo.MonthlyRanking(limit)
	}
	if err != nil {
		return nil, fmt.Errorf("ranking[%s]: %w", period, err)
	}
	items := make([]RankingItem, 0, len(records))
	for _, r := range records {
		nickname := "学员"
		if r.User != nil {
			nickname = r.User.Nickname
			if nickname == "" {
				nickname = r.User.Username
			}
		}
		items = append(items, RankingItem{UserID: r.UserID, Nickname: nickname, Duration: r.DurationMinutes})
	}
	s.logger.Info(constants.LogStudyRankingQueried, "period", period, "count", len(items))
	return items, nil
}

func (s *studyRecordService) Stats(userID uint) (*util.StudyStats, error) {
	total, err := s.studyRepo.SumByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("study stats user[%d]: %w", userID, err)
	}
	s.logger.Info(constants.LogStudyStatsQueried, "user_id", userID, "total", total)
	return &util.StudyStats{TotalMinutes: total}, nil
}
