package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/util"
)

type mockSeatRepo struct {
	seats map[uint]*model.Seat
	seq   uint
}

func newMockSeatRepo() *mockSeatRepo { return &mockSeatRepo{seats: map[uint]*model.Seat{}} }

func (m *mockSeatRepo) Create(s *model.Seat) error {
	m.seq++
	s.ID = m.seq
	m.seats[s.ID] = s
	return nil
}
func (m *mockSeatRepo) FindByID(id uint) (*model.Seat, error) {
	if s, ok := m.seats[id]; ok {
		return s, nil
	}
	return nil, util.ErrNotFound
}
func (m *mockSeatRepo) List(floor int, zone constants.SeatZone) ([]model.Seat, error) {
	return nil, nil
}
func (m *mockSeatRepo) Update(s *model.Seat) error { m.seats[s.ID] = s; return nil }
func (m *mockSeatRepo) Delete(id uint) error {
	if _, ok := m.seats[id]; !ok {
		return util.ErrNotFound
	}
	delete(m.seats, id)
	return nil
}
func (m *mockSeatRepo) UpdateStatus(id uint, status constants.SeatStatus) error {
	return m.UpdateStatusTx(nil, id, status)
}
func (m *mockSeatRepo) UpdateStatusTx(tx *gorm.DB, id uint, status constants.SeatStatus) error {
	if _, ok := m.seats[id]; !ok {
		return util.ErrNotFound
	}
	m.seats[id].Status = status
	return nil
}

func newTestSeatService() SeatService {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return NewSeatService(newMockSeatRepo(), logger)
}

func TestSeatServiceCreate(t *testing.T) {
	svc := newTestSeatService()
	tests := []struct {
		name    string
		seatNo  string
		zone    constants.SeatZone
		wantErr bool
	}{
		{"valid", "A-01", constants.ZoneSilent, false},
		{"empty no", "", constants.ZoneSilent, true},
		{"invalid zone", "A-02", "bad", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seat, err := svc.Create(tt.seatNo, 1, tt.zone, "单人桌", 1, 1, "")
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if !errors.Is(err, util.ErrValidation) {
					t.Fatalf("expected ErrValidation, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if seat.Status != constants.SeatIdle {
				t.Errorf("status = %s, want idle", seat.Status)
			}
		})
	}
}

func TestSeatServiceChangeStatus(t *testing.T) {
	svc := newTestSeatService()
	seat, err := svc.Create("A-09", 1, constants.ZoneWindow, "单座", 1, 1, "")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	updated, err := svc.ChangeStatus(seat.ID, constants.SeatUnavailable)
	if err != nil {
		t.Fatalf("change status failed: %v", err)
	}
	if updated.Status != constants.SeatUnavailable {
		t.Errorf("status = %s, want unavailable", updated.Status)
	}
	if _, err := svc.ChangeStatus(seat.ID, "invalid"); !errors.Is(err, util.ErrValidation) {
		t.Errorf("expected validation error, got %v", err)
	}
}
