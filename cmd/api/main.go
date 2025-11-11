package main

import (
	"log"

	"github.com/giang19062001/gin-golang-standard/internal/config"
	"github.com/giang19062001/gin-golang-standard/internal/database"
	"github.com/giang19062001/gin-golang-standard/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DbUrl)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer db.Close() // đóng khi app tắt

	// khởi tạo gin
	r := gin.Default()

	// Public folder /uploads để có thể đọc ảnh
	r.Static("/uploads", "./uploads")
	// Tạo router

	router.SetupRouter(r, cfg, db)

	log.Printf("Server running on :%s", cfg.Port)
	log.Printf("Swagger: http://localhost:%s/swagger/index.html", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
