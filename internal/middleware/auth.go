package middlewares

import (
	"net/http"
	"strings"

	"github.com/giang19062001/gin-golang-standard/internal/services"
	"github.com/giang19062001/gin-golang-standard/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string, userService services.IUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header là bắt buộc"})
			c.Abort()
			return
		}

		// * kiểm tra Authorization header có phải là dạng "Bearer {token}" hay không
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header định dạng không hợp lệ - thiếu Bearer"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// kiểm tra tính hợp lệ của token với jwtSecret
		token, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "token không hợp lệ",
			})
			c.Abort()
			return
		}

		// * lấy thông tin (claims) từ token
		claims, ok := utils.GetClaims(token)
		if !ok {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "không thể đọc thông tin của token",
			})
			c.Abort()
			return
		}

		// * khi decode JSON (hoặc JWT), tất cả các số (number) đều được mặc định parse thành "float64", không phải "int"
		userId := claims["userId"].(float64)      // đây là cách parse 1 interface{} -> sang 1 kiểu thông thường
		user, err := userService.Get(int(userId)) // đây là ép kiểu thông thường từ "float" sang "int"
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "không thể xác thực",
			})
			c.Abort()
			return
		}
		// * lưu dữ liệu tạm vào context => Dữ liệu này có thể lấy lại ở bất kỳ Controller nào chạy sau middleware
		c.Set("user", user)
		c.Next()

	}
}
