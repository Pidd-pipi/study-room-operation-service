package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/util"
)

// BookingRepository 预约仓储。
type BookingRepository interface {
	Create(booking *model.Booking) error
	CreateTx(tx *gorm.DB, booking *model.Booking) error
	FindByID(id uint) (*model.Booking, error)
	FindByCode(qrCode string) (*model.Booking, error)
	List(userID uint, page, pageSize int, status constants.BookingStatus) ([]model.Booking, int64, error)
	Update(booking *model.Booking) error
	UpdateTx(tx *gorm.DB, booking *model.Booking) error
	TransitionStatusTx(tx *gorm.DB, id uint, from, to constants.BookingStatus) error
	CountConflict(seatID uint, date string, startHour, endHour int) (int64, error)
	CountConflictTx(tx *gorm.DB, seatID uint, date string, startHour, endHour int) (int64, error)
}

type bookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository 构造预约仓储。
func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.CreateTx(nil, booking)
}

func (r *bookingRepository) CreateTx(tx *gorm.DB, booking *model.Booking) error {
	if err := dbOrTx(r.db, tx).Create(booking).Error; err != nil {
		return fmt.Errorf("create booking: %w", err)
	}
	return nil
}

func (r *bookingRepository) FindByID(id uint) (*model.Booking, error) {
	var b model.Booking
	err := r.db.Preload("User").Preload("Seat").First(&b, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find booking by id: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find booking by id: %w", err)
	}
	return &b, nil
}

func (r *bookingRepository) FindByCode(qrCode string) (*model.Booking, error) {
	var b model.Booking
	err := r.db.Preload("User").Preload("Seat").Where("qr_code = ?", qrCode).First(&b).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find booking by code: %w", util.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find booking by code: %w", err)
	}
	return &b, nil
}

func (r *bookingRepository) List(userID uint, page, pageSize int, status constants.BookingStatus) ([]model.Booking, int64, error) {
	var list []model.Booking
	var total int64
	q := r.db.Model(&model.Booking{})
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}
	if err := q.Preload("User").Preload("Seat").Offset(page * pageSize).Limit(pageSize).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}
	return list, total, nil
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.UpdateTx(nil, booking)
}

func (r *bookingRepository) UpdateTx(tx *gorm.DB, booking *model.Booking) error {
	if err := dbOrTx(r.db, tx).Save(booking).Error; err != nil {
		return fmt.Errorf("update booking: %w", err)
	}
	return nil
}

func (r *bookingRepository) TransitionStatusTx(tx *gorm.DB, id uint, from, to constants.BookingStatus) error {
	res := dbOrTx(r.db, tx).Model(&model.Booking{}).
		Where("id = ? AND status = ?", id, from).
		Update("status", to)
	if res.Error != nil {
		return fmt.Errorf("transition booking status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("transition booking status: %w", util.ErrConflict)
	}
	return nil
}

// CountConflict 检测同一座位同一日期时段冲突。
func (r *bookingRepository) CountConflict(seatID uint, date string, startHour, endHour int) (int64, error) {
	return r.CountConflictTx(nil, seatID, date, startHour, endHour)
}

func (r *bookingRepository) CountConflictTx(tx *gorm.DB, seatID uint, date string, startHour, endHour int) (int64, error) {
	var count int64
	err := dbOrTx(r.db, tx).Model(&model.Booking{}).
		Where("seat_id = ? AND booking_date = ? AND status IN ? AND start_hour < ? AND end_hour > ?",
			seatID, date, []string{string(constants.BookingPending), string(constants.BookingCheckedIn)}, endHour, startHour).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count booking conflict: %w", err)
	}
	return count, nil
}
