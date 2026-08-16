package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/dto"
	"github.com/ld/studyroom/internal/repository"
	"github.com/ld/studyroom/internal/util"
)

// AuditLogHandler 审计日志接口处理器。
type AuditLogHandler struct {
	logRepo repository.AuditLogRepository
}

// NewAuditLogHandler 构造审计日志处理器。
func NewAuditLogHandler(logRepo repository.AuditLogRepository) *AuditLogHandler {
	return &AuditLogHandler{logRepo: logRepo}
}

func (h *AuditLogHandler) List(c *gin.Context) {
	var q dto.ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	q.Normalize()
	list, total, err := h.logRepo.List(q.Page, q.PageSize)
	if err != nil {
		c.Error(fmt.Errorf("handler list audit logs: %w", err))
		return
	}
	util.OK(c, gin.H{"list": list, "total": total, "page": q.Page, "page_size": q.PageSize})
}
