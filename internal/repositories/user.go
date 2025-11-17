package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/giang19062001/gin-golang-standard/internal/models"
)

type userRepository struct {
	db *sql.DB
}

// ! Vì Struct đã implement Interface này nên nó thể trả ra Interface đó
func NewUserRepository(db *sql.DB) IUserRepository {
	return &userRepository{db: db}
}

type IUserRepository interface {
	GetAllUser() ([]models.User, error)
	Insert(*models.User) error
	UpdateAvatar(string, int) error
	Get(int) (*models.User, error)
	GetByEmail(string) (*models.User, error)
	GetUserOfEvent(int) ([]models.User, error)
}

func (repo *userRepository) GetAllUser() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id, name, email, password, avatar FROM users"
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // đóng kết quả sau khi duyệt xong

	users := []models.User{}
	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Avatar)
		if err != nil {
			log.Printf("Lỗi scan user: %v", err)
			continue
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil { // kiểm lỗi nếu có khi duyệt mảng
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) Insert(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users(email, password, name) VALUES (?, ?, ?)"
	result, err := repo.db.ExecContext(ctx, query, user.Email, user.Password, user.Name)
	if err != nil {
		return err
	}

	// lấy id vừa được thêm
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// gán cho user
	user.Id = int(id)
	return nil

}

func (repo *userRepository) UpdateAvatar(avatar string, userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "UPDATE users SET avatar = ? WHERE id = ? "
	_, err := repo.db.ExecContext(ctx, query, avatar, userId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) GetUserOfEvent(eventId int) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := ` SELECT u.id, u.name, u.email
			FROM users u
			JOIN attendees a 
			ON u.id = a.user_id
			WHERE a.event_id = ? `

	rows, err := repo.db.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() { // duyệt từng dòng record
		var user models.User
		// Scan : gán từng record vào struct
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil { // kiểm lỗi nếu có khi duyệt mảng
		return nil, err
	}

	return users, nil

}

func (repo *userRepository) Get(id int) (*models.User, error) {
	query := "SELECT  id, email, name, password, avatar FROM users WHERE id = ?"
	return repo.getCommon(query, id)
}

func (repo *userRepository) GetByEmail(email string) (*models.User, error) {
	query := "SELECT  id, email, name, password, avatar FROM users WHERE email = ?"
	return repo.getCommon(query, email)
}

// args truyền bấy nhiêu đối số với nhiều kiểu dữ liệu khác nhau tùy ý
func (repo *userRepository) getCommon(query string, args ...interface{}) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	// args...  liệt kê tất cả đối số vào trong hàm này
	err := repo.db.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.Avatar)
	log.Print(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
