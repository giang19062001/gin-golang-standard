package services

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/models"
	"github.com/giang19062001/gin-golang-standard/internal/repositories"
	"github.com/giang19062001/gin-golang-standard/internal/utils"
)

type userService struct {
	repo      repositories.IUserRepository
	jwtSecret string
}

type IUserService interface {
	GetAllUser() ([]models.User, error)
	Get(id int) (*models.User, error)
	RegisterUser(*dto.RegisterDto) (*models.User, error)
	LoginUser(*dto.LoginInDto) (string, error)
	GetUserOfEvent(int) ([]models.User, error)
	UpdateAvatar(multipart.File, *multipart.FileHeader, int) (string, error)
}

func NewUserService(repo repositories.IUserRepository, jwtSecret string) IUserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

func (ser *userService) LoginUser(loginDto *dto.LoginInDto) (string, error) {
	// kiểm tra user
	existingUser, err := ser.repo.GetByEmail(loginDto.Email)
	if existingUser == nil {
		return "", errors.New("user ko tồn tại")
	}
	if err != nil {
		return "", errors.New("có lỗi xảy ra")
	}

	// kiểm tra mật khẩu
	err = utils.ComparePassword(existingUser.Password, loginDto.Password)
	if err != nil {
		return "", errors.New("sai mật khẩu")
	}

	// tạo token
	tokenStr, err := utils.GenerateToken(existingUser.Id, ser.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil

}

func (ser *userService) RegisterUser(userDto *dto.RegisterDto) (*models.User, error) {

	existingUser, _ := ser.repo.GetByEmail(userDto.Email)
	if existingUser != nil {
		return nil, errors.New("email này đã được user khác sử dụng - không thể đăng ký")
	}

	// hash mật khẩu
	hashedPassword, err := utils.HashPassword(userDto.Password)
	if err != nil {
		return nil, err
	}
	// gán cho struct hiện tại đề insert
	userDto.Password = hashedPassword
	user := models.User{
		Email:    userDto.Email,
		Name:     userDto.Name,
		Password: userDto.Password,
	}
	err = ser.repo.Insert(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (ser *userService) GetAllUser() ([]models.User, error) {
	users, err := ser.repo.GetAllUser()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ser *userService) Get(id int) (*models.User, error) {
	user, _ := ser.repo.Get(id)
	if user == nil {
		return nil, errors.New("user không tồn tại")
	}
	return user, nil
}

func (ser *userService) GetUserOfEvent(eventId int) ([]models.User, error) {
	users, err := ser.repo.GetUserOfEvent(eventId)
	if err != nil {
		return []models.User{}, errors.New("lỗi trích xuất event")
	}
	return users, nil
}

func (ser *userService) UpdateAvatar(file multipart.File, header *multipart.FileHeader, userId int) (string, error) {
	// kiểm tra user
	user, _ := ser.repo.Get(userId)
	log.Print(userId)
	if user == nil {
		return "", errors.New("user không tồn tại")
	}

	prefix := fmt.Sprintf("user-%d", userId)
	// Xóa file cũ
	err := utils.DeleteFilesByPrefix(prefix)
	if err != nil {
		return "", err
	}
	// lưu file vào thư mục
	avatar := utils.SaveAndGetURL(prefix, file, header)

	// lưu vào db
	err = ser.repo.UpdateAvatar(avatar, userId)
	if err != nil {
		return "", errors.New("lỗi trích xuất event")
	}
	return avatar, nil
}
