package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/handler"
	"github.com/ld/studyroom/internal/middleware"
)

func registerViolationRoutes(v1 *gin.RouterGroup, h *handler.ViolationHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole) {
	violations := v1.Group("/violations", auth)
	violations.POST("", middleware.RequireRole(adminRoles...), h.Create)
	violations.GET("/me", h.My)
	violations.GET("/me/blocked", h.BlockedStatus)
}
