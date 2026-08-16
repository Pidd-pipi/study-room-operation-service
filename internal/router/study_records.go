package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/handler"
)

func registerStudyRoutes(v1 *gin.RouterGroup, h *handler.StudyRecordHandler, auth gin.HandlerFunc, adminRoles []constants.UserRole) {
	study := v1.Group("/study", auth)
	study.GET("/ranking", h.Ranking)
	study.GET("/me/stats", h.Stats)
	study.GET("/me/achievements", h.Achievements)
	study.GET("/achievements", h.AchievementList)
}
