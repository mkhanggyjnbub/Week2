package Middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

type Claims struct {
	UserID               uint   // ID của user
	Email                string // vai trò
	Role                 string // vai trò
	jwt.RegisteredClaims        // các trường chuẩn JWT, ví dụ Expiry
}

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// lấy token từ cookie
// 		tokenString, err := c.Cookie("token")
// 		if err != nil {
// 			c.Redirect(http.StatusFound, "/login")
// 			c.Abort()
// 			return
// 		}

// 		// parse token
// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return JwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.Redirect(http.StatusFound, "/login")
// 			c.Abort()
// 			return
// 		}

// 		// Lưu thông tin user vào context (có thể dùng sau này)
// 		c.Set("userID", claims.UserID)
// 		c.Set("role", claims.Role)

// 		// cho qua
// 		c.Next()
// 	}
// }

// Middleware xác thực và phân quyền
func AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Lấy từ Header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Nếu không có Header thì lấy từ Cookie
		if tokenString == "" {
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenString = cookie
			}
		}

		// 3. Nếu vẫn rỗng thì báo lỗi
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Thiếu token"})
			c.Abort()
			return
		}

		// 4. Parse token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
			c.Abort()
			return
		}

		// 5. Kiểm tra role
		authorized := false
		for _, r := range roles {
			if claims.Role == r {
				authorized = true
				break
			}
		}
		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Không đủ quyền"})
			c.Abort()
			return
		}

		// 6. Lưu claims vào context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
