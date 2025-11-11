package services

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/models"
	"github.com/giang19062001/gin-golang-standard/internal/repositories"
	"github.com/giang19062001/gin-golang-standard/internal/utils"
)

type eventService struct {
	repo        repositories.IEventRepository
	userService IUserService
}

type IEventService interface {
	CreateEvent(*dto.CreateEventDto) (*models.Event, error)
	UpdateEvent(int, *dto.UpdateEventDto) (*models.Event, error)
	GetAllEvents() ([]models.Event, error)
	GetEventsByUser(int) ([]models.Event, error)
	GetEvent(int) (*models.Event, error)
	DeleteEvent(int) (string, error)
	UploadImages([]*multipart.FileHeader, int) ([]string, error)
}

func NewEventService(repo repositories.IEventRepository, userService IUserService) IEventService {
	return &eventService{repo: repo, userService: userService}
}

func (ser *eventService) CreateEvent(newEvent *dto.CreateEventDto) (*models.Event, error) {
	user, _ := ser.userService.Get(newEvent.OwnerId) // * inject 1 server khác
	if user == nil {
		return nil, errors.New("user không tồn tại")
	}

	id, err := ser.repo.Insert(newEvent)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, errors.New("lỗi khi tạo event mới")
	}

	event := models.Event{
		Id:          id,
		OwnerId:     newEvent.OwnerId,
		Name:        newEvent.Name,
		Description: newEvent.Description,
		Date:        newEvent.Date,
		Location:    newEvent.Location,
	}
	return &event, nil
}

func (ser *eventService) UpdateEvent(id int, changesEvent *dto.UpdateEventDto) (*models.Event, error) {

	event, err := ser.repo.Get(id)
	if event == nil {
		return nil, errors.New("event không tồn tại")
	}
	if err != nil {
		return nil, errors.New("lỗi trích xuất event")
	}

	if err := ser.repo.Update(id, changesEvent); err != nil {
		return nil, errors.New("lỗi cập nhập event")
	}

	result := &models.Event{
		Id:          id,
		OwnerId:     event.OwnerId,
		Name:        changesEvent.Name,
		Description: changesEvent.Description,
		Date:        changesEvent.Date,
		Location:    changesEvent.Location,
	}
	return result, nil

}

func (ser *eventService) GetAllEvents() ([]models.Event, error) {
	events, err := ser.repo.GetAll()
	if err != nil {

		return []models.Event{}, errors.New("lỗi trích xuất event")
	}
	return events, nil
}

func (ser *eventService) GetEvent(id int) (*models.Event, error) {
	// lấy info của event
	event, err := ser.repo.Get(id)
	if event == nil {
		return nil, errors.New("event không tồn tại")
	}
	if err != nil {
		return nil, err
	}
	// lấy danh sách ảnh của event
	eventImgs, err := ser.repo.GetImgsById(id)
	if err != nil {
		return nil, err
	}

	eventMerge := &models.Event{
		Id:          event.Id,
		OwnerId:     event.OwnerId,
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		Location:    event.Location,
		EventImgs:   eventImgs,
	}
	return eventMerge, nil
}

func (ser *eventService) DeleteEvent(id int) (string, error) {
	event, _ := ser.repo.Get(id)
	if event == nil {
		return "", errors.New("event không tồn tại")
	}

	if err := ser.repo.Delete(id); err != nil {
		return "", errors.New("lỗi trích xuất event")
	}
	return "Event đã được xóa", nil
}

func (ser *eventService) GetEventsByUser(userId int) ([]models.Event, error) {
	events, err := ser.repo.GetEventsByUser(userId)
	if err != nil {
		return []models.Event{}, errors.New("lỗi trích xuất event")
	}
	return events, nil
}

func (ser *eventService) UploadImages(files []*multipart.FileHeader, eventId int) ([]string, error) {
	var filepaths []string
	prefix := fmt.Sprintf("event-%d", eventId)

	// kiểm tra event
	event, err := ser.repo.Get(eventId)
	if event == nil {
		return []string{}, errors.New("event không tồn tại")
	}
	if err != nil {
		return []string{}, errors.New("lỗi trích xuất event")
	}

	// Xóa file cũ trong folder
	err = utils.DeleteFilesByPrefix(prefix)
	if err != nil {
		return []string{}, err
	}
	// Xóa file cũ trong database
	err = ser.repo.DeleteImgs(eventId)
	if err != nil {
		return []string{}, errors.New("lỗi khi xóa ảnh cũ của event trong cơ sở dữ liệu")
	}

	// add file mới vào folder
	for _, header := range files {
		// mở đọc file
		f, err := header.Open()
		if err != nil {
			return []string{}, errors.New("mở file bị lỗi")
		}

		// lưu file vào thư mục
		filename := utils.SaveAndGetURL(prefix, f, header)
		filepaths = append(filepaths, filename)
	}

	// add file mới vào database
	for _, path := range filepaths {
		var eventImg = &models.EventImg{
			EventId:  eventId,
			Filepath: path,
		}
		_, err := ser.repo.InsertImgs(eventImg)
		if err != nil {
			return []string{}, errors.New("lỗi khi lưu ảnh mới của event trong cơ sở dữ liệu")
		}
	}
	return filepaths, nil
}
