package usecase

import (
	"testing"

	"hris-payroll/domain"
)

func TestCreateEmployee_RejectsInvalidEmailDomain(t *testing.T) {
	u := &EmployeeUsecase{}

	_, err := u.CreateEmployee(domain.Employee{
		Name:         "Test User",
		Email:        "test@other.com",
		Role:         "EMPLOYEE",
		BaseSalary:   1000,
		DepartmentID: 1,
	})

	if err == nil {
		t.Fatal("expected invalid email domain error")
	}
}

func TestCreateEmployee_AcceptsCompanyEmail(t *testing.T) {
	u := &EmployeeUsecase{}

	emp, err := u.CreateEmployee(domain.Employee{
		Name:         "Test User",
		Email:        "test@company.co.id",
		Role:         "EMPLOYEE",
		BaseSalary:   1000,
		DepartmentID: 1,
	})

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if emp.Email != "test@company.co.id" {
		t.Fatalf("expected employee email to be preserved")
	}
}
