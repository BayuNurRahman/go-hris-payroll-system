package domain

import (
	"context"
	"strings"
	"time"
)

type UserRole string

const (
	RoleHRD      UserRole = "HRD"
	RoleEmployee UserRole = "EMPLOYEE"
)

type RegisterReq struct {
	Name         string   `json:"name" binding:"required"`
	Email        string   `json:"email" binding:"required,email"`
	Password     string   `json:"password" binding:"required,min=6"`
	Role         UserRole `json:"role" binding:"required,oneof=HRD EMPLOYEE"`
	BaseSalary   float64  `json:"base_salary" binding:"required,gt=0"`
	DepartmentID uint     `json:"department_id" binding:"required"`
}

func (r *RegisterReq) ValidateEmailDomain() bool {
	return strings.HasSuffix(strings.ToLower(r.Email), "@company.co.id")
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenBlacklist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type AuthUsecase interface {
	Register(ctx context.Context, req RegisterReq) error
	Login(ctx context.Context, req LoginReq) (string, error)
	Logout(ctx context.Context, token string) error
}