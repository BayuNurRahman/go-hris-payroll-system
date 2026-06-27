package main

import (
	"fmt"
	"os"

	"hris-payroll/config"
	"hris-payroll/domain"
	authDelivery "hris-payroll/internal/auth/delivery"
	authUsecase "hris-payroll/internal/auth/usecase"
	"hris-payroll/internal/bootstrap"
	"hris-payroll/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	payrollDelivery "hris-payroll/internal/payroll/delivery"
	payrollRepository "hris-payroll/internal/payroll/repository"
	payrollUsecase "hris-payroll/internal/payroll/usecase"

	leaveDelivery "hris-payroll/internal/leave/delivery"
	leaveRepository "hris-payroll/internal/leave/repository"
	leaveUsecase "hris-payroll/internal/leave/usecase"
)

func main() {
	_ = godotenv.Load()

	// 1. Inisialisasi Koneksi ke PostgreSQL di Docker Desktop
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Println("Gagal terhubung ke database:", err)
		return
	}

	// 2. Jalankan Auto Migration untuk semua model sesuai kriteria ujian
	err = db.AutoMigrate(
		&domain.DepartmentBudget{},
		&domain.Employee{},
		&domain.Payroll{},
		&domain.LeaveRequest{},
		&domain.TokenBlacklist{},
	)
	if err != nil {
		fmt.Println("Gagal melakukan migrasi database:", err)
		return
	}
	fmt.Println("Auto Migration Berhasil: Semua tabel siap digunakan!")

	if err := bootstrap.SeedData(db); err != nil {
		fmt.Println("Gagal melakukan seed data:", err)
	}

	// 3. Inisialisasi Layer (Dependency Injection) untuk Auth Modul
	aUsecase := authUsecase.NewAuthUsecase(db)
	aHandler := authDelivery.NewAuthHandler(aUsecase)

	// 4. Dependency Injection - MODUL PAYROLL
	pRepo := payrollRepository.NewPayrollRepository(db)
	pUsecase := payrollUsecase.NewPayrollUsecase(pRepo)
	pHandler := payrollDelivery.NewPayrollHandler(pUsecase)

	// 5. Dependency Injection - MODUL LEAVE
	lRepo := leaveRepository.NewLeaveRepository(db)
	lUsecase := leaveUsecase.NewLeaveUsecase(lRepo)
	lHandler := leaveDelivery.NewLeaveHandler(lUsecase)

	// 5. Inisialisasi Engine Gin
	r := gin.Default()

	// Routing Definition
	api := r.Group("/api/v1")
	{
		// --- ENDPOINT PUBLIC / AUTHENTICATION ---
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", aHandler.LoginHandler)
			authGroup.POST("/register", middleware.JWTAuth(db), middleware.RequireRole("HRD"), aHandler.RegisterHandler)
			authGroup.POST("/logout", middleware.JWTAuth(db), aHandler.LogoutHandler)
		}

		// Payroll Routes (Dilindungi JWT & Hanya untuk HRD)
		payrollGroup := api.Group("/payroll")
		payrollGroup.Use(middleware.JWTAuth(db), middleware.RequireRole("HRD"))
		{
			payrollGroup.POST("/process", pHandler.ProcessPayrollHandler)
		}

		// Leave Routes (Butuh Login & Terbuka untuk Employee & HRD)
		leaveGroup := api.Group("/leaves")
		leaveGroup.Use(middleware.JWTAuth(db))
		{
			leaveGroup.POST("", lHandler.CreateLeaveHandler)
			leaveGroup.GET("/:id", lHandler.GetLeaveDetailHandler)
		}

		api.GET("/ping", middleware.JWTAuth(db), func(c *gin.Context) {
			role, _ := c.Get("user_role")
			c.JSON(200, gin.H{
				"message": "Anda berhasil menembus middleware keamanan!",
				"role":    role,
			})
		})
	}

	// 6. Jalankan Server di Port 8080 (Berdasarkan .env)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server HRIS berjalan lancar di port %s...\n", port)
	r.Run(":" + port)
}
