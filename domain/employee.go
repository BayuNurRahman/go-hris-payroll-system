package domain

import "time"

type Employee struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	Name         string           `gorm:"size:100;not null" json:"name"`
	Email        string           `gorm:"uniqueIndex;not null" json:"email"`
	Password     string           `gorm:"not null" json:"-"`
	Role         string           `gorm:"size:20;not null" json:"role"`
	BaseSalary   float64          `gorm:"type:decimal(12,2);not null" json:"base_salary"`
	DepartmentID uint             `gorm:"not null" json:"department_id"`
	Department   DepartmentBudget `gorm:"foreignKey:DepartmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"department,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type DepartmentBudget struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name	   string    `gorm:"size:100;not null" json:"name"`
	BudgetLeft float64   `gorm:"type:decimal(12,2);not null" json:"budget_left"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type EmployeeUsecase interface {
	CreateEmployee(employee Employee) (Employee, error)
	GetEmployeeByID(id uint) (Employee, error)
}