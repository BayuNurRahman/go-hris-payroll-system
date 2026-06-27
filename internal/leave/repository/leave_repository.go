package repository

import (
	"context"
	"hris-payroll/domain"

	"gorm.io/gorm"
)

type LeaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

func (r *LeaveRepository) CreateLeave(ctx context.Context, leave *domain.LeaveRequest) error {
	return r.db.WithContext(ctx).Create(leave).Error
}

func (r *LeaveRepository) GetLeaveByID(ctx context.Context, leaveID uint, userID uint, userRole string) (domain.LeaveRequest, error) {
	var leave domain.LeaveRequest
	var err error

	if userRole == "HRD" {
		// HRD bebas melihat detail cuti milik siapapun
		err = r.db.WithContext(ctx).Preload("Employee").First(&leave, leaveID).Error
	} else {
		// MITIGASI IDOR: Query Defensif mengunci parameter ID Cuti sekaligus ID Pemiliknya
		err = r.db.WithContext(ctx).Preload("Employee").
			Where("id = ? AND employee_id = ?", leaveID, userID).
			First(&leave).Error
	}

	return leave, err
}
