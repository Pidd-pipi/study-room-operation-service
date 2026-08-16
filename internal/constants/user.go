package constants

import "errors"

// UserRole 用户角色枚举，前后端共享定义（frontend/src/constants/user.ts 对应实现）。
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// ErrInvalidUserRole 非法角色。
var ErrInvalidUserRole = errors.New("invalid user role")

func (r UserRole) Valid() bool {
	switch r {
	case RoleUser, RoleAdmin:
		return true
	}
	return false
}
