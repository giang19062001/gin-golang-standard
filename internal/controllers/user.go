package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/services"
	"github.com/giang19062001/gin-golang-standard/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type UserController struct {
	service services.IUserService
}

func NewUserController(service services.IUserService) *UserController {
	return &UserController{service: service}
}

// @Summary      Đăng nhập
// @Description  Đăng nhập tài khoản người dùng mới với email, mật khẩu
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.LoginInDto  true  "Thông tin người dùng"
// @Success      201   {object}  dto.TokenResponse
// @Router       /auth/login [post]
func (ctl *UserController) LoginUser(c *gin.Context) {
	var loginInDto dto.LoginInDto
	if err := c.ShouldBindJSON(&loginInDto); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	tokenString, err := ctl.service.LoginUser(&loginInDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.TokenResponse{Token: tokenString})

}

// @Summary      Đăng ký người dùng mới
// @Description  Tạo tài khoản người dùng mới với email, mật khẩu và tên
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.RegisterDto  true  "Thông tin người dùng"
// @Success      201   {object}  models.User    "Người dùng được tạo thành công"
// @Router       /auth/register [post]
func (ctl *UserController) RegisterUser(c *gin.Context) {
	var userDto *dto.RegisterDto

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	user, err := ctl.service.RegisterUser(userDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, user)

}

// @Summary Trả toàn bộ users
// @Description Lấy danh sách tất cả users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (ctl *UserController) GetAllUser(c *gin.Context) {
	users, err := ctl.service.GetAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Xuất toàn bộ users ra file Excel
// @Description Lấy danh sách tất cả users từ database và trả về file Excel
// @Tags Users
// @Produce application/octet-stream
// @Success 200 {file} file.xlsx
// @Router /users/export [get]
func (ctl *UserController) ExportUsersExcel(c *gin.Context) {
	users, err := ctl.service.GetAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// tạo file excel với users data
	f, err := ctl.service.ExportUsersExcel(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// lấy ngày hôm nay với định dạng YYYY-MM-DD
	currentDate := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("users-%s.xlsx", currentDate)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	// ghi file excel ra response
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể ghi file Excel: " + err.Error()})
		return
	}
}

// @Summary Import users từ file Excel
// @Description Nhận file Excel, đọc danh sách users và thêm vào database
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File Excel chứa danh sách users"
// @Success 200 {string} string "Import thành công"
// @Router /users/import [post]
func (ctl *UserController) ImportUsersExcel(c *gin.Context) {
	// Lấy file upload từ form-data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không có file upload"})
		return
	}

	// Mở file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mở file: " + err.Error()})
		return
	}
	defer f.Close()

	// Đọc Excel bằng excelize
	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc file Excel: " + err.Error()})
		return
	}

	// Lấy tất cả hàng trong sheet "Users"
	rows, err := excelFile.GetRows("Users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc sheet Users: " + err.Error()})
		return
	}

	if len(rows) == 0 || len(rows) == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Excel rỗng"})
		return
	}

	// Kiểm tra header bắt buộc
	header := rows[0]
	if len(header) < 3 || header[0] != "Name" || header[1] != "Email" || header[2] != "Password" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Excel bắt buộc phải có cột header theo thứ tự từ trái sang phải là Name, Email, Password"})
		return
	}

	// gọi hàm xử lý dữ liệu import và thêm vào database
	err = ctl.service.ImportUsersExcel(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Import users thành công"})
}

// @Summary Trả thông tin user đang đăng nhập  **Yêu cầu xác thực**
// @Description Lấy thông tin user đang đăng nhập
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users/profile [get]
// @Security BearerAuth
func (ctl *UserController) GetProfile(c *gin.Context) {
	claims, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"claims": claims,
	})
}

// @Summary      Upload avatar **FILE**
// @Description  Upload file avatar (image) cho user - max 10MB
// @Tags         Users
// @Accept       multipart/form-data
// @Produce      json
// @Param        userId  formData  int   true   "ID của user"
// @Param        avatar  formData  file  true   "File avatar để upload"
// @Success      200     {object}  dto.DataResponse{data=map[string]interface{}}
// @Router       /users/avatar [put]
func (ctl *UserController) UpdateAvatar(c *gin.Context) {
	// lấy user
	userIdStr := c.PostForm("userId")
	if userIdStr == "" {
		c.JSON(400, gin.H{"error": "userId không thể trổng"})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "userId phải là số"})
		return
	}

	// lấy avatar
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "file không thể trống",
		})
		return
	}
	defer file.Close()

	// kiểm tra định dạng file
	if err := utils.ValidateFile(header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// lưu
	avatar, err := ctl.service.UpdateAvatar(file, header, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.DataResponse{
		Data: map[string]interface{}{
			"userId": userId,
			"avatar": avatar,
		},
	})

}
