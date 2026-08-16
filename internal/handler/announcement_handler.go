package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/dto"
	"github.com/ld/studyroom/internal/middleware"
	"github.com/ld/studyroom/internal/service"
	"github.com/ld/studyroom/internal/util"
)

// AnnouncementHandler 公告接口处理器。
type AnnouncementHandler struct {
	announcementSvc service.AnnouncementService
}

// NewAnnouncementHandler 构造公告处理器。
func NewAnnouncementHandler(announcementSvc service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{announcementSvc: announcementSvc}
}

func (h *AnnouncementHandler) Create(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	var req dto.AnnouncementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	announcement, err := h.announcementSvc.Create(req.Title, req.Content, claims.UserID)
	if err != nil {
		c.Error(fmt.Errorf("handler create announcement: %w", err))
		return
	}
	util.OK(c, announcement)
}

func (h *AnnouncementHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的公告ID", err))
		return
	}
	var req dto.AnnouncementUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	announcement, err := h.announcementSvc.Update(uint(id), req.Title, req.Content, isActive)
	if err != nil {
		c.Error(fmt.Errorf("handler update announcement: %w", err))
		return
	}
	util.OK(c, announcement)
}

func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的公告ID", err))
		return
	}
	if err := h.announcementSvc.Delete(uint(id)); err != nil {
		c.Error(fmt.Errorf("handler delete announcement: %w", err))
		return
	}
	util.OKMessage(c, "公告已删除", nil)
}

func (h *AnnouncementHandler) List(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	activeOnly := c.Query("active") == "1"
	list, total, err := h.announcementSvc.List(q.Page, q.PageSize, activeOnly)
	if err != nil {
		c.Error(fmt.Errorf("handler list announcements: %w", err))
		return
	}
	util.OK(c, gin.H{"list": list, "total": total, "page": q.Page, "page_size": q.PageSize})
}
