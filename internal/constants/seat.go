package constants

// SeatStatus 座位状态枚举，前后端共享定义（frontend/src/constants/seat.ts 对应实现）。
type SeatStatus string

const (
	SeatIdle        SeatStatus = "idle"
	SeatBooked      SeatStatus = "booked"
	SeatInUse       SeatStatus = "in_use"
	SeatUnavailable SeatStatus = "unavailable"
)

func (s SeatStatus) Valid() bool {
	switch s {
	case SeatIdle, SeatBooked, SeatInUse:
		return true
	}
	return false
}

// SeatZone 座位分区枚举。
type SeatZone string

const (
	ZoneSilent     SeatZone = "silent"
	ZoneDiscussion SeatZone = "discussion"
	ZoneWindow     SeatZone = "window"
)

func (z SeatZone) Valid() bool {
	switch z {
	case ZoneSilent, ZoneDiscussion, ZoneWindow:
		return true
	}
	return false
}
