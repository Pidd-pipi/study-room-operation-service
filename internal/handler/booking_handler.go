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

// BookingHandler 预约接口处理器。
type BookingHandler struct {
	bookingSvc service.BookingService
}

// NewBookingHandler 构造预约处理器。
func NewBookingHandler(bookingSvc service.BookingService) *BookingHandler {
	return &BookingHandler{bookingSvc: bookingSvc}
}

func (h *BookingHandler) Create(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	var req dto.BookingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	booking, err := h.bookingSvc.Create(claims.UserID, req.SeatID, req.BookingDate, req.StartHour, req.EndHour)
	if err != nil {
		c.Error(fmt.Errorf("handler create booking: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgBookingCreateSuccess, booking)
}

func (h *BookingHandler) List(c *gin.Context) {
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
	status := constants.BookingStatus(c.Query("status"))
	list, total, err := h.bookingSvc.List(claims.UserID, q.Page, q.PageSize, status)
	if err != nil {
		c.Error(fmt.Errorf("handler list bookings: %w", err))
		return
	}
	util.OK(c, gin.H{"list": list, "total": total, "page": q.Page, "page_size": q.PageSize})
}

func (h *BookingHandler) CheckIn(c *gin.Context) {
	var req dto.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.Validation(constants.MsgInvalidRequest, err))
		return
	}
	booking, err := h.bookingSvc.CheckIn(req.QRCode)
	if err != nil {
		c.Error(fmt.Errorf("handler check in: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgCheckInSuccess, booking)
}

func (h *BookingHandler) CheckOut(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的预约ID", err))
		return
	}
	booking, err := h.bookingSvc.CheckOut(uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler check out: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgCheckOutSuccess, booking)
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	claims, err := middleware.CurrentUser(c)
	if err != nil {
		c.Error(util.Unauthorized(constants.MsgUnauthorized, err))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的预约ID", err))
		return
	}
	booking, err := h.bookingSvc.Cancel(claims.UserID, uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler cancel booking: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgBookingCancelSuccess, booking)
}

func (h *BookingHandler) MarkNoShow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.Validation("无效的预约ID", err))
		return
	}
	booking, err := h.bookingSvc.MarkNoShow(uint(id))
	if err != nil {
		c.Error(fmt.Errorf("handler mark no show: %w", err))
		return
	}
	util.OKMessage(c, constants.MsgBookingNoShow, booking)
}
