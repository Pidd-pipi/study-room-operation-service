package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/handler"
	"github.com/ld/studyroom/internal/middleware"
)

func registerSeatRoutes(v1 *gin.RouterGroup, h *handler.SeatHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole) {
	seats := v1.Group("/seats", auth)
	seats.GET("", h.List)
	seats.POST("", middleware.RequireRole(adminRoles...), h.Create)
	seats.PUT("/:id", middleware.RequireRole(adminRoles...), h.Update)
	seats.DELETE("/:id", middleware.RequireRole(adminRoles...), h.Delete)
	seats.PUT("/:id/status", middleware.RequireRole(adminRoles...), h.ChangeStatus)
}
