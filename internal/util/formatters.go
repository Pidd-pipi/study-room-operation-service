package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/ld/studyroom/internal/constants"
)

// formatters.go 集中日期、时长、状态文本、分区文本、座位状态文本格式化。
func FormatTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func FormatDate(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func FormatDuration(minutes int) string {
	if minutes < 60 {
		return fmt.Sprintf("%d 分钟", minutes)
	}
	return fmt.Sprintf("%d 小时 %d 分钟", minutes/60, minutes%60)
}

func BookingStatusText(s constants.BookingStatus) string {
	switch s {
	case constants.BookingPending:
		return "待签到"
	case constants.BookingCheckedIn:
		return "使用中"
	case constants.BookingCompleted:
		return "已完成"
	case constants.BookingCancelled:
		return "已取消"
	case constants.BookingNoShow:
		return "爽约"
	}
	return "未知"
}

func SeatStatusText(s constants.SeatStatus) string {
	switch s {
	case constants.SeatIdle:
		return "空闲"
	case constants.SeatBooked:
		return "已预约"
	case constants.SeatInUse:
		return "使用中"
	case constants.SeatUnavailable:
		return "维护中"
	}
	return "未知"
}

func SeatZoneText(z constants.SeatZone) string {
	switch z {
	case constants.ZoneSilent:
		return "静音区"
	case constants.ZoneDiscussion:
		return "讨论区"
	case constants.ZoneWindow:
		return "靠窗区"
	}
	return "未知"
}

func UserRoleText(r constants.UserRole) string {
	switch r {
	case constants.RoleUser:
		return "学员"
	case constants.RoleAdmin:
		return "管理员"
	}
	return "未知"
}

func Lower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
