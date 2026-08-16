package service

import (
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/repository"
	"github.com/ld/studyroom/internal/util"
)

// UserService 用户业务逻辑。
type UserService interface {
	Register(username, password, nickname, phone string, role constants.UserRole) (*model.User, error)
	Login(username, password string) (*model.User, string, error)
	UpdateProfile(id uint, nickname, avatar, phone string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	AddStudyMinutes(id uint, minutes int) error
	AddStudyMinutesTx(tx *gorm.DB, id uint, minutes int) error
	Ranking(limit int) ([]model.User, error)
}

type userService struct {
	userRepo  repository.UserRepository
	logger    *slog.Logger
	jwtSecret string
	ttlHours  int
}

// NewUserService 构造用户服务。
func NewUserService(userRepo repository.UserRepository, logger *slog.Logger, jwtSecret string, ttlHours int) UserService {
	return &userService{userRepo: userRepo, logger: logger, jwtSecret: jwtSecret, ttlHours: ttlHours}
}

func (s *userService) Register(username, password, nickname, phone string, role constants.UserRole) (*model.User, error) {
	if username == "" || len(password) < 6 {
		return nil, fmt.Errorf("register: %w", util.ErrValidation)
	}
	if !role.Valid() {
		role = constants.RoleUser
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user := &model.User{Username: username, PasswordHash: string(hash), Nickname: nickname, Phone: phone, Role: role}
	if err := s.userRepo.Create(user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, fmt.Errorf("register user[username=%s]: %w", username, util.ErrConflict)
		}
		return nil, fmt.Errorf("register user[username=%s]: %w", username, err)
	}
	s.logger.Info(constants.LogUserRegisterSuccess, "user_id", user.ID, "username", username)
	return user, nil
}

func (s *userService) Login(username, password string) (*model.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		s.logger.Warn(constants.LogUserLoginFailed, "username", username, "reason", "not found")
		return nil, "", fmt.Errorf("login: %w", util.ErrUnauthorized)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		s.logger.Warn(constants.LogUserLoginFailed, "username", username, "reason", "password mismatch")
		return nil, "", fmt.Errorf("login: %w", util.ErrUnauthorized)
	}
	token, err := util.GenerateToken(s.jwtSecret, s.ttlHours, user.ID, user.Username, user.Role, nil)
	if err != nil {
		return nil, "", fmt.Errorf("login generate token: %w", err)
	}
	s.logger.Info(constants.LogUserLoginSuccess, "user_id", user.ID, "username", username)
	return user, token, nil
}

func (s *userService) UpdateProfile(id uint, nickname, avatar, phone string) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update profile user[id=%d]: %w", id, err)
	}
	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if phone != "" {
		user.Phone = phone
	}
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("update profile user[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogUserProfileUpdated, "user_id", id)
	return user, nil
}

func (s *userService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) AddStudyMinutes(id uint, minutes int) error {
	return s.AddStudyMinutesTx(nil, id, minutes)
}

func (s *userService) AddStudyMinutesTx(tx *gorm.DB, id uint, minutes int) error {
	user, err := s.userRepo.FindByIDTx(tx, id)
	if err != nil {
		return fmt.Errorf("add study minutes user[id=%d]: %w", id, err)
	}
	user.TotalStudyMinutes += minutes
	if err := s.userRepo.UpdateTx(tx, user); err != nil {
		return fmt.Errorf("add study minutes user[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogStudyRecordCreated, "user_id", id, "minutes", minutes)
	return nil
}

func (s *userService) Ranking(limit int) ([]model.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	users, err := s.userRepo.ListRanking("total", limit)
	if err != nil {
		return nil, fmt.Errorf("ranking: %w", err)
	}
	s.logger.Info(constants.LogStudyRankingQueried, "limit", limit)
	return users, nil
}
