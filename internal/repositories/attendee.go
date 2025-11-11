package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/giang19062001/gin-golang-standard/internal/dto"
	"github.com/giang19062001/gin-golang-standard/internal/models"
)

type attendeeRepository struct {
	db *sql.DB
}

type IAttendeeRepository interface {
	Insert(*dto.CreateAttendeeDto) (int, error)
	GetAttendeeByEventAndUser(int, int) (*models.Attendee, error)
	Delete(int, int) error
}

func NewAttendeRepository(db *sql.DB) IAttendeeRepository {
	return &attendeeRepository{db: db}
}

func (repo *attendeeRepository) Insert(dto *dto.CreateAttendeeDto) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO attendees (event_id, user_id) VALUES (?, ?)"
	result, err := repo.db.ExecContext(ctx, query, dto.EventId, dto.UserId)
	if err != nil {
		return 0, err
	}

	// lấy id vừa được thêm
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// trả về số id vừa insert
	return int(id), nil
}

func (repo *attendeeRepository) GetAttendeeByEventAndUser(eventId, userId int) (*models.Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, event_id, user_id FROM attendees WHERE event_id = ? AND user_id = ? "
	var attendee models.Attendee
	err := repo.db.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.Id, &attendee.EventId, &attendee.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attendee, nil
}

func (repo *attendeeRepository) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "DELETE FROM attendees WHERE user_id = ? AND event_id = ?"
	_, err := repo.db.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}
	return nil
}
