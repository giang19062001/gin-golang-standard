package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	MaxUploadSize = int64(10 << 20) //  =  10 * (2^20) = 10 * 1,048,576 = 10,485,760 bytes = 10MB
	UploadDir     = "./uploads"
	AllowedTypes  = []string{".jpg", ".jpeg", ".png"}
)

// Validate file
func ValidateFile(header *multipart.FileHeader) error {
	if header.Size > MaxUploadSize {
		return fmt.Errorf("file too large, max %dMB", MaxUploadSize>>20) // (10,485,760 / (2^20)) = 10
	}
	ext := strings.ToLower(filepath.Ext(header.Filename)) // Lấy phần đuôi mở rộng .jpg, .png
	if !contains(AllowedTypes, ext) {
		return fmt.Errorf("file type not allowed, only: jpg, png")
	}
	return nil
}

// kiểm tra đuôi file
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// tạo tên uuid
func GenerateFileName(prefix, originalName string) string {
	ext := filepath.Ext(originalName) // lấy phần mở rộng .jpg, .png, ...
	newName := uuid.New().String() + ext
	log.Print(prefix + "-" + newName)
	return prefix + "-" + newName
}

// lưu file local
func SaveFile(file multipart.File, header *multipart.FileHeader, savePath string) (string, error) {
	defer file.Close()

	// Tạo folder nếu chưa có
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil { // * os.ModePerm = 0777 → full quyền đọc/ghi/thực thi
		return "", err
	}

	// tạo file mới hoặc ghi đè file cũ nếu tồn tại
	out, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy nội dung từ upload vào file mới
	_, err = io.Copy(out, file)
	return savePath, err
}

// xóa file dựa vào prefix đầu (userId-)
func DeleteFilesByPrefix(prefix string) error {
	files, err := os.ReadDir(UploadDir)
	if err != nil {
		return err
	}

	found := false // đánh dấu có xóa được file nào không

	for _, f := range files {
		// Bỏ qua folder và chỉ xét file có prefix
		if !f.IsDir() && strings.HasPrefix(f.Name(), prefix+"-") {
			oldPath := filepath.Join(UploadDir, f.Name())

			if err := os.Remove(oldPath); err != nil {
				log.Printf("Lỗi khi xóa file %s: %v", f.Name(), err)
				return err
			}

			log.Printf("Đã xóa file cũ: %s", f.Name())
			found = true
		}
	}

	if !found {
		log.Printf("Không tìm thấy file nào có prefix: %s", prefix)
	}

	return nil
}

func SaveAndGetURL(prefix string, file multipart.File, header *multipart.FileHeader) string {

	// 1. Tạo tên file mới
	filename := GenerateFileName(prefix, header.Filename)
	savePath := filepath.Join(UploadDir, filename)

	// 2. Lưu file mới
	if _, err := SaveFile(file, header, savePath); err != nil {
		log.Printf("Lưu file thất bại: %v", err)
		return ""
	}

	// 3. Trả về URL
	return "/uploads/" + filename
}
