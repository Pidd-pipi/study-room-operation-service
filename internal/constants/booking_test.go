package constants

import "testing"

func TestCanBookingTransition(t *testing.T) {
	tests := []struct {
		name string
		from BookingStatus
		to   BookingStatus
		want bool
	}{
		{"pending to checked_in", BookingPending, BookingCheckedIn, true},
		{"pending to cancelled", BookingPending, BookingCancelled, true},
		{"pending to no_show", BookingPending, BookingNoShow, true},
		{"checked_in to completed", BookingCheckedIn, BookingCompleted, true},
		{"completed to cancelled", BookingCompleted, BookingCancelled, false},
		{"cancelled to checked_in", BookingCancelled, BookingCheckedIn, false},
	}
	for _, tt := range tests {
		if got := CanBookingTransition(tt.from, tt.to); got != tt.want {
			t.Errorf("CanBookingTransition(%s,%s) = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestSeatZoneValid(t *testing.T) {
	tests := []struct {
		zone SeatZone
		want bool
	}{
		{ZoneSilent, true},
		{ZoneDiscussion, true},
		{ZoneWindow, true},
		{"outside", false},
	}
	for _, tt := range tests {
		if got := tt.zone.Valid(); got != tt.want {
			t.Errorf("SeatZone(%s).Valid() = %v, want %v", tt.zone, got, tt.want)
		}
	}
}
