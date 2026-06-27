package domain

import (
	"context"
	"time"
)

type Payroll struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `gorm:"not null" json:"employee_id"`
	Bonus      float64   `gorm:"type:decimal(12,2);not null" json:"bonus"`
	TotalPaid  float64   `gorm:"type:decimal(12,2);not null" json:"total_paid"`
	PaidAt     time.Time `json:"paid_at"`
}

type ProcessPayrollReq struct {
	EmployeeID uint    `json:"employee_id" binding:"required"`
	Bonus      float64 `json:"bonus" binding:"required,gt=-1"`
}

type ProcessPayrollDTO = ProcessPayrollReq

type PayrollRepository interface {
	ExecutePayrollTransaction(ctx context.Context, req ProcessPayrollReq) error
}

type PayrollUsecase interface {
	ProcessBatchPayroll(ctx context.Context, requests []ProcessPayrollReq) ([]string, error)
}
