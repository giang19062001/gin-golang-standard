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
	date := time.Now().Format("2006-01-02")
	infoPath := filepath.Join(logDir, "info-"+date+".log")
	errorPath := filepath.Join(logDir, "error-"+date+".log")

	infoWriter := getWriter(infoPath)
	errorWriter := getWriter(errorPath)

	// cấu hình encoder cho file (JSON)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// cấu hình encoder cho console
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.TimeKey = "timestamp"
	consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// core cho file
	infoCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(infoWriter), zapcore.InfoLevel)
	errorCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(errorWriter), zapcore.ErrorLevel)

	// core cho console (ghi tất cả từ Info trở lên ra màn hình)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	// kết hợp tất cả core
	logger = zap.New(zapcore.NewTee(infoCore, errorCore, consoleCore))
	// logger = zap.New(zapcore.NewTee(infoCore, errorCore, consoleCore),
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
