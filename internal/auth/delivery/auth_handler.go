package delivery

import (
	"net/http"

	"hris-payroll/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// RegisterHandler - POST /api/v1/auth/register
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req domain.RegisterReq

	// 1. Validasi format JSON berdasarkan binding struct DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid: " + err.Error()})
		return
	}

	// 2. Panggil layer Usecase untuk mengeksekusi logika bisnis & custom validator
	if err := h.authUsecase.Register(c.Request.Context(), req); err != nil {
		// Mengembalikan error yang informatif (misal: domain salah atau email duplikat)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Karyawan berhasil didaftarkan ke dalam sistem!",
	})
}

// LoginHandler - POST /api/v1/auth/login
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req domain.LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan password wajib diisi dengan benar"})
		return
	}

	token, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   token,
	})
}

// LogoutHandler - POST /api/v1/auth/logout
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	// Ambil token dari context yang sudah divalidasi oleh JWTAuth middleware
	tokenString, exists := c.Get("token_string")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses logout"})
		return
	}

	if err := h.authUsecase.Logout(c.Request.Context(), tokenString.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memasukkan token ke daftar hitam"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout berhasil, sesi Anda telah dihapus",
	})
}
