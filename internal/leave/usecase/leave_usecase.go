package usecase

import (
	"context"
	"hris-payroll/domain"
	"hris-payroll/internal/leave/repository"
	"time"
)

type LeaveUsecase struct {
	leaveRepo *repository.LeaveRepository
}

func NewLeaveUsecase(repo *repository.LeaveRepository) *LeaveUsecase {
	return &LeaveUsecase{leaveRepo: repo}
}

func (u *LeaveUsecase) RequestLeave(ctx context.Context, employeeID uint, req domain.CreateLeaveReq) error {
	// Parsing string "YYYY-MM-DD" ke time.Time
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return err
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return err
	}

	leave := domain.LeaveRequest{
		EmployeeID: employeeID,
		Reason:     req.Reason,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     "PENDING",
		CreatedAt:  time.Now(),
	}

	return u.leaveRepo.CreateLeave(ctx, &leave)
}

func (u *LeaveUsecase) GetLeaveDetail(ctx context.Context, leaveID uint, userID uint, userRole string) (domain.LeaveRequest, error) {
	return u.leaveRepo.GetLeaveByID(ctx, leaveID, userID, userRole)
}
