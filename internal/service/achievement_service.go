package service

import (
	"fmt"
	"log/slog"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/repository"
	"github.com/ld/studyroom/internal/util"
)

// AchievementService 成就徽章业务逻辑。
type AchievementService interface {
	Create(code, name, conditionType string, threshold int, icon, description string) (*model.Achievement, error)
	List() ([]model.Achievement, error)
	MyAchievements(userID uint) ([]model.UserAchievement, error)
	CheckAndGrant(userID uint, studyMinutes int) error
}

type achievementService struct {
	achievementRepo repository.AchievementRepository
	logger          *slog.Logger
}

// NewAchievementService 构造成就徽章服务。
func NewAchievementService(achievementRepo repository.AchievementRepository, logger *slog.Logger) AchievementService {
	return &achievementService{achievementRepo: achievementRepo, logger: logger}
}

func (s *achievementService) Create(code, name, conditionType string, threshold int, icon, description string) (*model.Achievement, error) {
	if code == "" || name == "" || threshold <= 0 {
		return nil, fmt.Errorf("create achievement: %w", util.ErrValidation)
	}
	achievement := &model.Achievement{Code: code, Name: name, ConditionType: conditionType, Threshold: threshold, Icon: icon, Description: description}
	if err := s.achievementRepo.Create(achievement); err != nil {
		return nil, fmt.Errorf("create achievement: %w", err)
	}
	s.logger.Info(constants.LogAchievementCreated, "achievement_id", achievement.ID, "code", code)
	return achievement, nil
}

func (s *achievementService) List() ([]model.Achievement, error) {
	list, err := s.achievementRepo.List()
	if err != nil {
		return nil, fmt.Errorf("list achievements: %w", err)
	}
	s.logger.Info(constants.LogAchievementList, "count", len(list))
	return list, nil
}

func (s *achievementService) MyAchievements(userID uint) ([]model.UserAchievement, error) {
	list, err := s.achievementRepo.ListUserAchievements(userID)
	if err != nil {
		return nil, fmt.Errorf("my achievements user[%d]: %w", userID, err)
	}
	return list, nil
}

// CheckAndGrant 根据学习时长达成判定并发放徽章。
func (s *achievementService) CheckAndGrant(userID uint, studyMinutes int) error {
	achievements, err := s.achievementRepo.List()
	if err != nil {
		return fmt.Errorf("check and grant: %w", err)
	}
	for _, a := range achievements {
		ok := false
		if a.ConditionType == "study_hours" && studyMinutes >= a.Threshold*60 {
			ok = true
		}
		if a.ConditionType == "streak_days" && studyMinutes >= a.Threshold {
			ok = true
		}
		if !ok {
			continue
		}
		exists, err := s.achievementRepo.Exists(userID, a.ID)
		if err != nil {
			return fmt.Errorf("check and grant exists: %w", err)
		}
		if exists {
			continue
		}
		ua := &model.UserAchievement{UserID: userID, AchievementID: a.ID}
		if err := s.achievementRepo.CreateUserAchievement(ua); err != nil {
			return fmt.Errorf("check and grant create: %w", err)
		}
		s.logger.Info(constants.LogAchievementGranted, "user_id", userID, "achievement", a.Code)
	}
	return nil
}
