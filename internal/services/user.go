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
	"github.com/giang19062001/gin-golang-standard/pkg/logger"
	"github.com/xuri/excelize/v2"
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
	ExportUsersExcel([]models.User) (*excelize.File, error)
	ImportUsersExcel([][]string) error
}

func NewUserService(repo repositories.IUserRepository, jwtSecret string) IUserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

func (ser *userService) LoginUser(loginDto *dto.LoginInDto) (string, error) {
	logr := logger.With("userService")
	// kiểm tra user
	existingUser, err := ser.repo.GetByEmail(loginDto.Email)
	if existingUser == nil {
		logr.Error("user ko tồn tại")
		return "", errors.New("user ko tồn tại")
	}
	if err != nil {
		return "", errors.New("có lỗi xảy ra")
	}

	// kiểm tra mật khẩu
	err = utils.ComparePassword(existingUser.Password, loginDto.Password)
	if err != nil {
		logr.Error("user ko tồn tại")
		return "", errors.New("sai mật khẩu")
	}

	// tạo token
	tokenStr, err := utils.GenerateToken(existingUser.Id, ser.jwtSecret)
	if err != nil {
		return "", err
	}

	logr.Info("Đăng nhập thành công")
	return tokenStr, nil

}

func (ser *userService) RegisterUser(userDto *dto.RegisterDto) (*models.User, error) {
	logr := logger.With("userService")

	existingUser, _ := ser.repo.GetByEmail(userDto.Email)
	if existingUser != nil {
		logr.Error("email này đã được user khác sử dụng - không thể đăng ký")
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
func (ser *userService) ImportUsersExcel(rows [][]string) error {
	logr := logger.With("userService")

	// Parse dữ liệu thành slice User
	var users []models.User
	for i, row := range rows {
		if i == 0 {
			// bỏ qua header
			continue
		}
		if len(row) < 3 {
			continue
		}
		user := models.User{
			Name:     row[0],
			Email:    row[1],
			Password: row[2],
		}
		users = append(users, user)
	}

	// Log ra danh sách users
	for _, u := range users {
		log.Printf("User: Name=%s, Email=%s, Password=%s", u.Name, u.Email, u.Password)
		// kiểm tra nếu trong file import có bất kì email nào  trùng với email đã có trong database thì báo lỗi và dừng import
		existingUser, _ := ser.repo.GetByEmail(u.Email)
		if existingUser != nil {
			// email đã tồn tại -> ném lỗi
			errText := fmt.Sprintf("email '%s' này đã được user khác sử dụng - không thể thêm", u.Email)
			logr.Error(errText)
			return errors.New(errText)
		}
	}

	// nếu tất cả email đều hợp lệ thì insert nhiều
	usersHahed := []models.User{}
	for _, u := range users {
		// hash mật khẩu
		hashedPassword, _ := utils.HashPassword(u.Password)
		// gán cho struct hiện tại đề insert
		u.Password = hashedPassword
		userHash := models.User{
			Email:    u.Email,
			Name:     u.Name,
			Password: u.Password,
		}
		usersHahed = append(usersHahed, userHash)
	}
	// thêm vào db
	err := ser.repo.InsertMany(usersHahed)
	if err != nil {
		return err
	}
	return nil
}
func (ser *userService) ExportUsersExcel(users []models.User) (*excelize.File, error) {
	logr := logger.With("userService")

	f := excelize.NewFile()
	// tạo sheet mới
	sheet := "Users"
	f.NewSheet(sheet)

	// Xóa sheet mặc định
	f.DeleteSheet("Sheet1")

	// Đặt sheet "Users" làm active
	idx, err := f.GetSheetIndex(sheet)
	if err != nil {
		logr.Error("Có lỗi khi tạo file excel: " + err.Error())
		return nil, err
	}
	f.SetActiveSheet(idx)

	// gán header
	f.SetCellValue(sheet, "A1", "ID")
	f.SetCellValue(sheet, "B1", "Name")
	f.SetCellValue(sheet, "C1", "Email")

	// gán giá trị từng dòng
	for i, u := range users {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), u.Id)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), u.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), u.Email)
	}
	logr.Info("Tạo file excel thành công")
	return f, nil
}

func (ser *userService) GetAllUser() ([]models.User, error) {
	users, err := ser.repo.GetAllUser()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ser *userService) Get(id int) (*models.User, error) {
	logr := logger.With("userService")

	user, _ := ser.repo.Get(id)
	if user == nil {
		logr.Error("user không tồn tại")
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
	logr := logger.With("userService")

	// kiểm tra user
	user, _ := ser.repo.Get(userId)
	if user == nil {
		logr.Error("user không tồn tại")
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
