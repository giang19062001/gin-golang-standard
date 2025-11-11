package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/services"
	"github.com/giang19062001/gin-golang-standard/internal/utils"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	service services.IEventService
}

func NewEventController(service services.IEventService) *EventController {
	return &EventController{service: service}
}

// @Summary Tạo event mới
// @Description Thêm event mới vào database
// @Tags Events
// @Accept json
// @Produce json
// @Param event body dto.CreateEventDto true "Thông tin event vừa thêm"
// @Success 201 {object} models.Event
// @Router /events [post]
func (ctl *EventController) CreateEvent(c *gin.Context) {
	newEvent := &dto.CreateEventDto{}
	if err := c.ShouldBindJSON(newEvent); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return

	}
	event, err := ctl.service.CreateEvent(newEvent)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// @Summary Cập nhập event
// @Description  Cập nhập event đã tồn tại bằng ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Id của event"
// @Param event body dto.UpdateEventDto true "Dữ liệu event sau cập nhập"
// @Success 200 {object} models.Event
// @Router /events/{id} [put]
func (ctl *EventController) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "id event không hợp lệ",
		})
		return
	}

	changesEvent := &dto.UpdateEventDto{}
	if err := c.ShouldBindJSON(changesEvent); err != nil {
		c.JSON(http.StatusBadGateway, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	event, err := ctl.service.UpdateEvent(id, changesEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, event)
}

// @Summary Trả toàn bộ events
// @Description Lấy danh sách tất cả events
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {array} models.Event
// @Router /events [get]
func (ctl *EventController) GetAllEvents(c *gin.Context) {
	events, err := ctl.service.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, events)
}

// @Summary Lấy event theo ID
// @Description Trích xuất event cụ thể bằng ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Id của event"
// @Success 200 {object} models.Event
// @Router /events/{id} [get]
func (ctl *EventController) GetEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "id event không hợp lệ",
		})
		return
	}
	event, err := ctl.service.GetEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, event)
}

// @Summary Xóa event
// @Description Xóa event theo Id
// @Tags Events
// @Param id path int true "Id của event"
// @Success 200 {object} dto.MessageResponse
// @Router /events/{id} [delete]
func (ctl *EventController) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "id event không hợp lệ",
		})
		return
	}

	msg, err := ctl.service.DeleteEvent(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"message": msg,
	})
}

// @Summary      Upload ảnh cho event **MULTIPLE FILE**
// @Description  Upload mỗi file <= 10MB. Nếu có 1 file sai định dạng → KHÔNG LƯU BẤT KÌ FILE NÀO
// @Tags         Events
// @Accept       multipart/form-data
// @Produce      json
// @Param        eventId  formData  int   true   "ID của event"
// @Param        files  formData  file  true  "Danh sách file ảnh"
// @Success      200     {object}  dto.DataResponse{data=map[string]interface{}}
// @Router       /events/images [post]
func (ctl *EventController) UploadImages(c *gin.Context) {
	// lấy event
	eventIdStr := c.PostForm("eventId")
	if eventIdStr == "" {
		c.JSON(400, gin.H{"error": "eventId không thể trổng"})
		return
	}

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "eventId phải là số"})
		return
	}

	// lấy danh sách file ảnh
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Form không hợp lệ"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Danh sách files upload bị rỗng"})
		return
	}

	// validate các files đầu vào
	var validFiles []*multipart.FileHeader
	for _, file := range files {
		// kiểm tra .ext của file có hợp lệ hay không
		if err := utils.ValidateFile(file); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("File %s: %s", file.Filename, err.Error()),
			})
			return // stop -> trả lỗi về response
		}
		// nếu không lỗi -> đẩy vào mảng các files hợp lệ
		validFiles = append(validFiles, file)
	}

	// các files ok -> lưu
	filepaths, err := ctl.service.UploadImages(validFiles, eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.DataResponse{
		Data: map[string]interface{}{
			"filepaths": filepaths,
			"eventId":   eventId,
		},
	})
}
