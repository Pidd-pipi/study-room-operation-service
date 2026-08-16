package constants

// BookingStatus 预约状态枚举，前后端共享定义（frontend/src/constants/booking.ts 对应实现）。
type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingCheckedIn BookingStatus = "checked_in"
	BookingCompleted BookingStatus = "completed"
	BookingCancelled BookingStatus = "cancelled"
	BookingNoShow    BookingStatus = "no_show"
)

func (s BookingStatus) Valid() bool {
	switch s {
	case BookingPending, BookingCheckedIn, BookingCompleted, BookingCancelled, BookingNoShow:
		return true
	}
	return false
}

// BookingStatusFlow 预约状态机：新增状态值需同步前端 constants、按钮显隐、日志模板、错误码、formatters。
var BookingStatusFlow = map[BookingStatus][]BookingStatus{
	BookingPending:   {BookingCheckedIn},
	BookingCheckedIn: {BookingCompleted},
	BookingCompleted: {},
	BookingCancelled: {},
	BookingNoShow:    {},
}

func CanBookingTransition(from, to BookingStatus) bool {
	for _, next := range BookingStatusFlow[from] {
		if next == to {
			return true
		}
	}
	return false
}
