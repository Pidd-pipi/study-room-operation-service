package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/handler"
	"github.com/ld/studyroom/internal/middleware"
)

func registerBookingRoutes(v1 *gin.RouterGroup, h *handler.BookingHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole, limiter gin.HandlerFunc) {
	bookings := v1.Group("/bookings", auth)
	bookings.GET("", h.List)
	bookings.POST("", limiter, h.Create)
	bookings.POST("/check-in", limiter, h.CheckIn)
	bookings.POST("/:id/check-out", h.CheckOut)
	bookings.POST("/:id/cancel", h.Cancel)
	bookings.POST("/:id/no-show", middleware.RequireRole(adminRoles...), h.MarkNoShow)
}
