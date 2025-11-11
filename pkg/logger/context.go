package logger

import (
	"runtime"
	"strings"

	"go.uber.org/zap"
)

func With(service string) *zap.SugaredLogger {
	// Lấy thông tin hàm đang gọi With()
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return Info.With("service", service, "func", "unknown")
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return Info.With("service", service, "func", "unknown")
	}

	fullName := fn.Name()
	// fullName: "github.com/.../internal/services.(*userService).LoginUser"

	funcName := fullName
	if i := strings.LastIndex(fullName, "."); i >= 0 {
		funcName = fullName[i+1:] // cắt lấy từ dấu chấm cuối
	}
	// Nếu là method của struct → bỏ "(type)."
	if i := strings.Index(funcName, "-"); i >= 0 {
		funcName = funcName[i+1:]
	}

	return Info.With("service", service, "func", funcName)
}

// log Info + Warn + Error + Fatal → ghi vào cả 2 fileinfo-*.log và error-*.log
func WithService(service string) *zap.SugaredLogger {
	return Info.With("service", service)
}

// log chỉ Error + Fatal → ghi vào cả 2 file, nhưng không ghi các log Info/Warn
func WithServiceError(service string) *zap.SugaredLogger {
	return Error.With("service", service)
}
