package services

import (
	"errors"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/models"
	"github.com/giang19062001/gin-golang-standard/internal/repositories"
	"github.com/giang19062001/gin-golang-standard/pkg/logger"
)

type attendeeService struct {
	repo         repositories.IAttendeeRepository
	userService  IUserService
	eventService IEventService
	emailService IEmailService
}

type IAttendeeService interface {
	RegisterAttendee(*dto.CreateAttendeeDto) (*models.Attendee, error)
	Delete(int, int) error
	GetEventsByUser(int) ([]models.Event, error)
	GetUsersByEvent(int) ([]models.User, error)
}

func NewAttendeeService(repo repositories.IAttendeeRepository, userService IUserService, eventService IEventService, emailService IEmailService) IAttendeeService {
	return &attendeeService{repo: repo, userService: userService, eventService: eventService, emailService: emailService}
}

func (ser *attendeeService) RegisterAttendee(dto *dto.CreateAttendeeDto) (*models.Attendee, error) {
	logr := logger.With("attendeeService")

	// kiểm tra event
	event, err := ser.eventService.GetEvent(dto.EventId)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event không tồn tại")
	}

	// kiểm tra user
	user, err := ser.userService.Get(dto.UserId)

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user không tồn tại")

	}

	// kiểm tra user này đã tham dự sự kiện này chưa
	existingAttendee, err := ser.repo.GetAttendeeByEventAndUser(event.Id, user.Id)
	if err != nil {
		return nil, err
	}
	if existingAttendee != nil {
		return nil, errors.New("user hiện tại đã tham dự event này rồi")
	}

	id, err := ser.repo.Insert(dto)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, errors.New("lỗi khi thêm dữ liệu mới")
	}
	// gửi sự kiện cho serivce khac -> service đó đãm nhiệm trọng trách gửi email
	logr.Info("Email : ", user)

	ser.emailService.SendEmailHandler(user.Email, user.Name)

	// trả về thông tin đăng ký event
	attende := &models.Attendee{
		Id:      id,
		UserId:  dto.UserId,
		EventId: dto.EventId,
	}
	return attende, nil
}

func (ser *attendeeService) Delete(userId, eventId int) error {
	// kiểm tra event
	event, err := ser.eventService.GetEvent(eventId)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event không tồn tại")
	}
	// kiểm tra user
	user, err := ser.userService.Get(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user không tồn tại")

	}
	// kiểm tra user, sự kiện đã có đăng ký trước đó chưa
	existingAttendee, err := ser.repo.GetAttendeeByEventAndUser(event.Id, user.Id)
	if err != nil {
		return err
	}
	if existingAttendee == nil {
		return errors.New("user hiện tại chưa từng tham dự event này - không thể xóa")
	}

	// xóa
	err = ser.repo.Delete(userId, eventId)
	if err != nil {
		return err
	}

	return nil
}

func (ser *attendeeService) GetEventsByUser(userId int) ([]models.Event, error) {
	// kiểm tra user
	user, err := ser.userService.Get(userId)
	if err != nil {
		return []models.Event{}, err
	}
	if user == nil {
		return []models.Event{}, errors.New("user không tồn tại")

	}

	events, err := ser.eventService.GetEventsByUser(userId)
	if err != nil {
		return []models.Event{}, err
	}
	if len(events) == 0 {
		return []models.Event{}, errors.New("user này chưa từng đăng ký event nào")
	}

	return events, nil
}

func (ser *attendeeService) GetUsersByEvent(eventId int) ([]models.User, error) {
	// kiểm tra event
	event, err := ser.eventService.GetEvent(eventId)
	if err != nil {
		return []models.User{}, err
	}
	if event == nil {
		return []models.User{}, errors.New("event không tồn tại")
	}

	users, err := ser.userService.GetUserOfEvent(eventId)
	if err != nil {
		return []models.User{}, err
	}
	if len(users) == 0 {
		return []models.User{}, errors.New("event này chưa từng có user nào đăng ký")
	}

	return users, nil
}
