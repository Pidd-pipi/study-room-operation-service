package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/handler"
	"github.com/ld/studyroom/internal/middleware"
)

func registerAnnouncementRoutes(v1 *gin.RouterGroup, h *handler.AnnouncementHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole) {
	announcements := v1.Group("/announcements", auth)
	announcements.GET("", h.List)
	announcements.POST("", middleware.RequireRole(adminRoles...), h.Create)
	announcements.PUT("/:id", middleware.RequireRole(adminRoles...), h.Update)
	announcements.DELETE("/:id", middleware.RequireRole(adminRoles...), h.Delete)
}
