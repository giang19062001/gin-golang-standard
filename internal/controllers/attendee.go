package controllers

import (
	"net/http"
	"strconv"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/services"
	"github.com/gin-gonic/gin"
)

type attendeeController struct {
	service services.IAttendeeService
}

func NewAttendeeController(service services.IAttendeeService) *attendeeController {
	return &attendeeController{service: service}
}

// @Summary Đăng kí tham dự sự kiên
// @Description Chèn dữ liệu tham dự của người dùng với sự kiện
// @Tags Attendees
// @Accept json
// @Produce json
// @Param event body dto.CreateAttendeeDto true "Thông tin đăng ký tham dự"
// @Success 201 {object} models.Attendee
// @Router /attendees/register [post]
func (ctl *attendeeController) RegisterAttendee(c *gin.Context) {
	dto := &dto.CreateAttendeeDto{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendee, err := ctl.service.RegisterAttendee(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, attendee)
}

// @Summary Xóa dự liệu tham dự sự kiện
// @Description Xóa dự liệu tham dự sự kiện của người tham dự
// @Tags Attendees
// @Param        eventId path      int     true  "Id của event"
// @Param        userId  path      int     true  "Id của user"
// @Success      200     {object}  dto.MessageResponse
// @Router       /attendees/{eventId}/{userId} [delete]
func (ctl *attendeeController) DeleteAttendee(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id event không hợp lệ",
		})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id user không hợp lệ",
		})
		return
	}

	err = ctl.service.Delete(userId, eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Xóa dữ liệu tham dự thành công",
	})
}

// @Summary Trả toạn bộ sự kiện mà người tham dự đã tham gia
// @Description Lấy danh sách tất cả sự kiện đã tham gia của người tham dự này
// @Tags Attendees
// @Accept json
// @Produce json
// @Param  userId path int true  "Id của user"
// @Success 200 {array} models.Event
// @Router /attendees/events/{userId} [get]
func (ctl *attendeeController) GetEventsByUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id user không hợp lệ",
		})
		return
	}

	event, err := ctl.service.GetEventsByUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, event)
}

// @Summary Trả toàn bộ người tham dự sự kiện
// @Description Lấy danh sách tất cả người tham dự sự kiện này
// @Tags Attendees
// @Accept json
// @Produce json
// @Param eventId path int true "Id của event"
// @Success 200 {array} models.User
// @Router /attendees/users/{eventId} [get]
func (ctl *attendeeController) GetUsersByEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id event không hợp lệ",
		})
		return
	}

	event, err := ctl.service.GetUsersByEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, event)
}
