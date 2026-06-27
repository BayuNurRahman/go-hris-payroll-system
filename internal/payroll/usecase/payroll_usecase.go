package usecase

import (
	"context"
	"fmt"
	"hris-payroll/domain"
	"hris-payroll/internal/payroll/repository"
)

type PayrollUsecase struct {
	payrollRepo *repository.PayrollRepository
}

func NewPayrollUsecase(repo *repository.PayrollRepository) *PayrollUsecase {
	return &PayrollUsecase{payrollRepo: repo}
}

func (u *PayrollUsecase) ProcessBatchPayroll(ctx context.Context, requests []domain.ProcessPayrollReq) ([]string, error) {
	var summaryLogs []string
	var hasFailure bool

	for _, req := range requests {
		err := u.payrollRepo.ExecutePayrollTransaction(ctx, req)
		if err != nil {
			hasFailure = true
			summaryLogs = append(summaryLogs, fmt.Sprintf("ID Karyawan %d GAGAL: %s", req.EmployeeID, err.Error()))
		} else {
			summaryLogs = append(summaryLogs, fmt.Sprintf("ID Karyawan %d SUKSES: Gaji berhasil diproses", req.EmployeeID))
		}
	}

	if hasFailure {
		return summaryLogs, fmt.Errorf("satu atau lebih transaksi payroll gagal")
	}

	return summaryLogs, nil
}
