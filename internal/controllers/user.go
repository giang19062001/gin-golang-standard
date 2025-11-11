package controllers

import (
	"net/http"
	"strconv"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/services"
	"github.com/giang19062001/gin-golang-standard/internal/utils"
	"github.com/gin-gonic/gin"
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
