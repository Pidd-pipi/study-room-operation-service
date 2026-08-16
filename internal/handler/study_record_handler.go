package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/middleware"
	"github.com/ld/studyroom/internal/service"
	"github.com/ld/studyroom/internal/util"
)

// StudyRecordHandler 学习时长接口处理器。
type StudyRecordHandler struct {
	studySvc service.StudyRecordService
	achSvc   service.AchievementService
}

// NewStudyRecordHandler 构造学习时长处理器。
func NewStudyRecordHandler(studySvc service.StudyRecordService, achSvc service.AchievementService) *StudyRecordHandler {
	return &StudyRecordHandler{studySvc: studySvc, achSvc: achSvc}
}

// Ranking 学习排行榜（day/week/month）。
func (h *StudyRecordHandler) Ranking(c *gin.Context) {
	period := c.Query("period")
	if period == "" {
		period = "month"
	}
	limit := 20
	ranking, err := h.studySvc.Ranking(period, limit)
	if err != nil {
		c.Error(fmt.Errorf("handler ranking: %w", err))
		return
	}
	util.OK(c, ranking)
}

// Stats 我的学习统计。
func (h *StudyRecordHandler) Stats(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	stats, err := h.studySvc.Stats(claims.UserID)
	if err != nil {
		c.Error(fmt.Errorf("handler study stats: %w", err))
		return
	}
	util.OK(c, stats)
}

// Achievements 我的成就徽章。
func (h *StudyRecordHandler) Achievements(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	list, err := h.achSvc.MyAchievements(claims.UserID)
	if err != nil {
		c.Error(fmt.Errorf("handler my achievements: %w", err))
		return
	}
	util.OK(c, list)
}

// AchievementList 成就定义列表。
func (h *StudyRecordHandler) AchievementList(c *gin.Context) {
	list, err := h.achSvc.List()
	if err != nil {
		c.Error(fmt.Errorf("handler achievement list: %w", err))
		return
	}
	util.OK(c, list)
}
