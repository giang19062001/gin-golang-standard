package models

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"` // "-" khi lấy dữ liệu json thì password sẽ được ẩn
	Avatar   string `json:"avatar"`
}
