package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/dto"
	"github.com/ld/studyroom/internal/service"
	"github.com/ld/studyroom/internal/util"
)

// SeatHandler 座位接口处理器。
type SeatHandler struct {
	seatSvc service.SeatService
}

// NewSeatHandler 构造座位处理器。
func NewSeatHandler(seatSvc service.SeatService) *SeatHandler {
	return &SeatHandler{seatSvc: seatSvc}
}

func (h *SeatHandler) Create(c *gin.Context) {
	var req dto.SeatCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	seat, err := h.seatSvc.Create(req.SeatNo, req.Floor, req.Zone, req.SeatType, req.X, req.Y, req.Remark)
	if err != nil {
		c.Error(fmt.Errorf("handler create seat: %w", err))
		return
	}
	util.OK(c, seat)
}

func (h *SeatHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的座位ID", err))
		return
	}
	var req dto.SeatUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	seat, err := h.seatSvc.Update(uint(id), req.SeatNo, req.Floor, req.Zone, req.SeatType, req.X, req.Y, req.Remark)
	if err != nil {
		c.Error(fmt.Errorf("handler update seat: %w", err))
		return
	}
	util.OK(c, seat)
}

func (h *SeatHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的座位ID", err))
		return
	}
	if err := h.seatSvc.Delete(uint(id)); err != nil {
		c.Error(fmt.Errorf("handler delete seat: %w", err))
		return
	}
	util.OKMessage(c, "座位已删除", nil)
}

func (h *SeatHandler) List(c *gin.Context) {
	floor, _ := strconv.Atoi(c.Query("floor"))
	zone := constants.SeatZone(c.Query("zone"))
	seats, err := h.seatSvc.List(floor, zone)
	if err != nil {
		c.Error(fmt.Errorf("handler list seats: %w", err))
		return
	}
	util.OK(c, seats)
}

func (h *SeatHandler) ChangeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的座位ID", err))
		return
	}
	var req dto.SeatStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	seat, err := h.seatSvc.ChangeStatus(uint(id), req.Status)
	if err != nil {
		c.Error(fmt.Errorf("handler change seat status: %w", err))
		return
	}
	util.OK(c, seat)
}
