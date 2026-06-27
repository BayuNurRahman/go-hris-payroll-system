package usecase

import (
	"fmt"
	"strings"

	"hris-payroll/domain"
)

type EmployeeUsecase struct{}

func (u *EmployeeUsecase) CreateEmployee(employee domain.Employee) (domain.Employee, error) {
	if employee.Email == "" {
		return domain.Employee{}, fmt.Errorf("email is required")
	}
	if !strings.HasSuffix(strings.ToLower(employee.Email), "@company.co.id") {
		return domain.Employee{}, fmt.Errorf("email must use @company.co.id")
	}
	return employee, nil
}

func (u *EmployeeUsecase) GetEmployeeByID(id uint) (domain.Employee, error) {
	if id == 0 {
		return domain.Employee{}, fmt.Errorf("invalid employee id")
	}
	return domain.Employee{ID: id, Name: "Sample Employee"}, nil
}
