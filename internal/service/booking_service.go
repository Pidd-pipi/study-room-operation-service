package service

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/repository"
	"github.com/ld/studyroom/internal/util"
)

// BookingService 预约业务逻辑。
type BookingService interface {
	Create(userID, seatID uint, bookingDate string, startHour, endHour int) (*model.Booking, error)
	List(userID uint, page, pageSize int, status constants.BookingStatus) ([]model.Booking, int64, error)
	CheckIn(qrCode string) (*model.Booking, error)
	CheckOut(bookingID uint) (*model.Booking, error)
	Cancel(userID, bookingID uint) (*model.Booking, error)
	MarkNoShow(bookingID uint) (*model.Booking, error)
	CountConflict(seatID uint, date string, startHour, endHour int) (int64, error)
}

type bookingService struct {
	bookingRepo   repository.BookingRepository
	seatRepo      repository.SeatRepository
	violationRepo repository.ViolationRepository
	studyRepo     repository.StudyRecordRepository
	userSvc       UserService
	db            *gorm.DB
	logger        *slog.Logger
}

// NewBookingService 构造预约服务。
func NewBookingService(bookingRepo repository.BookingRepository, seatRepo repository.SeatRepository, violationRepo repository.ViolationRepository, studyRepo repository.StudyRecordRepository, userSvc UserService, db *gorm.DB, logger *slog.Logger) BookingService {
	return &bookingService{bookingRepo: bookingRepo, seatRepo: seatRepo, violationRepo: violationRepo, studyRepo: studyRepo, userSvc: userSvc, db: db, logger: logger}
}

func (s *bookingService) Create(userID, seatID uint, bookingDate string, startHour, endHour int) (*model.Booking, error) {
	if bookingDate == "" || startHour < 0 || endHour <= startHour || endHour > 24 {
		return nil, fmt.Errorf("create booking time[%s %d-%d]: %w", bookingDate, startHour, endHour, util.ErrValidation)
	}
	if blocked, err := s.violationRepo.ActiveBlockedUntil(userID); err != nil {
		return nil, fmt.Errorf("create booking check blocked user[%d]: %w", userID, err)
	} else if blocked != nil {
		return nil, fmt.Errorf("create booking user[%d] blocked until %s: %w", userID, blocked.Format("2006-01-02 15:04"), util.ErrBlocked)
	}
	seat, err := s.seatRepo.FindByID(seatID)
	if err != nil {
		return nil, fmt.Errorf("create booking seat[id=%d]: %w", seatID, err)
	}
	if seat.Status == constants.SeatUnavailable {
		return nil, fmt.Errorf("create booking seat[id=%d] unavailable: %w", seatID, util.ErrConflict)
	}

	booking := &model.Booking{
		BookingNo: generateBookingNo(), UserID: userID, SeatID: seatID,
		BookingDate: bookingDate, StartHour: startHour, EndHour: endHour,
		Status: constants.BookingPending, QRCode: generateQRCode(),
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		conflict, countErr := s.bookingRepo.CountConflictTx(tx, seatID, bookingDate, startHour, endHour)
		if countErr != nil {
			return fmt.Errorf("create booking conflict check: %w", countErr)
		}
		if conflict > 0 {
			s.logger.Warn(constants.LogBookingConflict, "seat_id", seatID, "date", bookingDate, "start", startHour, "end", endHour)
			return fmt.Errorf("create booking seat[%d] date[%s]: %w", seatID, bookingDate, util.ErrSeatConflict)
		}
		if createErr := s.bookingRepo.CreateTx(tx, booking); createErr != nil {
			return fmt.Errorf("create booking: %w", createErr)
		}
		if statusErr := s.seatRepo.UpdateStatusTx(tx, seatID, constants.SeatBooked); statusErr != nil {
			return fmt.Errorf("create booking update seat: %w", statusErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogBookingCreateSuccess, "booking_id", booking.ID, "user", userID, "seat", seatID)
	return booking, nil
}

func (s *bookingService) List(userID uint, page, pageSize int, status constants.BookingStatus) ([]model.Booking, int64, error) {
	list, total, err := s.bookingRepo.List(userID, page, pageSize, status)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}
	s.logger.Info(constants.LogBookingListQueried, "user", userID, "total", total)
	return list, total, nil
}

func (s *bookingService) CheckIn(qrCode string) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByCode(qrCode)
	if err != nil {
		return nil, fmt.Errorf("check in booking code[%s]: %w", qrCode, err)
	}
	if booking.Status != constants.BookingPending {
		s.logger.Warn(constants.LogBookingCheckInFailed, "booking_id", booking.ID, "status", booking.Status)
		return nil, fmt.Errorf("check in booking[id=%d] status[%s]: %w", booking.ID, booking.Status, util.ErrConflict)
	}
	now := time.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if transitionErr := s.bookingRepo.TransitionStatusTx(tx, booking.ID, booking.Status, constants.BookingCheckedIn); transitionErr != nil {
			return fmt.Errorf("check in booking[id=%d] status[%s]: %w", booking.ID, booking.Status, transitionErr)
		}
		booking.Status = constants.BookingCheckedIn
		booking.CheckInTime = &now
		if updateErr := s.bookingRepo.UpdateTx(tx, booking); updateErr != nil {
			return fmt.Errorf("check in booking[id=%d]: %w", booking.ID, updateErr)
		}
		if statusErr := s.seatRepo.UpdateStatusTx(tx, booking.SeatID, constants.SeatInUse); statusErr != nil {
			return fmt.Errorf("check in update seat: %w", statusErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogBookingCheckInSuccess, "booking_id", booking.ID)
	return booking, nil
}

func (s *bookingService) CheckOut(bookingID uint) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, fmt.Errorf("check out booking[id=%d]: %w", bookingID, err)
	}
	if booking.Status != constants.BookingCheckedIn {
		return nil, fmt.Errorf("check out booking[id=%d] status[%s]: %w", bookingID, booking.Status, util.ErrConflict)
	}
	now := time.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if transitionErr := s.bookingRepo.TransitionStatusTx(tx, booking.ID, booking.Status, constants.BookingCompleted); transitionErr != nil {
			return fmt.Errorf("check out booking[id=%d] status[%s]: %w", bookingID, booking.Status, transitionErr)
		}
		booking.CheckOutTime = &now
		booking.Status = constants.BookingCompleted
		if updateErr := s.bookingRepo.UpdateTx(tx, booking); updateErr != nil {
			return fmt.Errorf("check out booking[id=%d]: %w", bookingID, updateErr)
		}
		duration := 0
		if booking.CheckInTime != nil {
			duration = int(now.Sub(*booking.CheckInTime).Minutes())
			if duration < 1 {
				duration = 1
			}
		}
		record := &model.StudyRecord{UserID: booking.UserID, BookingID: booking.ID, StudyDate: booking.BookingDate, DurationMinutes: duration}
		if createErr := s.studyRepo.CreateTx(tx, record); createErr != nil {
			return fmt.Errorf("check out create study record: %w", createErr)
		}
		if addErr := s.userSvc.AddStudyMinutesTx(tx, booking.UserID, duration); addErr != nil {
			return fmt.Errorf("check out add study minutes: %w", addErr)
		}
		if statusErr := s.seatRepo.UpdateStatusTx(tx, booking.SeatID, constants.SeatIdle); statusErr != nil {
			return fmt.Errorf("check out update seat: %w", statusErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogBookingComplete, "booking_id", booking.ID, "duration", int(now.Sub(*booking.CheckInTime).Minutes()))
	return booking, nil
}

func (s *bookingService) Cancel(userID, bookingID uint) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, fmt.Errorf("cancel booking[id=%d]: %w", bookingID, err)
	}
	if booking.UserID != userID {
		return nil, fmt.Errorf("cancel booking[id=%d] not owner: %w", bookingID, util.ErrForbidden)
	}
	if booking.Status == constants.BookingCompleted {
		return nil, fmt.Errorf("cancel booking[id=%d] status[%s]: %w", bookingID, booking.Status, util.ErrConflict)
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if transitionErr := s.bookingRepo.TransitionStatusTx(tx, booking.ID, booking.Status, constants.BookingCancelled); transitionErr != nil {
			return fmt.Errorf("cancel booking[id=%d] status[%s]: %w", bookingID, booking.Status, transitionErr)
		}
		if statusErr := s.seatRepo.UpdateStatusTx(tx, booking.SeatID, constants.SeatIdle); statusErr != nil {
			return fmt.Errorf("cancel booking update seat: %w", statusErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	booking.Status = constants.BookingCancelled
	s.logger.Info(constants.LogBookingCancelSuccess, "booking_id", bookingID)
	return booking, nil
}

func (s *bookingService) MarkNoShow(bookingID uint) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return nil, fmt.Errorf("mark no show booking[id=%d]: %w", bookingID, err)
	}
	if !constants.CanBookingTransition(booking.Status, constants.BookingNoShow) {
		return nil, fmt.Errorf("mark no show booking[id=%d] status[%s]: %w", bookingID, booking.Status, util.ErrConflict)
	}
	blockedUntil := time.Now().Add(72 * time.Hour)
	violation := &model.ViolationRecord{
		UserID: booking.UserID, BookingID: booking.ID, ViolationType: "no_show",
		Points: 3, BlockedUntil: &blockedUntil, Remark: "爽约未签到，限制预约 72 小时",
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if transitionErr := s.bookingRepo.TransitionStatusTx(tx, booking.ID, booking.Status, constants.BookingNoShow); transitionErr != nil {
			return fmt.Errorf("mark no show booking[id=%d] status[%s]: %w", bookingID, booking.Status, transitionErr)
		}
		if createErr := s.violationRepo.CreateTx(tx, violation); createErr != nil {
			return fmt.Errorf("mark no show create violation: %w", createErr)
		}
		if statusErr := s.seatRepo.UpdateStatusTx(tx, booking.SeatID, constants.SeatIdle); statusErr != nil {
			return fmt.Errorf("mark no show update seat: %w", statusErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	booking.Status = constants.BookingNoShow
	s.logger.Info(constants.LogBookingNoShow, "booking_id", bookingID, "user", booking.UserID)
	return booking, nil
}

func (s *bookingService) CountConflict(seatID uint, date string, startHour, endHour int) (int64, error) {
	return s.bookingRepo.CountConflict(seatID, date, startHour, endHour)
}

func generateBookingNo() string {
	return fmt.Sprintf("BK%s%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%10000)
}

func generateQRCode() string {
	return fmt.Sprintf("QR-%s-%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%1000000)
}
