package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/repository"
	"github.com/ld/studyroom/internal/util"
)

// SeatService 座位业务逻辑。
type SeatService interface {
	Create(seatNo string, floor int, zone constants.SeatZone, seatType string, x, y int, remark string) (*model.Seat, error)
	Update(id uint, seatNo string, floor int, zone constants.SeatZone, seatType string, x, y int, remark string) (*model.Seat, error)
	Delete(id uint) error
	List(floor int, zone constants.SeatZone) ([]model.Seat, error)
	ChangeStatus(id uint, status constants.SeatStatus) (*model.Seat, error)
	GetByID(id uint) (*model.Seat, error)
}

type seatService struct {
	seatRepo repository.SeatRepository
	logger   *slog.Logger
}

// NewSeatService 构造座位服务。
func NewSeatService(seatRepo repository.SeatRepository, logger *slog.Logger) SeatService {
	return &seatService{seatRepo: seatRepo, logger: logger}
}

func (s *seatService) Create(seatNo string, floor int, zone constants.SeatZone, seatType string, x, y int, remark string) (*model.Seat, error) {
	if seatNo == "" {
		return nil, fmt.Errorf("create seat: %w", util.ErrValidation)
	}
	if floor <= 0 {
		floor = 1
	}
	seat := &model.Seat{SeatNo: seatNo, Floor: floor, Zone: zone, SeatType: seatType, Status: constants.SeatIdle, X: x, Y: y, Remark: remark}
	if err := s.seatRepo.Create(seat); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, fmt.Errorf("create seat[no=%s]: %w", seatNo, util.ErrConflict)
		}
		return nil, fmt.Errorf("create seat[no=%s]: %w", seatNo, err)
	}
	s.logger.Info(constants.LogSeatCreated, "seat_id", seat.ID, "seat_no", seatNo)
	return seat, nil
}

func (s *seatService) Update(id uint, seatNo string, floor int, zone constants.SeatZone, seatType string, x, y int, remark string) (*model.Seat, error) {
	seat, err := s.seatRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update seat[id=%d]: %w", id, err)
	}
	if seatNo != "" {
		seat.SeatNo = seatNo
	}
	if floor > 0 {
		seat.Floor = floor
	}
	if zone.Valid() {
		seat.Zone = zone
	}
	if seatType != "" {
		seat.SeatType = seatType
	}
	seat.X = x
	seat.Y = y
	if remark != "" {
		seat.Remark = remark
	}
	if err := s.seatRepo.Update(seat); err != nil {
		return nil, fmt.Errorf("update seat[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogSeatUpdated, "seat_id", id)
	return seat, nil
}

func (s *seatService) Delete(id uint) error {
	if err := s.seatRepo.Delete(id); err != nil {
		return fmt.Errorf("delete seat[id=%d]: %w", id, err)
	}
	s.logger.Info(constants.LogSeatDeleted, "seat_id", id)
	return nil
}

func (s *seatService) List(floor int, zone constants.SeatZone) ([]model.Seat, error) {
	seats, err := s.seatRepo.List(floor, zone)
	if err != nil {
		return nil, fmt.Errorf("list seats: %w", err)
	}
	s.logger.Info(constants.LogSeatListQueried, "count", len(seats))
	return seats, nil
}

func (s *seatService) ChangeStatus(id uint, status constants.SeatStatus) (*model.Seat, error) {
	if !status.Valid() {
		return nil, fmt.Errorf("change seat status[%s]: %w", status, util.ErrValidation)
	}
	seat, err := s.seatRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("change seat status[id=%d]: %w", id, err)
	}
	if err := s.seatRepo.UpdateStatus(id, status); err != nil {
		return nil, fmt.Errorf("change seat status[id=%d]: %w", id, err)
	}
	seat.Status = status
	s.logger.Info(constants.LogSeatStatusChanged, "seat_id", id, "status", status)
	return seat, nil
}

func (s *seatService) GetByID(id uint) (*model.Seat, error) {
	return s.seatRepo.FindByID(id)
}
