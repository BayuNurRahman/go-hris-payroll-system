package repository

import (
	"context"
	"fmt"
	"time"

	"hris-payroll/domain"

	"gorm.io/gorm"
)

type PayrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) *PayrollRepository {
	return &PayrollRepository{db: db}
}

func (r *PayrollRepository) ExecutePayrollTransaction(ctx context.Context, req domain.ProcessPayrollReq) error {
	// Memulai Transaksi Database (ACID)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Ambil data karyawan untuk mendapatkan BaseSalary dan DepartmentID
		var emp domain.Employee
		if err := tx.First(&emp, req.EmployeeID).Error; err != nil {
			return fmt.Errorf("karyawan dengan ID %d tidak ditemukan", req.EmployeeID)
		}

		// 2. Terapkan Row-Level Locking (FOR UPDATE) pada budget departemen terkait
		// Ini mencegah race condition jika ada proses payroll lain yang berjalan bersamaan
		var dept domain.DepartmentBudget
		err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&dept, emp.DepartmentID).Error
		if err != nil {
			return fmt.Errorf("departemen dengan ID %d tidak ditemukan", emp.DepartmentID)
		}

		// 3. Hitung Gaji Total = Gaji Pokok + Bonus
		totalPaid := emp.BaseSalary + req.Bonus

		// 4. Validasi kecukupan anggaran divisi
		if dept.BudgetLeft < totalPaid {
			// Mengembalikan error otomatis memicu ROLLBACK seluruh transaksi
			return fmt.Errorf("transaksi batal: sisa budget divisi %s tidak mencukupi (Sisa: %.2f, Kebutuhan: %.2f)",
				dept.Name, dept.BudgetLeft, totalPaid)
		}

		// 5. Potong anggaran divisi dan update ke database
		dept.BudgetLeft -= totalPaid
		if err := tx.Save(&dept).Error; err != nil {
			return err
		}

		// 6. Catat riwayat ke tabel payrolls
		payrollRecord := domain.Payroll{
			EmployeeID: emp.ID,
			Bonus:      req.Bonus,
			TotalPaid:  totalPaid,
			PaidAt:     time.Now(),
		}
		if err := tx.Create(&payrollRecord).Error; err != nil {
			return err
		}

		// Jika mengembalikan nil, GORM otomatis melakukan COMMIT transaksi
		return nil
	})
}
