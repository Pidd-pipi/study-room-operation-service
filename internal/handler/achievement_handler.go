package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/dto"
	"github.com/ld/studyroom/internal/service"
	"github.com/ld/studyroom/internal/util"
)

// AchievementHandler 成就徽章管理接口处理器。
type AchievementHandler struct {
	achSvc service.AchievementService
}

// NewAchievementHandler 构造成就徽章管理处理器。
func NewAchievementHandler(achSvc service.AchievementService) *AchievementHandler {
	return &AchievementHandler{achSvc: achSvc}
}

func (h *AchievementHandler) Create(c *gin.Context) {
	var req dto.AchievementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	achievement, err := h.achSvc.Create(req.Code, req.Name, req.ConditionType, req.Threshold, req.Icon, req.Description)
	if err != nil {
		c.Error(fmt.Errorf("handler create achievement: %w", err))
		return
	}
	util.OK(c, achievement)
}
