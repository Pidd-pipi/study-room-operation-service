package dto

import "github.com/ld/studyroom/internal/constants"

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Username string             `json:"username" binding:"required,min=3,max=64"`
	Password string             `json:"password" binding:"required,min=6,max=64"`
	Nickname string             `json:"nickname" binding:"max=64"`
	Phone    string             `json:"phone" binding:"max=32"`
	Role     constants.UserRole `json:"role"`
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应。
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// UpdateProfileRequest 修改资料请求。
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=64"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Phone    string `json:"phone" binding:"max=32"`
}

// ListQuery 分页参数。
type ListQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=200"`
}

func (q *ListQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
}

// SeatCreateRequest 座位创建请求。
type SeatCreateRequest struct {
	SeatNo   string             `json:"seat_no" binding:"required,max=32"`
	Floor    int                `json:"floor"`
	Zone     constants.SeatZone `json:"zone" binding:"required"`
	SeatType string             `json:"seat_type" binding:"max=32"`
	X        int                `json:"x"`
	Y        int                `json:"y"`
	Remark   string             `json:"remark" binding:"max=255"`
}

// SeatUpdateRequest 座位更新请求。
type SeatUpdateRequest struct {
	SeatNo   string             `json:"seat_no" binding:"max=32"`
	Floor    int                `json:"floor"`
	Zone     constants.SeatZone `json:"zone"`
	SeatType string             `json:"seat_type" binding:"max=32"`
	X        int                `json:"x"`
	Y        int                `json:"y"`
	Remark   string             `json:"remark" binding:"max=255"`
}

// SeatStatusRequest 座位状态变更请求。
type SeatStatusRequest struct {
	Status constants.SeatStatus `json:"status" binding:"required"`
}

// BookingCreateRequest 预约创建请求。
type BookingCreateRequest struct {
	SeatID      uint   `json:"seat_id" binding:"required"`
	BookingDate string `json:"booking_date" binding:"required"`
	StartHour   int    `json:"start_hour" binding:"min=0,max=23"`
	EndHour     int    `json:"end_hour" binding:"min=1,max=24"`
}

// CheckInRequest 签到请求。
type CheckInRequest struct {
	QRCode string `json:"qr_code" binding:"required"`
}

// AnnouncementCreateRequest 公告创建请求。
type AnnouncementCreateRequest struct {
	Title   string `json:"title" binding:"required,max=128"`
	Content string `json:"content" binding:"required"`
}

// AnnouncementUpdateRequest 公告更新请求。
type AnnouncementUpdateRequest struct {
	Title    string `json:"title" binding:"max=128"`
	Content  string `json:"content"`
	IsActive *bool  `json:"is_active"`
}

// AchievementCreateRequest 成就创建请求。
type AchievementCreateRequest struct {
	Code          string `json:"code" binding:"required,max=64"`
	Name          string `json:"name" binding:"required,max=64"`
	ConditionType string `json:"condition_type" binding:"required,max=32"`
	Threshold     int    `json:"threshold" binding:"required,min=1"`
	Icon          string `json:"icon" binding:"max=128"`
	Description   string `json:"description" binding:"max=255"`
}

// ViolationCreateRequest 违约登记请求。
type ViolationCreateRequest struct {
	UserID        uint   `json:"user_id" binding:"required"`
	BookingID     uint   `json:"booking_id"`
	ViolationType string `json:"violation_type" binding:"required,max=32"`
	Points        int    `json:"points"`
	BlockedHours  int    `json:"blocked_hours"`
	Remark        string `json:"remark" binding:"max=255"`
}
