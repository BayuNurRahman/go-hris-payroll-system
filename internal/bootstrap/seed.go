package bootstrap

import (
	"fmt"
	"os"

	"hris-payroll/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	var deptCount int64
	if err := db.Model(&domain.DepartmentBudget{}).Count(&deptCount).Error; err != nil {
		return err
	}

	if deptCount == 0 {
		if err := db.Create(&domain.DepartmentBudget{ID: 1, Name: "Engineering", BudgetLeft: 100000000}).Error; err != nil {
			return err
		}
		if err := db.Create(&domain.DepartmentBudget{ID: 2, Name: "HR", BudgetLeft: 50000000}).Error; err != nil {
			return err
		}
	}

	var employeeCount int64
	if err := db.Model(&domain.Employee{}).Count(&employeeCount).Error; err != nil {
		return err
	}

	if employeeCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SEED_PASSWORD")), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := db.Create(&domain.Employee{
			Name:         "HRD Admin",
			Email:        "hrd@company.co.id",
			Password:     string(hashedPassword),
			Role:         "HRD",
			BaseSalary:   15000000,
			DepartmentID: 2,
		}).Error; err != nil {
			return err
		}
	}

	fmt.Println("Seed data initialized")
	return nil
}
