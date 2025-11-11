package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/models"
)

type eventRepository struct {
	db *sql.DB
}

type IEventRepository interface {
	Insert(*dto.CreateEventDto) (int, error)
	InsertImgs(*models.EventImg) (int, error)
	Delete(int) error
	DeleteImgs(int) error
	Update(int, *dto.UpdateEventDto) error
	GetAll() ([]models.Event, error)
	Get(int) (*models.Event, error)
	GetEventsByUser(int) ([]models.Event, error)
	GetImgsById(int) ([]models.EventImg, error)
}

func NewEventRepository(db *sql.DB) IEventRepository {
	return &eventRepository{db: db}
}
func (repo *eventRepository) Insert(event *dto.CreateEventDto) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO events (owner_id, name, description, date, location) VALUES (?, ?, ?, ?, ?)"
	result, err := repo.db.ExecContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location)
	if err != nil {
		return 0, err
	}

	// Lấy dữ liệu dòng đầu tiên sau khi thêm sẽ trả ra ID mới và gán giá trị đó cho Id của event
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// trả ra số id vừa inserted được
	return int(id), nil
}
func (repo *eventRepository) InsertImgs(eventImg *models.EventImg) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO events_img (event_id, filepath) VALUES (?, ?)"
	result, err := repo.db.ExecContext(ctx, query, eventImg.EventId, eventImg.Filepath)
	if err != nil {
		return 0, err
	}

	// Lấy dữ liệu dòng đầu tiên sau khi thêm sẽ trả ra ID mới và gán giá trị đó cho Id của event
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// trả ra số id vừa inserted được
	return int(id), nil
}

func (repo *eventRepository) DeleteImgs(eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "DELETE FROM events_img WHERE event_id = ? "
	_, err := repo.db.ExecContext(ctx, query, eventId)
	if err != nil {
		return err
	}
	return nil
}
func (repo *eventRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "DELETE FROM events WHERE id = ? "
	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *eventRepository) Update(id int, event *dto.UpdateEventDto) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "UPDATE events SET name = ?, description = ?, date = ?, location = ? WHERE id = ? "
	_, err := repo.db.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *eventRepository) GetAll() ([]models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, owner_id, name, description, date, location FROM events"

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close() // đóng kết quả sau khi duyệt xong

	events := []models.Event{} // * nếu truct lớn -> nên dùng con trỏ  []*models.Event{} để giảm bộ nhớ hệ thống

	for rows.Next() { // duyệt từng dòng record
		var event models.Event
		// Scan : gán từng record vào struct
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil { // kiểm lỗi nếu có khi duyệt mảng
		return nil, err
	}

	return events, nil

}

func (repo *eventRepository) GetImgsById(eventId int) ([]models.EventImg, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, filepath, event_id FROM events_img WHERE event_id = ?"

	rows, err := repo.db.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}

	defer rows.Close() // đóng kết quả sau khi duyệt xong

	eventImgs := []models.EventImg{}

	for rows.Next() { // duyệt từng dòng record
		var event models.EventImg
		// Scan : gán từng record vào struct
		err := rows.Scan(&event.Id, &event.Filepath, &event.EventId)
		if err != nil {
			return nil, err
		}
		eventImgs = append(eventImgs, event)
	}

	if err = rows.Err(); err != nil { // kiểm lỗi nếu có khi duyệt mảng
		return nil, err
	}

	return eventImgs, nil

}
func (repo *eventRepository) Get(id int) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, owner_id, name, description, date, location FROM events WHERE id = ?"
	var event models.Event

	err := repo.db.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // lỗi ko tìm thấy dữ liệu -> ko phải lỗi hệ thống -> ko xem nó là lỗi
		}
		return nil, err
	}
	return &event, nil
}

func (repo *eventRepository) GetEventsByUser(userId int) ([]models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = ?
	`
	rows, err := repo.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
