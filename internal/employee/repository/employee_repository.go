package repository

import "hris-payroll/domain"

type EmployeeRepository struct{}

func (r *EmployeeRepository) Create(employee domain.Employee) (domain.Employee, error) {
	return employee, nil
}

func (r *EmployeeRepository) FindByID(id uint) (domain.Employee, error) {
	return domain.Employee{ID: id, Name: "Sample Employee"}, nil
}
