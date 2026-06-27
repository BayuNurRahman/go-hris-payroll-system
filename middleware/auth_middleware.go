package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"hris-payroll/domain"
)

func JWTAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Header otorisasi diperlukan"})
			c.Abort()
			return
		}

		// Ambil string token setelah kata "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid (wajib Bearer [token])"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// KRITERIA: Periksa apakah token ada di database Blacklist
		var blacklisted domain.TokenBlacklist
		err := db.Where("token = ?", tokenString).First(&blacklisted).Error
		if err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi telah berakhir, silakan login kembali (Token Blacklisted)"})
			c.Abort()
			return
		}

		// Parse & Validasi Token JWT
		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token kedaluwarsa atau tidak valid"})
			c.Abort()
			return
		}

		// Simpan data user ke context Gin agar bisa dibaca di layer Handler
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("user_role", claims["user_role"].(string))
		c.Set("token_string", tokenString) // Berguna saat proses logout

		c.Next()
	}
}