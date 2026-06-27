package delivery

import (
	"net/http"

	"hris-payroll/domain"
	"hris-payroll/internal/payroll/usecase"

	"github.com/gin-gonic/gin"
)

type PayrollHandler struct {
	usecase *usecase.PayrollUsecase
}

func NewPayrollHandler(u *usecase.PayrollUsecase) *PayrollHandler {
	return &PayrollHandler{usecase: u}
}

// ProcessPayrollHandler - POST /api/v1/payroll/process
func (h *PayrollHandler) ProcessPayrollHandler(c *gin.Context) {
	// Mengharuskan input berupa JSON Array dari ProcessPayrollReq
	var reqs []domain.ProcessPayrollReq

	if err := c.ShouldBindJSON(&reqs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input payroll harus berupa list/array JSON yang valid"})
		return
	}

	if len(reqs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "List data payroll tidak boleh kosong"})
		return
	}

	// Eksekusi proses batch
	logs, err := h.usecase.ProcessBatchPayroll(c.Request.Context(), reqs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Proses payroll gagal",
			"results": logs,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Proses sinkronisasi transaksi payroll selesai dikerjakan",
		"results": logs,
	})
}
