package util

import (
	"testing"

	"github.com/ld/studyroom/internal/constants"
)

func TestSeatAndBookingText(t *testing.T) {
	if !constants.SeatUnavailable.Valid() {
		t.Fatal("SeatUnavailable.Valid() should be true")
	}
	if got := SeatStatusText(constants.SeatIdle); got != "空闲" {
		t.Fatalf("SeatStatusText(idle) = %s, want 空闲", got)
	}
	if got := SeatZoneText(constants.ZoneWindow); got != "靠窗区" {
		t.Fatalf("SeatZoneText(window) = %s, want 靠窗区", got)
	}
	if got := BookingStatusText(constants.BookingPending); got != "待签到" {
		t.Fatalf("BookingStatusText(pending) = %s, want 待签到", got)
	}
}
