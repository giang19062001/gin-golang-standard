// pkg/database/database.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("lỗi không thể mở cổng kết nối: %v", err)
	}

	db.SetMaxOpenConns(100)                 // Max connections đang mở
	db.SetMaxIdleConns(10)                  // Max idle connections
	db.SetConnMaxLifetime(time.Hour)        // Đóng connection cũ sau 1 giờ
	db.SetConnMaxIdleTime(10 * time.Minute) // Idle quá 10 phút thì đóng

	// Test kết nối ngay
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("lỗi kết nối database: %v", err)
	}

	log.Println("MySQL kết nối thành công")
	return db, nil
}
