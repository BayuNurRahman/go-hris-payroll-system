package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Otorisasi tidak ditemukan"})
			c.Abort()
			return
		}

		// Cek apakah role user ada di daftar role yang diizinkan
		for _, role := range allowedRoles {
			if role == userRole.(string) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Anda tidak memiliki hak akses untuk fitur ini"})
		c.Abort()
	}
}