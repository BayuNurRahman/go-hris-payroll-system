package domain

import (
	"context"
	"time"
)

type LeaveRequest struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `gorm:"not null" json:"employee_id"`
	Employee   Employee  `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Reason     string    `gorm:"type:text;not null" json:"reason"`
	StartDate  time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate    time.Time `gorm:"type:date;not null" json:"end_date"`
	Status     string    `gorm:"size:20;not null" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateLeaveReq struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}

type LeaveRepository interface {
	CreateLeave(ctx context.Context, leave *LeaveRequest) error
	GetLeaveByID(ctx context.Context, leaveID uint, userID uint, userRole string) (LeaveRequest, error)
}

type LeaveUsecase interface {
	RequestLeave(ctx context.Context, employeeID uint, req CreateLeaveReq) error
	GetLeaveDetail(ctx context.Context, leaveID uint, userID uint, userRole string) (LeaveRequest, error)
}
