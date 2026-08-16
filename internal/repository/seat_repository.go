package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/util"
)

// SeatRepository 座位仓储。
type SeatRepository interface {
	Create(seat *model.Seat) error
	FindByID(id uint) (*model.Seat, error)
	List(floor int, zone constants.SeatZone) ([]model.Seat, error)
	Update(seat *model.Seat) error
	Delete(id uint) error
	UpdateStatus(id uint, status constants.SeatStatus) error
	UpdateStatusTx(tx *gorm.DB, id uint, status constants.SeatStatus) error
}

type seatRepository struct {
	db *gorm.DB
}

// NewSeatRepository 构造座位仓储。
func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) Create(seat *model.Seat) error {
	if err := r.db.Create(seat).Error; err != nil {
		if isDuplicate(err) {
			return fmt.Errorf("create seat: %v", ErrDuplicate)
		}
		return fmt.Errorf("create seat: %w", err)
	}
	return nil
}

func (r *seatRepository) FindByID(id uint) (*model.Seat, error) {
	var s model.Seat
	err := r.db.First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find seat by id: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find seat by id: %w", err)
	}
	return &s, nil
}

func (r *seatRepository) List(floor int, zone constants.SeatZone) ([]model.Seat, error) {
	var seats []model.Seat
	q := r.db.Order("floor asc, id asc")
	if floor > 0 {
		q = q.Where("floor = ?", floor)
	}
	if zone != "" {
		q = q.Where("zone = ?", zone)
	}
	if err := q.Find(&seats).Error; err != nil {
		return nil, fmt.Errorf("list seats: %w", err)
	}
	return seats, nil
}

func (r *seatRepository) Update(seat *model.Seat) error {
	if err := r.db.Save(seat).Error; err != nil {
		return fmt.Errorf("update seat: %w", err)
	}
	return nil
}

func (r *seatRepository) Delete(id uint) error {
	res := r.db.Delete(&model.Seat{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete seat: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete seat: %w", util.ErrNotFound)
	}
	return nil
}

func (r *seatRepository) UpdateStatus(id uint, status constants.SeatStatus) error {
	return r.UpdateStatusTx(nil, id, status)
}

func (r *seatRepository) UpdateStatusTx(tx *gorm.DB, id uint, status constants.SeatStatus) error {
	res := dbOrTx(r.db, tx).Model(&model.Seat{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return fmt.Errorf("update seat status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("update seat status: %w", util.ErrNotFound)
	}
	return nil
}
