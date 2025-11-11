// pkg/logger/logger.go
package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Info   *zap.SugaredLogger
	Error  *zap.SugaredLogger
	logger *zap.Logger
)

func Init() {
	logDir := "./logs"
	// tạo thư mục LOGS
	// 	7 → read + write + execute (rwx) cho owner
	// 5 → read + execute (rx) cho group
	// 5 → read + execute (rx) cho others
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// Tạo filename theo ngày
	date := time.Now().Format("2006-01-02")
	infoPath := filepath.Join(logDir, "info-"+date+".log")
	errorPath := filepath.Join(logDir, "error-"+date+".log")

	infoWriter := getWriter(infoPath)
	errorWriter := getWriter(errorPath)

	// chuẩn cấu hình
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	infoCore := zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), zapcore.InfoLevel)
	errorCore := zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), zapcore.ErrorLevel)

	// Tee kết hợp 2 core
	logger = zap.New(zapcore.NewTee(infoCore, errorCore))
	// logger = zap.New(zapcore.NewTee(infoCore, errorCore),
	// 	zap.AddCaller(), // thêm thông tin caller - người gọi thông thường caller
	// )

	Info = logger.Sugar()
	Error = logger.Sugar()
}

func getWriter(filename string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,  // Mỗi file tối đa 100MB trước khi tạo file mới: info-2025-11-11.log → khi đầy → info-2025-11-11.log.001
		MaxBackups: 30,   // Giữ tối đa 30 file cũ (kể cả đã nén)
		MaxAge:     30,   // File quá 30 ngày sẽ bị xóa tự động
		Compress:   true, // File cũ sẽ được nén .gz
		LocalTime:  true,
	})
}

// Sync khi shutdown
func Sync() {
	_ = Info.Sync()
	_ = Error.Sync()
}
