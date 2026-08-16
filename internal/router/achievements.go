package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/handler"
	"github.com/ld/studyroom/internal/middleware"
)

func registerAchievementRoutes(v1 *gin.RouterGroup, h *handler.AchievementHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole) {
	achievements := v1.Group("/achievements", auth)
	achievements.POST("", middleware.RequireRole(adminRoles...), h.Create)
}
