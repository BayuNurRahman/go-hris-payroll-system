package usecase

import (
	"context"
	"errors"
	"os"
	"time"

	"hris-payroll/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authUsecase struct {
	db *gorm.DB
}

func NewAuthUsecase(db *gorm.DB) domain.AuthUsecase {
	return &authUsecase{db: db}
}

func (u *authUsecase) Register(ctx context.Context, req domain.RegisterReq) error {
	// 1. Jalankan Custom Validator domain email
	if !req.ValidateEmailDomain() {
		return errors.New("email harus menggunakan domain resmi @company.co.id")
	}

	// 2. Cek apakah email sudah terdaftar
	var existingEmployee domain.Employee
	err := u.db.WithContext(ctx).Where("email = ?", req.Email).First(&existingEmployee).Error
	if err == nil {
		return errors.New("email sudah terdaftar di sistem")
	}

	// 3. Hash password menggunakan Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4. Map DTO ke Core Entity Employee
	newEmployee := domain.Employee{
		Name:         req.Name,
		Email:        req.Email,
		Password:     string(hashedPassword),
		Role:         string(req.Role),
		BaseSalary:   req.BaseSalary,
		DepartmentID: req.DepartmentID,
	}

	// 5. Simpan ke database
	return u.db.WithContext(ctx).Create(&newEmployee).Error
}

func (u *authUsecase) Login(ctx context.Context, req domain.LoginReq) (string, error) {
	var emp domain.Employee
	// 1. Cari user berdasarkan email
	if err := u.db.WithContext(ctx).Where("email = ?", req.Email).First(&emp).Error; err != nil {
		return "", errors.New("email atau password salah")
	}

	// 2. Validasi password menggunakan Bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(emp.Password), []byte(req.Password)); err != nil {
		return "", errors.New("email atau password salah")
	}

	// 3. Generate JWT Token jika password cocok
	secretKey := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id":   emp.ID,
		"user_role": emp.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Berlaku 24 Jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *authUsecase) Logout(ctx context.Context, token string) error {
	// Masukkan token ke tabel Blacklist agar tidak bisa dipakai lagi
	blacklist := domain.TokenBlacklist{
		Token:     token,
		ExpiredAt: time.Now().Add(time.Hour * 24), // Set masa aman kadaluarsa
	}
	return u.db.WithContext(ctx).Create(&blacklist).Error
}
