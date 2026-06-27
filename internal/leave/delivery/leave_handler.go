package delivery

import (
	"net/http"
	"strconv"

	"hris-payroll/domain"
	"hris-payroll/internal/leave/usecase"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	usecase *usecase.LeaveUsecase
}

func NewLeaveHandler(u *usecase.LeaveUsecase) *LeaveHandler {
	return &LeaveHandler{usecase: u}
}

// CreateLeaveHandler - POST /api/v1/leaves
func (h *LeaveHandler) CreateLeaveHandler(c *gin.Context) {
	var req domain.CreateLeaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil ID pengaju dari JWT Token Context
	userID := c.MustGet("user_id").(uint)

	if err := h.usecase.RequestLeave(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengajukan cuti"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pengajuan cuti berhasil dikirim, menunggu persetujuan HRD"})
}

// GetLeaveDetailHandler - GET /api/v1/leaves/:id
func (h *LeaveHandler) GetLeaveDetailHandler(c *gin.Context) {
	leaveIDParam := c.Param("id")
	leaveID, err := strconv.Atoi(leaveIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Cuti tidak valid"})
		return
	}

	userID := c.MustGet("user_id").(uint)
	userRole := c.MustGet("user_role").(string)

	leave, err := h.usecase.GetLeaveDetail(c.Request.Context(), uint(leaveID), userID, userRole)
	if err != nil {
		// Proteksi IDOR mengembalikan status 404/403 jika data mencoba diintip oleh user lain
		c.JSON(http.StatusForbidden, gin.H{"error": "Data tidak ditemukan atau Anda tidak memiliki hak akses (IDOR Blocked)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data cuti berhasil diambil",
		"data":    leave,
	})
}
