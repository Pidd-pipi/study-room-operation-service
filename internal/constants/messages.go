package constants

// messages.go 同时承载前端提示文案、后端返回文案与日志文案。
const (
	MsgOK                   = "ok"
	MsgInvalidRequest       = "请求参数不合法"
	MsgUnauthorized         = "未登录或登录已过期"
	MsgForbidden            = "无权限执行该操作"
	MsgNotFound             = "资源不存在"
	MsgConflict             = "资源状态冲突"
	MsgRateLimited          = "请求过于频繁，请稍后再试"
	MsgInternalError        = "服务器内部错误"
	MsgLoginSuccess         = "登录成功"
	MsgRegisterSuccess      = "注册成功"
	MsgPasswordIncorrect    = "用户名或密码错误"
	MsgSeatConflict         = "该座位在所选时段已被预约"
	MsgBookingNotFound      = "预约不存在"
	MsgBookingNotCheckIn    = "当前状态不可签到"
	MsgBookingNoShow        = "预约已标记为爽约"
	MsgViolationBlocked     = "您处于违约黑名单中，暂不能预约"
	MsgCheckInSuccess       = "签到成功，开始计时"
	MsgCheckOutSuccess      = "离场成功，学习时长已记录"
	MsgBookingCreateSuccess = "预约成功"
	MsgBookingCancelSuccess = "预约已取消"
)
