package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/dto"
	"github.com/ld/studyroom/internal/middleware"
	"github.com/ld/studyroom/internal/service"
	"github.com/ld/studyroom/internal/util"
)

// ViolationHandler 违约记录接口处理器。
type ViolationHandler struct {
	violationSvc service.ViolationService
}

// NewViolationHandler 构造违约记录处理器。
func NewViolationHandler(violationSvc service.ViolationService) *ViolationHandler {
	return &ViolationHandler{violationSvc: violationSvc}
}

// Create 违约登记（管理员）。
func (h *ViolationHandler) Create(c *gin.Context) {
	var req dto.ViolationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	record, err := h.violationSvc.Create(req.UserID, req.BookingID, req.ViolationType, req.Points, req.BlockedHours, req.Remark)
	if err != nil {
		c.Error(fmt.Errorf("handler create violation: %w", err))
		return
	}
	util.OK(c, record)
}

// My 我的违约记录。
func (h *ViolationHandler) My(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	list, total, err := h.violationSvc.List(claims.UserID, q.Page, q.PageSize)
	if err != nil {
		c.Error(fmt.Errorf("handler my violations: %w", err))
		return
	}
	util.OK(c, gin.H{"list": list, "total": total, "page": q.Page, "page_size": q.PageSize})
}

// BlockedStatus 黑名单状态。
func (h *ViolationHandler) BlockedStatus(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	blockedUntil, err := h.violationSvc.BlockedUntil(claims.UserID)
	if err != nil {
		c.Error(fmt.Errorf("handler blocked status: %w", err))
		return
	}
	util.OK(c, gin.H{"blocked": blockedUntil != nil, "blocked_until": blockedUntil})
}
