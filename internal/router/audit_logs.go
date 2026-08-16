package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/handler"
	"github.com/ld/studyroom/internal/middleware"
)

func registerAuditLogRoutes(v1 *gin.RouterGroup, h *handler.AuditLogHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole) {
	logs := v1.Group("/audit-logs", auth)
	logs.GET("", middleware.RequireRole(adminRoles...), h.List)
}
